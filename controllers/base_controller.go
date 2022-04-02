/*

 */

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"opensea-data-backend/models"
)

const (
	MsgPong          = "pong"
	PER_PAGE         = 15
	PER_PAGE_LONG    = 50
	SQS_MAX_MESSAGES = 2
)

// BaseController
type BaseController struct {
	Name string
	DB   *models.DB

	log *logrus.Entry
}

func NewBaseController(name string, db *models.DB, log *logrus.Logger) *BaseController {
	return &BaseController{
		name,
		db,
		log.WithField("controller", name),
	}
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK,
		data,
	)
}

func BadRequestResponse(c *gin.Context, msg string, detail string) {
	c.JSON(
		http.StatusBadRequest,
		gin.H{
			"error": ErrorResponse{
				Code:    14000,
				Message: msg,
				Detail:  "(14000) " + msg + ": " + detail,
			},
		},
	)
}

func InternalErrorResponse(c *gin.Context, msg string, detail string) {
	c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": ErrorResponse{
				Code:    15000,
				Message: msg,
				Detail:  "(15000) " + msg + ": " + detail,
			},
		},
	)
}

func AuthErrorResponse(c *gin.Context, msg string, detail string) {
	c.JSON(
		http.StatusUnauthorized,
		gin.H{
			"error": ErrorResponse{
				Code:    14000,
				Message: msg,
				Detail:  "(14000) " + msg + ": " + detail,
			},
		},
	)
}

func CustomSuccessResponse(c *gin.Context, data interface{}, code int) {
	c.JSON(
		code,
		gin.H{
			"response": data,
		},
	)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// HandlePing handles the ping request for health check.
func (bc *BaseController) HandlePing(c *gin.Context) {
	bc.log.Debug("handling ping...")
	SuccessResponse(c, "pong")
}
