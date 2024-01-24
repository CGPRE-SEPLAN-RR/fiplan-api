package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/model"
	"github.com/labstack/echo/v4"
)

// @ContaContabilHandler godoc
// @Summary     Lista as contas
// @Description Teste
// @Tags        Conta
// @Accept      json
// @Produce     json
// @Param       ano_exercicio query    int8 true "Ano de Exercício"
// @Success     200           {array}  model.ContaContabil
// @Failure     400           {object} Erro
// @Failure     404           {object} Erro
// @Failure     500           {object} Erro
// @Router      /conta [get]
func (s *Server) ContaContabilHandler(c echo.Context) error {
	/*** Parâmetros ***/
	var anoExercicio int16
	/*** Parâmetros ***/

	/*** Validação dos Parâmetros ***/
	if err := echo.QueryParamsBinder(c).MustInt16("ano_exercicio", &anoExercicio).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Por favor, forneça o ano de exercício no parâmetro 'ano_exercicio'.",
		)
	}

	if err := Validate.Var(anoExercicio, fmt.Sprintf("gte=2010,lte=%d", time.Now().Year())); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Por favor, forneça um ano de exercício válido entre 2010 e %d para o parâmetro 'ano_exercicio'.", time.Now().Year()),
		)
	}
	/*** Validação dos Parâmetros ***/

	/*** Consulta no Banco de Dados ***/
	var contas []model.ContaContabil

	sqlQuery := `SELECT CODG_CONTA_CONTABIL,NOME_CONTA_CONTABIL
							 FROM ACWTB0032
							 WHERE CD_EXERCICIO = :1
							 ORDER BY CODG_CONTA_CONTABIL ASC`

	rows, err := s.db.Query(sqlQuery, anoExercicio)

	if err != nil {
		log.Printf("ContaContabilHandler: %v", err)

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Ocorreu um erro ao consultar o banco de dados.",
		)
	}

	defer rows.Close()

	for rows.Next() {
		var conta model.ContaContabil

		if err := rows.Scan(&conta.Codigo, &conta.Nome); err != nil {
			log.Printf("ContaContabilHandler: %v", err)

			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"Ocorreu um erro ao consultar uma linha no banco de dados.",
			)
		}

		contas = append(contas, conta)
	}

	if err := rows.Err(); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Ocorreu um erro de rede ou problema no resultado do banco de dados.",
		)
	}
	/*** Consulta no Banco de Dados ***/

	return c.JSON(http.StatusOK, contas)
}
