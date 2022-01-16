package controllers

import (
	"github.com/NFTActions/opensea-data-backend/services/opensea"
	"github.com/gin-gonic/gin"
)

type ActivityController struct {
	BaseController
	osvc *opensea.OpenseaService
}

func NewActivityController(bc *BaseController, osvc *opensea.OpenseaService) *ActivityController {
	return &ActivityController{
		BaseController{
			Name: "activity",
			DB:   bc.DB,
			log:  bc.log,
			conf: bc.conf,
		},
		osvc,
	}
}

func (ac *ActivityController) HandleActivityCreate(c *gin.Context) {
	recentActivities, err := ac.osvc.GetRecentOpenseaEvents()
	if err != nil {
		InternalErrorResponse(c, "error in GetRecentEvents", err.Error())
		return
	}

	if len(recentActivities) > 0 {
		err = ac.DB.BatchInsertActivity(recentActivities)
		if err != nil {
			InternalErrorResponse(c, "error in BatchInsert", err.Error())
			return
		}
	}

	SuccessResponse(c, gin.H{
		"activities": recentActivities,
	})
}
