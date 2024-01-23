package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/CGPRE-SEPLAN-RR/fiplan-api/docs"
)

// @title API do FIPLAN
// @version 0.0.0-alpha
// @description API para consulta de dados e relatórios do FIPLAN em JSON

// @contact.name Equipe da CGPRE (SEPLAN/RR)
// @contact.email cgpre@planejamento.rr.gov.br
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "swagger")
		},
	}))

	e.GET("/", s.HelloWorldHandler)
	e.GET("/health", s.HealthHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
