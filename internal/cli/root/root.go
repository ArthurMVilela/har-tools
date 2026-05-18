package root

import (
	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "har-tools",
	Run: execute,
}

var file string

func init() {
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path to HAR file to be read.")
	rootCmd.MarkPersistentFlagFilename("file", "har")
	rootCmd.MarkPersistentFlagRequired("file")
}

func Command() *cobra.Command {
	return rootCmd
}

func execute(cmd *cobra.Command, args []string) {
	har, err := encoding.LoadHARFromFile(file)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	out, err := encoding.EncodeToJSON(har, true)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	cmd.Println(string(out))
}
