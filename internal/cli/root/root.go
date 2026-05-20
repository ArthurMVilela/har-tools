package root

import (
	"github.com/ArthurMVilela/har-tools/internal/cli/entries"
	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "har-tools",
	PersistentPreRunE: persistencePreRun,
	Run:               execute,
}

var file string

func init() {
	rootCmd.AddCommand(entries.Command())

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Toogles debug mode, which logs in debug information.")

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path to HAR file to be read.")
	rootCmd.MarkPersistentFlagFilename("file", "har")
	rootCmd.MarkPersistentFlagRequired("file")
}

func Command() *cobra.Command {
	return rootCmd
}

func persistencePreRun(cmd *cobra.Command, args []string) error {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger().Level(zerolog.ErrorLevel)

	debugOn, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}

	if debugOn {
		logger = logger.Level(zerolog.DebugLevel)
		logger.Debug().Msg("Running on debug mode.")
	}

	ctx := logger.WithContext(cmd.Context())
	cmd.SetContext(ctx)

	return nil
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
