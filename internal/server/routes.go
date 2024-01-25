package server

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/CGPRE-SEPLAN-RR/fiplan-api/docs"
)

var Validate = validator.New()

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

	e.GET("/conta", s.ContaContabilHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/relatorio/fip_215", s.FIP215Handler)

	return e
}
