package server

import (
	"io"

	"github.com/NFTActions/opensea-data-backend/controllers"
	"github.com/NFTActions/opensea-data-backend/services/config"
	"github.com/NFTActions/opensea-data-backend/services/opensea"
)

type Service struct {
	config *config.AdminConfig
	osvc   *opensea.OpenseaService
}

type Controller struct {
	base *controllers.BaseController
	acti *controllers.ActivityController
}

func (server *Server) NewService() []io.Closer {
	var service Service

	service.config = config.NewService(server.db, server.log)
	service.osvc = opensea.NewOpenseaService(server.db, server.log, service.config)

	// add all service that need to be closed
	toClose := []io.Closer{}
	server.service = &service
	return toClose
}

func (server Server) NewController() *Controller {
	var controller Controller
	controller.base = controllers.NewBaseController("base", server.db, server.log, *server.config)
	controller.acti = controllers.NewActivityController(controller.base, server.service.osvc)
	return &controller
}
