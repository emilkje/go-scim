package main

import (
	"log/slog"
	"net/http"

	"github.com/intility/scim/server"
)

func main() {
	server := server.NewServer(logger{})
	err := http.ListenAndServe("127.0.0.1:1337", server)
	if err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

type logger struct{}

// Error implements scim.Logger.
func (l logger) Error(args ...interface{}) {
	slog.Error("SCIM ERROR", args...)
}
