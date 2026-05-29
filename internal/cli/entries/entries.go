package entries

import (
	"regexp"

	"github.com/ArthurMVilela/har-tools/internal/cli"
	"github.com/ArthurMVilela/har-tools/internal/cli/cmdflags"
	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/ArthurMVilela/har-tools/internal/filtering"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Command(deps *cli.CLIDependencies) *cobra.Command {
	entriesCmd := &cobra.Command{
		Use: "entries",
		Run: execute,
	}

	entriesCmd.Flags().String(cmdflags.EntriesRequestURLFilterFlag, "", "Filters out entries by the entries request's url by the given regex filter.")
	entriesCmd.Flags().String(cmdflags.EntriesXPathContentFilterFlag, "", "Apply xPath filter to entries' responses' content (body). It will only work on entries which response's content types are either json, xml or html.")
	entriesCmd.Flags().String(cmdflags.EntriesMimeTypeContentFilterFlag, "", "Filters out entries by the response's content MIME type by the given regex filter.")

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

	urlFilter, _ := cmd.Flags().GetString(cmdflags.EntriesRequestURLFilterFlag)
	if len(urlFilter) > 0 {
		if _, err := regexp.Compile(urlFilter); err != nil {
			cmd.PrintErrln(err)
			return
		}

		logger.Debug().Msgf("Applying filter: %s", urlFilter)
		filters = append(filters, filtering.URLFilter(urlFilter))
	}

	mimeFilter, _ := cmd.Flags().GetString("mime-filter-content")
	if len(mimeFilter) > 0 {
		if _, err := regexp.Compile(mimeFilter); err != nil {
			cmd.PrintErrln(err)
			return
		}

		logger.Debug().Msgf("Applying filter: %s", mimeFilter)
		filters = append(filters, filtering.MimeTypeContentFilter(mimeFilter))
	}

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

	prettyPrint, _ := cmd.Flags().GetBool("pretty")

	out, err := encoding.EncodeToJSON(filteredEntries, prettyPrint)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	cmd.Println(string(out))
}
