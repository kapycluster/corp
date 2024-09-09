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

	mux, err := handlers.Setup(ctx)
	if err != nil {
		l.Error("failed to setup panel server", "error", err.Error())
		return
	}

	l.Info("starting panel server", "address", "[::]:8080")
	http.ListenAndServe(":8080", mux)
}
