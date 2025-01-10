package main

import (
	"context"
	"fmt"
	"net/http"

	"kapycluster.com/corp/log"
	"kapycluster.com/corp/panel/config"
	"kapycluster.com/corp/panel/handlers"
)

func main() {
	ctx := log.NewContext(context.Background(), "panel")
	l := log.FromContext(ctx)

	c := config.NewConfig()

	if c.Kubernetes.KubeconfigsDir == "" {
		l.Error("kubeconfigs path is required", "key", config.AsEnv("kubernetes.kubeconfigs"))
		return
	}

	mux, err := handlers.Setup(ctx, c)
	if err != nil {
		l.Error("failed to setup panel server", "error", err.Error())
		return
	}

	listenAddr := fmt.Sprintf("%s:%d", c.Server.ListenHost, c.Server.ListenPort)

	l.Info("starting panel server", "address", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
