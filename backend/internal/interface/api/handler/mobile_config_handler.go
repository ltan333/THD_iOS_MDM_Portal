package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type MobileConfigHandler interface {
	// GetByID(c *gin.Context)
	GetXML(c *gin.Context)
}

type mobileConfigHandlerImpl struct {
	mobileConfigService service.MobileConfigService
}

func NewMobileConfigHandler(mobileConfigService service.MobileConfigService) MobileConfigHandler {
	return &mobileConfigHandlerImpl{mobileConfigService: mobileConfigService}
}

// GetXML godoc
// @Summary Export mobile config XML
// @Description Generate and return raw Apple mobileconfig XML content by ID
// @Tags Mobile Config
// @Produce xml
// @Param id path int true "Mobile config ID"
// @Success 200 {string} string "Raw XML"
// @Failure 400 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /api/mobile-configs/{id}/xml [get]
func (m *mobileConfigHandlerImpl) GetXML(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	cmd := service.GenerateMobileConfigXMLCommand{ID: uint(id)}
	xmlBytes, err := m.mobileConfigService.GenerateXML(c.Request.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại"))
			return
		}
		response.WriteErrorResponse(c, err)
		return
	}

	// Return raw XML
	c.Data(http.StatusOK, "text/xml; charset=utf-8", xmlBytes)
}
