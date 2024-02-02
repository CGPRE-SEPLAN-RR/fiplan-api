package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/CGPRE-SEPLAN-RR/fiplan-api/docs"
	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/handlers"
)

// @title API do FIPLAN
// @version 0.0.0-alpha
// @description API para consulta de dados e relatórios do FIPLAN em JSON

// @contact.name Equipe da CGPRE (SEPLAN/RR)
// @contact.email cgpre@planejamento.rr.gov.br
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "swagger")
		},
	}))

	e.GET("/conta", handlers.ContaContabilHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/relatorio/fip_215", handlers.RelatorioFIP215Handler)
	e.GET("/relatorio/fip_215m", handlers.RelatorioFIP215MHandler)

	return e
}
