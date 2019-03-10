package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

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

				fmt.Printf("visited file or dir: %q\n", abspath)

			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking the path %q: %v\n", dir, err)
		}

		return nil
	},
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
