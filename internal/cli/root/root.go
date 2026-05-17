package root

import (
	"encoding/json"
	"os"

	"github.com/ArthurMVilela/har-tools/internal/model"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "har-tools",
	Run: execute,
}

var file string

func init() {
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "Path to HAR file to be read.")
	rootCmd.MarkFlagFilename("file", "har")
	rootCmd.MarkFlagRequired("file")
}

func Command() *cobra.Command {
	return rootCmd
}

func execute(cmd *cobra.Command, args []string) {
	var har model.HAR

	rawFile, err := os.ReadFile(file)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	err = json.Unmarshal(rawFile, &har)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	cmd.Printf("%+v\n", har.Log.Browser)
	cmd.Printf("%+v\n", har.Log.Creator)
	cmd.Printf("%+v\n", har.Log.Pages)
}
