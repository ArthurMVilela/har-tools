package entries

import (
	"strings"

	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/ArthurMVilela/har-tools/pkg/model"
	"github.com/antchfx/jsonquery"
	"github.com/spf13/cobra"
)

var entriesCmd = &cobra.Command{
	Use: "entries",
	Run: execute,
}

func init() {
	entriesCmd.Flags().String("json-content-filter", "", "Apply xPath filter to entries json content")
}

func Command() *cobra.Command {
	return entriesCmd
}

func execute(cmd *cobra.Command, args []string) {
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

	jsonFilter, _ := cmd.Flags().GetString("json-content-filter")
	if len(jsonFilter) > 0 {
		entries, err = filterEntries(entries, jsonFilter)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	}

	out, err := encoding.EncodeToJSON(entries, true)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	cmd.Println(string(out))
}

func filterEntries(entries []model.Entry, filter string) ([]model.Entry, error) {
	var filtered []model.Entry

	for _, entry := range entries {
		if !strings.Contains(entry.Response.Content.MimeType, "json") {
			continue
		}
		reader := strings.NewReader(entry.Response.Content.Text)
		node, err := jsonquery.Parse(reader)
		if err != nil {
			continue
		}

		node, err = jsonquery.Query(node, filter)
		if err != nil {
			continue
		}
		if node != nil {
			filtered = append(filtered, entry)
		}
	}

	return filtered, nil
}
