package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var directory string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "walker",
	Short: "Generates vegeta friendly test data",
	Long: `Walks a directory tree matching files against a supplied regex
expression to feed into vegeta`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			}

			if !info.IsDir() {
				fmt.Printf("visited file: %q\n", path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking the path %q: %v\n", directory, err)
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

	rootCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory to walk")
	rootCmd.MarkFlagRequired("directory")
}
