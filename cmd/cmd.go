package cmd

import (
	"github.com/armiariyan/synapsis/internal/infrastructure/container"
	"github.com/armiariyan/synapsis/internal/server"
)

func Run() {
	server.StartService(container.New())
}
