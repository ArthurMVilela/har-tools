package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "har-tools",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Hello")
	},
}

func main() {
	rootCmd.Execute()
}
