package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArthurMVilela/har-tools/pkg/model"
	"github.com/rs/zerolog"
)

// CLIDependencies contains the shared dependies shared thourgh all commands
type CLIDependencies struct {
	out, err   io.Writer
	in         io.Reader
	fsys       fs.FS
	workingDir string
	har        *model.HAR
	logger     *zerolog.Logger
}

type CLIOption func(deps *CLIDependencies) error

func NewCLIDependencies(options ...CLIOption) (*CLIDependencies, error) {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger().Level(zerolog.ErrorLevel)
	wd, _ := os.Getwd()

	deps := &CLIDependencies{
		out:        os.Stdout,
		err:        os.Stderr,
		in:         os.Stdin,
		workingDir: wd,
		fsys:       os.DirFS("/"),
		logger:     &logger,
	}

	for _, option := range options {
		err := option(deps)
		if err != nil {
			return nil, fmt.Errorf("unable to create CLI dependencies: %w", err)
		}
	}

	return deps, nil
}

// WithOut creates a CLIOption that sets the output writer. By default STDOUT will be used.
func WithOut(out io.Writer) CLIOption {
	return func(deps *CLIDependencies) error {
		deps.out = out
		return nil
	}
}

// WithErr creates a CLIOption that sets the error writer. By default STDERR will be used.
func WithErr(err io.Writer) CLIOption {
	return func(deps *CLIDependencies) error {
		deps.err = err
		return nil
	}
}

// WithIn creates a CLIOption that sets the input reader. By default STDIN will be used.
func WithIn(in io.Reader) CLIOption {
	return func(deps *CLIDependencies) error {
		deps.in = in
		return nil
	}
}

// WithWorkingDirectory creates a CLIOption that sets the working directory. By default is the current working directory.
func WithWorkingDirectory(dir string) CLIOption {
	return func(deps *CLIDependencies) error {
		deps.workingDir = dir
		return nil
	}
}

// WithFilesystem creates a CLIOption that sets the filesystem used to read the HAR file. By default it picks the root filesystem.
func WithFilesystem(fsys fs.FS) CLIOption {
	return func(deps *CLIDependencies) error {
		deps.fsys = fsys
		return nil
	}
}

func (deps *CLIDependencies) Log() *zerolog.Logger {
	return deps.logger
}

func (deps *CLIDependencies) SetLogLevel(level zerolog.Level) {
	*deps.logger = deps.logger.Level(level)
}

// LoadHAR loads the HAR file given a path, decode it and save
// the content
func (deps *CLIDependencies) LoadHAR(path string) error {
	fullpath := filepath.Join(deps.workingDir, path)
	relpath := strings.TrimPrefix(fullpath, "/")

	deps.logger.Debug().Msgf("path: %s", fullpath)

	file, err := deps.fsys.Open(relpath)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	har := new(model.HAR)
	err = json.Unmarshal(data, har)
	if err != nil {
		return fmt.Errorf("unable to parse HAR file: %w", err)
	}

	deps.har = har

	return nil
}

// HAR returns the loaded har file
func (deps *CLIDependencies) HAR() (*model.HAR, error) {
	if deps.har == nil {
		return nil, errors.New("missing har")
	}

	return deps.har, nil
}
