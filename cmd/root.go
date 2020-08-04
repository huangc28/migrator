package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "migrator",
	Short: "automate the process of various database migration actions",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("some command invoked")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
