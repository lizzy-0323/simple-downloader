package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "downloader",
		Short: "A tool for downloading files",
		Long:  "A tool for downloading files, currently supports http and ftp protocols",
	}
	rootCmd.AddCommand(NewDownloadCmd())
	return rootCmd
}
