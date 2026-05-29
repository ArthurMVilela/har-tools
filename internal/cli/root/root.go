package root

import (
	"fmt"

	"github.com/ArthurMVilela/har-tools/internal/cli"
	"github.com/ArthurMVilela/har-tools/internal/cli/cmdflags"
	"github.com/ArthurMVilela/har-tools/internal/cli/entries"
	"github.com/ArthurMVilela/har-tools/internal/encoding"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Command(deps *cli.CLIDependencies) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "har-tools",
		PersistentPreRunE: persistencePreRun(deps),
		RunE:              execute(deps),
	}

	rootCmd.AddCommand(entries.Command(deps))

	rootCmd.PersistentFlags().BoolP(cmdflags.DebugFlag, cmdflags.DebugShortFlag, false, "Toogles debug mode, which logs in debug information.")
	rootCmd.PersistentFlags().StringP(cmdflags.FileFlag, cmdflags.FileShortFlag, "", "Path to HAR file to be processed.")
	rootCmd.PersistentFlags().BoolP(cmdflags.PrettyFlag, cmdflags.PrettyShortFlag, false, "Toogles pretty outpput: with line breaking and indentation. Ideal for human readability.")

	rootCmd.MarkPersistentFlagFilename(cmdflags.FileFlag, "har")
	rootCmd.MarkPersistentFlagRequired(cmdflags.FileFlag)

	return rootCmd
}

func persistencePreRun(deps *cli.CLIDependencies) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		debugOn, err := cmd.Flags().GetBool(cmdflags.DebugFlag)
		if err != nil {
			return fmt.Errorf("unable to get %s flag: %w", cmdflags.DebugFlag, err)
		}
		if debugOn {
			deps.SetLogLevel(zerolog.DebugLevel)
		}

		file, err := cmd.Flags().GetString(cmdflags.FileFlag)
		if err != nil {
			return fmt.Errorf("unable to get %s flag: %w", cmdflags.FileFlag, err)
		}

		err = deps.LoadHAR(file)
		if err != nil {
			return err
		}

		return nil
	}
}

func execute(deps *cli.CLIDependencies) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		har, err := deps.HAR()
		if err != nil {
			return fmt.Errorf("unable to get HAR: %w", err)
		}

		pretty, err := cmd.Flags().GetBool(cmdflags.PrettyFlag)
		if err != nil {
			return fmt.Errorf("unable to get %s flag: %w", cmdflags.FileFlag, err)
		}

		out, err := encoding.EncodeToJSON(har, pretty)
		if err != nil {
			return fmt.Errorf("unable to encode output: %w", err)
		}

		cmd.Println(string(out))
		return nil
	}
}
