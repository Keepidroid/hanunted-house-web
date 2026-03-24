package main

import (
	"log/slog"
	"os"

	"hanunted-house/server/internal/httpapi"
	"hanunted-house/server/internal/room"
	"hanunted-house/server/internal/ws"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	manager := room.NewManager()
	hub := ws.NewHub(manager, logger)
	api := httpapi.New(manager, hub, logger)

	addr := ":8080"
	if v := os.Getenv("SERVER_ADDR"); v != "" {
		addr = v
	}
	logger.Info("server starting", "addr", addr)
	if err := api.Router().Run(addr); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
