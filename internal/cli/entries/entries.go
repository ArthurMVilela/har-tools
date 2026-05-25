package entries

import (
	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/ArthurMVilela/har-tools/internal/filtering"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var entriesCmd = &cobra.Command{
	Use: "entries",
	Run: execute,
}

func init() {
	entriesCmd.Flags().String("xpath-filter-content", "", "Apply xPath filter to entries' responses' content (body). It will only work on entries which response's content types are either json, xml or html.")
}

func Command() *cobra.Command {
	return entriesCmd
}

func execute(cmd *cobra.Command, args []string) {
	logger := zerolog.Ctx(cmd.Context())

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	har, err := encoding.LoadHARFromFile(file)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	entries := har.Log.Entries

	var filters []filtering.EntryFilter

	jsonFilter, _ := cmd.Flags().GetString("xpath-filter-content")
	if len(jsonFilter) > 0 {
		filters = append(filters, filtering.XPathEntryContentFilter(jsonFilter))
	}

	processor, err := filtering.NewEntryProcessor(filtering.WithLogger(logger), filtering.WithEntryFilters(filters...))
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	filteredEntries, err := processor.ApplyFilters(entries)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	out, err := encoding.EncodeToJSON(filteredEntries, true)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	cmd.Println(string(out))
}
