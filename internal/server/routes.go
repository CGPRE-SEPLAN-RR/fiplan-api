package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/CGPRE-SEPLAN-RR/fiplan-api/docs"
	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal"
	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/database"
	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/model"
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

	e.GET("/conta", s.ContaContabilHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

func (s *Server) ContaContabilHandler(c echo.Context) error {
	var contas []model.ContaContabil

	sqlQuery := "SELECT CODG_CONTA_CONTABIL, NOME_CONTA_CONTABIL FROM ACWTB0032 FETCH FIRST 50 ROWS ONLY"

	rows, err := s.db.Query(sqlQuery)

	if err != nil {
		log.Println(err)
		return c.JSON(
			http.StatusInternalServerError,
			internal.BasicResponse("Erro ao consultar as contas contábeis"),
		)
	}

	defer rows.Close()

	for rows.Next() {
		var conta model.ContaContabil

		if err := rows.Scan(&conta.Codigo, &conta.Nome); err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				internal.BasicResponse("Erro ao consultar uma das contas contábeis"),
			)

		}

		contas = append(contas, conta)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			internal.BasicResponse("Não sei que erro é esse"),
		)
	}

	return c.JSON(http.StatusOK, contas)
}
