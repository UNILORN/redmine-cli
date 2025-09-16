package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "redmine",
	Short: "Redmine CLI tool",
	Long:  `A command-line interface for managing Redmine issues and projects`,
}

var profileFlag string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "p", "", "Profile to use for this command")
}