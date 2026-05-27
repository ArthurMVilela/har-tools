package root

import (
	"github.com/ArthurMVilela/har-tools/internal/cli/cmdflags"
	"github.com/ArthurMVilela/har-tools/internal/cli/entries"
	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "har-tools",
		PersistentPreRunE: persistencePreRun,
		RunE:              execute,
	}

	rootCmd.AddCommand(entries.Command())

	rootCmd.PersistentFlags().BoolP(cmdflags.DebugFlag, cmdflags.DebugShortFlag, false, "Toogles debug mode, which logs in debug information.")
	rootCmd.PersistentFlags().StringP(cmdflags.FileFlag, cmdflags.FileShortFlag, "", "Path to HAR file to be processed.")
	rootCmd.PersistentFlags().BoolP(cmdflags.PrettyFlag, cmdflags.PrettyShortFlag, false, "Toogles pretty outpput: with line breaking and indentation. Ideal for human readability.")

	rootCmd.MarkPersistentFlagFilename(cmdflags.FileFlag, "har")
	rootCmd.MarkPersistentFlagRequired(cmdflags.FileFlag)

	return rootCmd
}

func persistencePreRun(cmd *cobra.Command, args []string) error {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger().Level(zerolog.ErrorLevel)

	debugOn, err := cmd.Flags().GetBool(cmdflags.DebugFlag)
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

func execute(cmd *cobra.Command, args []string) error {
	filepath, err := cmd.Flags().GetString(cmdflags.FileFlag)
	if err != nil {
		return err
	}

	har, err := encoding.LoadHARFromFile(filepath)
	if err != nil {
		return err
	}

	out, err := encoding.EncodeToJSON(har, true)
	if err != nil {
		return err
	}

	cmd.Println(string(out))

	return nil
}
