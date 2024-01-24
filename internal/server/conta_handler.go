package server

import (
	"net/http"

	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal"
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
// @Failure     400           {object} echo.HTTPError
// @Failure     404           {object} echo.HTTPError
// @Failure     500           {object} echo.HTTPError
// @Router      /conta [get]
func (s *Server) ContaContabilHandler(c echo.Context) error {
	var anoExercicio int16

	if err := echo.QueryParamsBinder(c).MustInt16("ano_exercicio", &anoExercicio).BindError(); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			internal.BasicResponse("O parâmetro de query 'ano_exercicio' é obrigatório"),
		)
	}

	if err := Validate.Var(anoExercicio, "gte=2010"); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			internal.BasicResponse("O parâmetro de query 'ano_exercicio' deve ser um número no intervalo [2010,2024]"),
		)
	}

	var contas []model.ContaContabil

	sqlQuery := `SELECT CODG_CONTA_CONTABIL,NOME_CONTA_CONTABIL
							 FROM ACWTB0032
							 WHERE CD_EXERCICIO = :1
							 ORDER BY CODG_CONTA_CONTABIL ASC`

	rows, err := s.db.Query(sqlQuery, anoExercicio)

	if err != nil {
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
