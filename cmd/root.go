package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/spf13/cobra"
	"github.com/tsenart/vegeta/lib"
)

func init() {
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".tiff", "image/tiff")
	mime.AddExtensionType(".tif", "image/tiff")
}

var dir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "walker",
	Short: "Generates vegeta friendly test data",
	Long: `Walks a directory tree matching files against a supplied regex
expression to feed into vegeta`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				abspath, err := filepath.Abs(path)

				if err != nil {
					return fmt.Errorf("unable to get absolute path for %q", path)
				}

				target, err := asTarget(abspath)
				if err != nil {
					return err
				}

				var buf bytes.Buffer
				enc := vegeta.NewJSONTargetEncoder(&buf)
				err = enc.Encode(target)

				if err != nil {
					return err
				}
				fmt.Printf("%s\n", string(buf.Bytes()))

			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking the path %q: %v\n", dir, err)
		}

		return nil
	},
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&dir, "directory", "d", "", "Directory to walk")
	rootCmd.MarkFlagRequired("directory")
}
