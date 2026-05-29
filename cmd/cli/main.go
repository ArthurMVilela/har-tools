package main

import (
	"context"
	"os"

	"github.com/ArthurMVilela/har-tools/internal/cli"
	"github.com/ArthurMVilela/har-tools/internal/cli/root"
)

func main() {
	ctx := context.Background()

	deps, err := cli.NewCLIDependencies()
	if err != nil {
		os.Exit(1)
	}

	root.Command(deps).ExecuteContext(ctx)
}
