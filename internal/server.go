package internal

import (
	"ch-whisper/config"
	"ch-whisper/internal/core"
	"ch-whisper/pkg/api"
	"go.uber.org/fx"
)

func RunServer() {
	config.SetViper()
	fx.New(
		fx.Provide(
			core.NewModel,
			api.NewServer,
		),
		fx.Invoke(
			func(server *api.Server) { server.Serve() },
		),
	).Run()
}
