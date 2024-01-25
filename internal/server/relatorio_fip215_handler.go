package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/model"
	"github.com/labstack/echo/v4"
)

// @FIP215Handler godoc
// @Summary     FIP215 - Balancete Mensal de Verificação
// @Description Teste
// @Tags        Relatório
// @Accept      json
// @Produce     json
// @Param       ano_exercicio                    query    uint16 true  "Ano de Exercício"
// @Param       unidade_gestora                  query    uint16 false "Unidade Gestora"
// @Param       mes_referencia                   query    uint8  true  "Mês de Referência"                  Enums(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
// @Param       mes_contabil                     query    uint8  true  "Mês Contábil"                       Enums(1, 2, 3, 4)
// @Param       tipo_encerramento                query    uint8  false "Tipo de Encerramento"               Enums(1, 2)
// @Param       indicativo_conta_contabil_rp     query    bool   false "Indicativo de Conta Contábil de RP" Enums(true, false)
// @Param       indicativo_superavit_fincanceiro query    bool   false "Indicativo de Superávit Financeiro" Enums(true, false)
// @Param       indicativo_composicao_msc        query    bool   false "Indicativo de Composição da MSC"    Enums(true, false)
// @Success     200                              {array}  model.RelatorioFIP215
// @Failure     400                              {object} Erro
// @Failure     404                              {object} Erro
// @Failure     500                              {object} Erro
// @Router      /relatorio/fip_215 [get]
func (s *Server) FIP215Handler(c echo.Context) error {
	/*** Parâmetros ***/
	var anoExercicio uint16
	var unidadeGestora uint16
	var mesDeReferencia uint8
	var mesContabil uint8
	var tipoDeEncerramento uint8
	indicativoComposicaoMSC := false
	indicativoContaContabilRP := false
	indicativoSuperavitFinanceiro := false
	/*** Parâmetros ***/

	/*** Validação dos Parâmetros ***/
	if err := echo.QueryParamsBinder(c).MustUint16("ano_exercicio", &anoExercicio).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o ano de exercício no parâmetro 'ano_exercicio'.",
		)
	}

	if err := echo.QueryParamsBinder(c).Uint16("unidade_gestora", &unidadeGestora).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça a unidade gestora no parâmetro 'unidade_gestora'.",
		)
	}

	if err := echo.QueryParamsBinder(c).MustUint8("mes_referencia", &mesDeReferencia).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o mês de referência no parâmetro 'mes_referencia'.",
		)
	}

	if err := echo.QueryParamsBinder(c).MustUint8("mes_contabil", &mesContabil).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o mês contábil no parâmetro 'mes_contabil'.",
		)
	}

	if err := echo.QueryParamsBinder(c).Uint8("tipo_encerramento", &tipoDeEncerramento).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o tipo de encerramento no parâmetro 'tipo_encerramento'.",
		)
	}

	if err := echo.QueryParamsBinder(c).Bool("indicativo_conta_contabil_rp", &indicativoContaContabilRP).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o indicativo de conta contábil de RP no parâmetro 'indicativo_conta_contabil_rp'.",
		)
	}

	if err := echo.QueryParamsBinder(c).Bool("indicativo_superavit_fincanceiro", &indicativoSuperavitFinanceiro).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o indicativo de superávit financeiro no parâmetro 'indicativo_superavit_fincanceiro'.",
		)
	}

	if err := echo.QueryParamsBinder(c).Bool("indicativo_composicao_msc", &indicativoComposicaoMSC).BindError(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça o indicativo de composição da MSC no parâmetro 'indicativoComposicaoMSC'.",
		)
	}

	if err := Validate.Var(anoExercicio, fmt.Sprintf("gte=2010,lte=%d", time.Now().Year())); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Por favor, forneça um ano de exercício válido entre 2010 e %d para o parâmetro 'ano_exercicio'.", time.Now().Year()),
		)
	}

	if err := Validate.Var(mesDeReferencia, fmt.Sprintf("gte=1,lte=12")); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça um mês de referência válido entre 1 e 12 para o parâmetro 'mes_referencia'.",
		)
	}

	if err := Validate.Var(mesContabil, fmt.Sprintf("gte=1,lte=4")); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça um mês contábil válido entre 1 e 4 para o parâmetro 'mes_contabil'.",
		)
	}

	if err := Validate.Var(tipoDeEncerramento, fmt.Sprintf("gte=1,lte=2")); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Por favor, forneça um tipo de encerramento válido entre 1 e 2 para o parâmetro 'tipo_encerramento'.",
		)
	}
	/*** Validação dos Parâmetros ***/

	/*** Consulta no Banco de Dados ***/
	var relatorio model.RelatorioFIP215

	sqlQuery, err := LerSQL("relatorio_fip215.sql")

	if err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroLeituraSQL
	}

	rows, err := s.db.Query(sqlQuery)

	if err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroConsultaBancoDados
	}

	defer rows.Close()

	for rows.Next() {
		var dado model.DadoRelatorioFIP215

		if err := rows.Scan(
			&dado.CodigoUnidadeOrcamentaria,
			&dado.NomeUnidadeOrcamentaria,
			&dado.IDContaContabil,
			&dado.IDContaContabilExplosao,
			&dado.CodigoContaContabil,
			&dado.NomeContaContabil,
			&dado.SaldoAnterior,
			&dado.ValorCredito,
			&dado.ValorDebito,
		); err != nil {
			log.Printf("RelatorioFIP215Handler: %v", err)
			return ErroConsultaLinhaBancoDados
		}

		relatorio.Dados = append(relatorio.Dados, dado)
	}

	if err := rows.Err(); err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroRedeOuResultadoBancoDados
	}
	/*** Consulta no Banco de Dados ***/

	return c.JSON(http.StatusOK, relatorio)
}
