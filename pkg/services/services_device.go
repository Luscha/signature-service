package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"signature.service/pkg/device"
)

type CreateDeviceRequest struct {
	Label     string `form:"label" request:"label" xml:"label"`
	Algorithm string `form:"algorithm" request:"algorithm" xml:"algorithm" binding:"required"`
}

func (s *Server) CreateDevice(c *gin.Context) {
	var request CreateDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// TODO validate request.Algorithm

	w := s.workspaceFactory.NewWorksapce(c.Request.Context())
	device, err := w.CreateDevice(c.Request.Context(), device.Algorithm(request.Algorithm), request.Label)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, device)
}

func (s *Server) ListDevices(c *gin.Context) {
	// TODO Pagination

	w := s.workspaceFactory.NewWorksapce(c.Request.Context())
	devices, err := w.ListDevices(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, devices)
}

func (s *Server) GetDevice(c *gin.Context) {
	w := s.workspaceFactory.NewWorksapce(c.Request.Context())
	device, err := w.GetDevice(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, device)
}
