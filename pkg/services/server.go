package services

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"signature.service/pkg/logger"
	"signature.service/pkg/workspace"
)

type ServerConfig struct {
	Port string
}

type Server struct {
	r                *gin.Engine
	cfg              *ServerConfig
	workspaceFactory workspace.WorkspaceFactoryInterface
}

func NewServer(ctx context.Context, cfg *ServerConfig, workspaceFactory workspace.WorkspaceFactoryInterface) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/"),
		gin.Recovery(),
	)

	// tracing

	// logger
	r.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.New(
			logger.WithMinLevel(os.Getenv(logger.LOGGER_LEVEL_ENV)),
		)
		c.Request = c.Request.WithContext(logger.WithLogger(ctx, log))
		c.Next()
	})

	server := &Server{
		r:                r,
		cfg:              cfg,
		workspaceFactory: workspaceFactory,
	}

	if err := server.buildRouter(ctx); err != nil {
		logger.Main.WithError(err).Fatal("error building routes the server")
	}

	return server
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) buildRouter(ctx context.Context) error {
	s.r.GET("/", s.HealthCheck)
	s.r.POST("/devices", s.CreateDevice)
	s.r.GET("/devices", s.ListDevices)
	s.r.GET("/devices/:id", s.GetDevice)
	s.r.POST("/signatures", s.Sign)
	s.r.GET("/signatures", s.ListSignatures)

	return nil
}

func (s *Server) Run(ctx context.Context, port string) error {
	return s.r.Run(port)
}

func (s *Server) HealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
