package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type contaContabil struct {
	Codigo string `json:"codigo"`
	Nome   string `json:"nome"`
} // @name ContaContabil

type relatorioContas struct {
	Dados []contaContabil
} // @name RelatorioContas

// ContaContabilHandler godoc
//
// @Summary     Contas contábeis
// @Description Lista as contas contábeis
// @Tags        Conta
// @Accept      json
// @Produce     json
// @Param       ano_exercicio query    uint16 true "Ano de Exercício"
// @Success     200           {array}  relatorioContas
// @Failure     400           {object} Erro
// @Failure     500           {object} Erro
// @Router      /conta [get]
func (s *Server) ContaContabilHandler(c echo.Context) error {
	/*** Parâmetros ***/
	parametros := struct {
		AnoExercicio uint16
	}{}
	/*** Parâmetros ***/

	/*** Validação dos Parâmetros ***/
	valueBinder := echo.QueryParamsBinder(c)

	var errors []string

	if err := valueBinder.MustUint16("ano_exercicio", &parametros.AnoExercicio).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o ano de exercício no parâmetro 'ano_exercicio'.")
	}

	if err := Validate.Var(parametros.AnoExercicio, fmt.Sprintf("gte=2010,lte=%d", time.Now().Year())); err != nil {
		errors = append(errors, fmt.Sprintf("Por favor, forneça um ano de exercício válido entre 2010 e %d para o parâmetro 'ano_exercicio'.", time.Now().Year()))
	}

	if len(errors) > 0 {
		return ErroValidacaoParametro(errors)
	}
	/*** Validação dos Parâmetros ***/

	/*** Consulta no Banco de Dados ***/
	var sqlQuery strings.Builder

	var contasContabeis relatorioContas

	queryTemplate := `SELECT CODG_CONTA_CONTABIL,NOME_CONTA_CONTABIL
							      FROM ACWTB0032
							      WHERE CD_EXERCICIO = {{.AnoExercicio}}
							      ORDER BY CODG_CONTA_CONTABIL ASC`

	tmplContaExplosao, err := template.New("queryContasContabeis").Parse(queryTemplate)

	if err != nil {
		log.Printf("ContaContabilHandler: %v", err)
		return ErroMontagemTemplate
	}

	err = tmplContaExplosao.Execute(&sqlQuery, parametros)

	if err != nil {
		log.Printf("ContaContabilHandler: %v", err)
		return ErroExecucaoTemplate
	}

	compactSqlQuery := strings.Join(strings.Fields(sqlQuery.String()), " ")
	log.Printf("ContaContabilHandler: %s", compactSqlQuery)
	rows, err := s.db.Query(compactSqlQuery)

	sqlQuery.Reset()

	if err != nil {
		log.Printf("ContaContabilHandler: %v", err)
		return ErroConsultaBancoDados
	}

	defer rows.Close()

	for rows.Next() {
		var conta contaContabil

		if err := rows.Scan(&conta.Codigo, &conta.Nome); err != nil {
			log.Printf("ContaContabilHandler: %v", err)
			return ErroConsultaLinhaBancoDados
		}

		contasContabeis.Dados = append(contasContabeis.Dados, conta)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ContaContabilHandler: %v", err)
		return ErroRedeOuResultadoBancoDados
	}
	/*** Consulta no Banco de Dados ***/

	return c.JSON(http.StatusOK, contasContabeis)
}
