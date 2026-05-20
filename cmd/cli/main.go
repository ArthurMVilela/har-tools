package main

import (
	"context"

	"github.com/ArthurMVilela/har-tools/internal/cli/root"
)

func main() {
	ctx := context.Background()

	root.Command().ExecuteContext(ctx)
}
