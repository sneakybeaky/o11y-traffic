package cmd

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/spf13/cobra"
	"github.com/tsenart/vegeta/lib"
)

var dir string
var forever bool

var types = map[string]string{
	".jpeg": "image/jpeg",
	".tiff": "image/tiff",
	".tif":  "image/tiff",
	".gif":  "image/gif",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".svg":  "image/svg+xml",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "walker",
	Short: "Generates vegeta friendly test data",
	Long: `Walks a directory tree matching files against a supplied regex
expression to feed into vegeta`,
	RunE: func(cmd *cobra.Command, args []string) error {

		for ext, mimeType := range types {
			err := mime.AddExtensionType(ext, mimeType)

			if err != nil {
				return fmt.Errorf("unable to setup mime types: %v\n", err)
			}
		}

		found, err := findImages(dir)

		if err != nil {
			return fmt.Errorf("unable to find images: %v\n", err)
		}

		for {
			shuffle(found)
			asTargets(found)

			if !forever {
				break
			}
		}
		return nil
	},
}

func asTargets(paths []string) {

	var buf bytes.Buffer
	enc := vegeta.NewJSONTargetEncoder(&buf)

	for _, path := range paths {

		target, err := asTarget(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to process %q : %v", path, err)
			break
		}

		err = enc.Encode(target)

		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to encode %q : %v", path, err)
			break
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(buf.Bytes()))
		buf.Reset()

	}

}

func findImages(dir string) ([]string, error) {

	var found []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isImage := isImageFile(path); !isImage {
			return nil
		}

		abspath, err := filepath.Abs(path)

		if err != nil {
			return fmt.Errorf("unable to get absolute path for %q", path)
		}

		found = append(found, abspath)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return found, nil

}

func isImageFile(path string) bool {
	_, ok := types[filepath.Ext(path)]
	return ok
}

func asTarget(path string) (*vegeta.Target, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(createFormFile("file", filepath.Base(path)))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	writer.WriteField("name", fmt.Sprintf("%s %s %s", randomdata.Adjective(), randomdata.Noun(), randomdata.FullDate()))

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	target := vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080/api/images",
		Header: http.Header{"Content-Type": []string{writer.FormDataContentType()}},
		Body:   []byte(body.Bytes()),
	}

	return &target, nil

}

func createFormFile(fieldname, filename string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			fieldname, filename))
	h.Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	return h
}

func shuffle(vals []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}

func Execute() {

	rootCmd.Flags().StringVarP(&dir, "directory", "d", "", "Directory to walk")
	rootCmd.MarkFlagRequired("directory")

	rootCmd.Flags().BoolVarP(&forever, "forever", "f", true, "Run forever")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
