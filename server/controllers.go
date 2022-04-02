package server

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"io"

	"opensea-data-backend/controllers"
	"opensea-data-backend/services/config"
	"opensea-data-backend/services/opensea"
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

	service.config = config.NewAdminConfig(server.db, server.log)
	service.osvc = opensea.NewOpenseaService(server.db, server.log, service.config)

	StartJob(service.osvc)

	// add all service that need to be closed
	toClose := []io.Closer{}
	server.service = &service
	return toClose
}

func StartJob(osvc *opensea.OpenseaService) {
	job := cron.New()
	_, err := job.AddFunc("* * * * *", func() {
		_, err := osvc.GetRecentOpenseaEvents()
		if err != nil {
			fmt.Errorf("error when running job:%s", err)
			return
		}
	})
	if err != nil {
		return
	}
	job.Start()
}

func (server Server) NewController() *Controller {
	var controller Controller
	controller.base = controllers.NewBaseController("base", server.db, server.log)
	controller.acti = controllers.NewActivityController(controller.base, server.service.osvc)
	return &controller
}
