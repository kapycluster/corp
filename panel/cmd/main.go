package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kapycluster/corpy/log"
	"github.com/kapycluster/corpy/panel/config"
	"github.com/kapycluster/corpy/panel/handlers"
)

func main() {
	ctx := log.NewContext(context.Background(), "panel")
	l := log.FromContext(ctx)

	config := config.NewConfig()

	mux, err := handlers.Setup(ctx, config)
	if err != nil {
		l.Error("failed to setup panel server", "error", err.Error())
		return
	}

	listenAddr := fmt.Sprintf("%s:%d", config.Server.ListenHost, config.Server.ListenPort)

	l.Info("starting panel server", "address", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
