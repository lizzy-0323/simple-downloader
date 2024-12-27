package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "downloader",
		Short: "A tool for downloading datasets for ANN search.",
		Long:  "A tool for downloading datasets for ANN search.",
	}
	rootCmd.AddCommand(NewDownloadCmd())
	return rootCmd
}
