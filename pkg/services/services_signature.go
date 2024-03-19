package services

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"signature.service/pkg/workspace"
)

type SignatureRequest struct {
	Data   string `form:"data" request:"data" xml:"data" binding:"required"`
	Device string `form:"device" request:"device" xml:"device" binding:"required"`
}

func (s *Server) Sign(c *gin.Context) {
	var request SignatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	w := s.workspaceFactory.NewWorksapce(c.Request.Context())
	signature, err := w.Sign(c.Request.Context(), request.Device, []byte(request.Data))
	if err != nil {
		if errors.Is(err, workspace.ErrDeviceNotFound) {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signature)
}

func (s *Server) ListSignatures(c *gin.Context) {
	// TODO Pagination

	w := s.workspaceFactory.NewWorksapce(c.Request.Context())
	signatures, err := w.ListSignatures(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signatures)
}
