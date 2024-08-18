package main

import (
	"context"

	"github.com/kapycluster/corpy/log"
)

func main() {
	ctx := log.NewContext(context.Background(), "panel")
	l := log.FromContext(ctx)

	l.Info("starting panel server", "address", "[::]:8080")
}
