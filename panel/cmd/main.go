package main

import (
	"context"
	"net/http"

	"github.com/kapycluster/corpy/log"
	"github.com/kapycluster/corpy/panel/handlers"
)

func main() {
	ctx := log.NewContext(context.Background(), "panel")
	l := log.FromContext(ctx)

	mux := handlers.Setup(ctx)

	l.Info("starting panel server", "address", "[::]:8080")
	http.ListenAndServe(":8080", mux)
}
