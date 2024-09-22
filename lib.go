package gloeggnats

import (
	"fmt"

	"github.com/Meduzz/gloegg"
	"github.com/Meduzz/gloegg-nats/flags"
	"github.com/Meduzz/gloegg-nats/logs"
	"github.com/nats-io/nats.go"
)

type (
	Config struct {
		nc      *nats.Conn
		project string
		service string
	}

	SetupOpt func(*Config)
)

func Setup(nc *nats.Conn, project, service string) error {
	config := &Config{
		nc:      nc,
		project: project,
		service: service,
	}

	if config.project != "" {
		gloegg.AddMeta("@project", config.project)
	} else {
		return fmt.Errorf("no project specified")
	}

	if config.service != "" {
		gloegg.AddMeta("@service", config.service)
	} else {
		return fmt.Errorf("no service specified")
	}

	factory := logs.NewLogFactory(config.nc, config.project)
	gloegg.SetHandlerFactory(factory)

	flagSync := flags.NewFlagHandler(config.nc, config.project)
	flagSync.Setup()

	return nil
}
