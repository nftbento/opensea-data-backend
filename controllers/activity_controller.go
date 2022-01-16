package controllers

import (
	"sort"
	"strconv"
	"time"

	"github.com/NFTActions/opensea-data-backend/models"
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

func (ac *ActivityController) HandleGetActivitySummary(c *gin.Context) {
	periodString := c.DefaultQuery("period", "60")
	period, err := strconv.Atoi(periodString)
	if err != nil {
		BadRequestResponse(c, "period parameter is wrong", err.Error())
		return
	}

	activities, err := ac.DB.GetActivitiesAfter(time.Now().Add(-time.Duration(period) * time.Minute))
	if err != nil {
		InternalErrorResponse(c, "error in GetActivitiesAfter", err.Error())
		return
	}

	collectionMap := GetCollectionMap(activities)
	keys := make([]string, 0, len(collectionMap))
	for k := range collectionMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return collectionMap[keys[i]].Count > collectionMap[keys[j]].Count })

	sortedCollectionSummary := make([]Collection, 0, len(keys))
	for _, k := range keys {
		hottestCollection := collectionMap[k]
		sortedCollectionSummary = append(sortedCollectionSummary, hottestCollection)
	}
	SuccessResponse(c, gin.H{
		"collections": sortedCollectionSummary,
	})
}

func GetCollectionMap(activities []models.Activity) map[string]Collection {
	collectionMap := make(map[string]Collection, 0)
	for _, a := range activities {
		// First we get a "copy" of the entry
		if entry, ok := collectionMap[a.CollectionSlug]; ok {
			// Then we modify the copy
			entry.Count += 1
			entry.TotalSalesInGwei += a.TotalPrice

			// Then we reassign map entry
			collectionMap[a.CollectionSlug] = entry
		} else {
			newCollection := Collection{
				Name:             a.CollectionName,
				Count:            1,
				TotalSalesInGwei: a.TotalPrice,
				ImageUrl:         a.CollectionImageUrl,
				CreatedDate:      a.CollectionCreatedDate,
			}
			collectionMap[a.CollectionSlug] = newCollection
		}
	}
	return collectionMap
}

type Collection struct {
	Name             string    `json:"name"`
	Count            int       `json:"count"`
	TotalSalesInGwei int64     `json:"total_sales_in_gwei"`
	ImageUrl         string    `json:"image_url"`
	CreatedDate      time.Time `json:"created_date"`
}
