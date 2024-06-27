package server

import (
	"github.com/armiariyan/synapsis/internal/infrastructure/container"
	"github.com/armiariyan/synapsis/internal/server/http"
)

func StartService(container *container.Container) {
	http.StartH2CServer(container)
}
