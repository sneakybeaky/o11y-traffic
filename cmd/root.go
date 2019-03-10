package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var glob string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "walker",
	Short: "Generates vegeta friendly test data",
	Long: `Walks a directory tree matching files against a supplied regex
expression to feed into vegeta`,
	RunE: func(cmd *cobra.Command, args []string) error {

		matches, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, match := range matches {
			fmt.Printf("Found a matching file at %s\n", match)
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

	rootCmd.Flags().StringVarP(&glob, "glob", "g", "", "Glob to match")
	rootCmd.MarkFlagRequired("glob")
}
