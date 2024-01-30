package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/godror/godror"
	"github.com/labstack/echo/v4"
)

type dado struct {
	CodigoUnidadeOrcamentaria string  `json:"codigo_unidade_orcamentaria"`
	NomeUnidadeOrcamentaria   string  `json:"nome_unidade_orcamentaria"`
	IDContaContabil           string  `json:"id_conta_contabil"`
	IDContaContabilExplosao   string  `json:"id_conta_contabil_explosao"`
	CodigoContaContabil       string  `json:"codigo_conta_contabil"`
	NomeContaContabil         string  `json:"nome_conta_contabil"`
	SaldoAnterior             float64 `json:"saldo_anterior"`
	ValorCredito              float64 `json:"valor_credito"`
	ValorDebito               float64 `json:"valor_debito"`
}

type relatorioFIP215 struct {
	Dados []dado
} // @name RelatorioFIP215

// @RelatorioFIP215Handler godoc
// @Summary     FIP215 - Balancete Mensal de Verificação
// @Description Teste
// @Tags        Relatório
// @Accept      json
// @Produce     json
// @Param       ano_exercicio                    query    uint16 true  "Ano de Exercício"
// @Param       unidade_gestora                  query    uint16 false "Código da Unidade Gestora"
// @Param       unidade_orcamentaria             query    uint16 false "Código da Unidade Orçamentária"
// @Param       mes_referencia                   query    uint8  true  "Mês de Referência"                                                                              Enums(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
// @Param       mes_contabil                     query    uint8  true  "Mês Contábil (1-Execução / 2-Apuração / 3-Encerramento / 4-Todos)"                              Enums(1, 2, 3, 4)
// @Param       tipo_poder                       query    uint8  false "Tipo de Poder (1-Executivo / 2-Legislativo / 3-Judiciário / 4-Ministério Público / 5-Todos)"    Enums(1, 2, 3, 4, 5)
// @Param       tipo_administracao               query    uint8  false "Tipo de Administração (1-Diretas / 2-Indiretas / 3-Todas)"                                      Enums(1, 2, 3)
// @Param       tipo_encerramento                query    uint8  false "Tipo de Encerramento (1-Encerra ao Final do Exercício / 2-Transfere para o Exercício Seguinte)" Enums(1, 2)
// @Param       consolidado_rpps                 query    bool   false "Consolidado RPPS"                                                                               Enums(true, false)
// @Param       indicativo_conta_contabil_rp     query    uint8  false "Indicativo de Conta Contábil de RP"                                                             Enums(1, 2)
// @Param       indicativo_superavit_fincanceiro query    uint8  false "Indicativo de Superávit Financeiro"                                                             Enums(1, 2)
// @Param       indicativo_composicao_msc        query    uint8  false "Indicativo de Composição da MSC (1-Sim / 2-Não)"                                                Enums(1, 2)
// @Success     200                              {object} relatorioFIP215
// @Failure     400                              {object} Erro
// @Failure     404                              {object} Erro
// @Failure     500                              {object} Erro
// @Router      /relatorio/fip_215 [get]
func (s *Server) RelatorioFIP215Handler(c echo.Context) error {
	/*** Parâmetros ***/
	parametros := struct {
		AnoExercicio                  uint16 // Unused
		UnidadeGestora                uint16 // Unused
		UnidadeOrcamentaria           uint16 // Unused
		MesReferencia                 uint8
		MesContabil                   uint8
		TipoPoder                     uint8
		TipoAdministracao             uint8 // Unused
		TipoEncerramento              uint8
		ConsolidadoRPPS               bool
		IndicativoComposicaoMSC       uint8
		IndicativoContaContabilRP     uint8
		IndicativoSuperavitFinanceiro uint8

		Meses                     []string
		MesReferenciaNome         string
		MesAnteriorReferenciaNome string
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

	if err := valueBinder.Uint16("unidade_gestora", &parametros.UnidadeGestora).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça a unidade gestora no parâmetro 'unidade_gestora'.")
	}

	if err := valueBinder.Uint16("unidade_orcamentaria", &parametros.UnidadeOrcamentaria).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça a unidade orçamentária no parâmetro 'unidade_orcamentaria'.")
	}

	if err := valueBinder.MustUint8("mes_referencia", &parametros.MesReferencia).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o mês de referência no parâmetro 'mes_referencia'.")
	}

	if err := Validate.Var(parametros.MesReferencia, "gte=1,lte=12"); err != nil {
		errors = append(errors, "Por favor, forneça um mês de referência válido entre 1 e 12 para o parâmetro 'mes_referencia'.")
	} else {
		var i uint8
		for i = 1; i < parametros.MesReferencia; i++ {
			parametros.Meses = append(parametros.Meses, MesParaNome[i])
		}

		parametros.MesReferenciaNome = MesParaNome[parametros.MesReferencia]
		parametros.MesAnteriorReferenciaNome = MesParaNome[parametros.MesReferencia-1]
	}

	if err := valueBinder.MustUint8("mes_contabil", &parametros.MesContabil).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o mês contábil no parâmetro 'mes_contabil'.")
	}

	if err := Validate.Var(parametros.MesContabil, "gte=1,lte=4"); err != nil {
		errors = append(errors, "Por favor, forneça um mês contábil válido entre 1 e 4 para o parâmetro 'mes_contabil'.")
	}

	if err := valueBinder.Uint8("tipo_poder", &parametros.TipoPoder).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o tipo de poder no parâmetro 'tipo_poder'.")
	}

	if err := Validate.Var(parametros.TipoPoder, "gte=0,lte=5"); err != nil {
		errors = append(errors, "Por favor, forneça um tipo de poder válido entre 1 e 5 para o parâmetro 'tipo_poder'.")
	}

	if err := valueBinder.Uint8("tipo_administracao", &parametros.TipoAdministracao).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o tipo de administração no parâmetro 'tipo_administracao'.")
	}

	if err := Validate.Var(parametros.TipoAdministracao, "gte=0,lte=3"); err != nil {
		errors = append(errors, "Por favor, forneça um tipo de administração válido entre 1 e 3 para o parâmetro 'tipo_administracao'.")
	}

	if err := valueBinder.Uint8("tipo_encerramento", &parametros.TipoEncerramento).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o tipo de encerramento no parâmetro 'tipo_encerramento'.")
	}

	if err := Validate.Var(parametros.TipoEncerramento, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um tipo de encerramento válido entre 1 e 2 para o parâmetro 'tipo_encerramento'.")
	}

	if err := valueBinder.Bool("consolidado_rpps", &parametros.ConsolidadoRPPS).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o consolidado do RPPS no parâmetro 'consolidado_rpps'.")
	}

	if err := valueBinder.Uint8("indicativo_conta_contabil_rp", &parametros.IndicativoContaContabilRP).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o indicativo de conta contábil de RP no parâmetro 'indicativo_conta_contabil_rp'.")
	}

	if err := Validate.Var(parametros.IndicativoContaContabilRP, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um indicativo de conta contábil de RP válido entre 1 e 2 para o parâmetro 'indicativo_conta_contabil_rp'.")
	}

	if err := valueBinder.Uint8("indicativo_superavit_fincanceiro", &parametros.IndicativoSuperavitFinanceiro).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o indicativo de superávit financeiro no parâmetro 'indicativo_superavit_fincanceiro'.")
	}

	if err := Validate.Var(parametros.IndicativoSuperavitFinanceiro, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um indicativo de superávit financeiro válido entre 1 e 2 para o parâmetro 'indicativo_superavit_fincanceiro'.")
	}

	if err := valueBinder.Uint8("indicativo_composicao_msc", &parametros.IndicativoComposicaoMSC).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o indicativo de composição da MSC no parâmetro 'indicativoComposicaoMSC'.")
	}

	if err := Validate.Var(parametros.IndicativoComposicaoMSC, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um indicativo de composição da MSC válido entre 1 e 2 para o parâmetro 'indicativo_composicao_msc'.")
	}

	if len(errors) > 0 {
		return ErroValidacaoParametro(errors)
	}
	/*** Validação dos Parâmetros ***/

	/*** Consulta no Banco de Dados ***/
	var relatorio relatorioFIP215

	queryTemplate := `SELECT
										RESULTADO_SALDO_INICIAL.CD_UNIDADE_ORCAMENTARIA,
										RESULTADO_SALDO_INICIAL.DS_UNIDADE_ORCAMENTARIA,
										RESULTADO_SALDO_INICIAL.IDEN_CONTA_CONTABIL,
										RESULTADO_SALDO_INICIAL.CONTA_EXPLOSAO,
										RESULTADO_SALDO_INICIAL.CODG_CONTA_CONTABIL,
										RESULTADO_SALDO_INICIAL.NOME_CONTA_CONTABIL,
										RESULTADO_SALDO_INICIAL.SALDO_ANTERIOR,
										RESULTADO_SALDO_INICIAL.VALOR_CREDITO,
										RESULTADO_SALDO_INICIAL.VALOR_DEBITO

										FROM
										(
										  SELECT
										  RESULTADO_SALDO_MENSAL.CD_EXERCICIO,
										  RESULTADO_SALDO_MENSAL.ID_UNIDADE_ORCAMENTARIA,
										  RESULTADO_SALDO_MENSAL.CD_UNIDADE_ORCAMENTARIA,
										  RESULTADO_SALDO_MENSAL.DS_UNIDADE_ORCAMENTARIA,
										  RESULTADO_SALDO_MENSAL.IDEN_CONTA_CONTABIL,
										  RESULTADO_SALDO_MENSAL.CONTA_EXPLOSAO,
										  RESULTADO_SALDO_MENSAL.CODG_CONTA_CONTABIL,
										  RESULTADO_SALDO_MENSAL.NOME_CONTA_CONTABIL,
										  NVL(RESULTADO_SALDO_MENSAL.VALOR_CREDITO, 0) AS VALOR_CREDITO,
										  NVL(RESULTADO_SALDO_MENSAL.VALOR_DEBITO, 0) AS VALOR_DEBITO,
										  COALESCE(SUM(SI.SALDO_ABERTURA), 0) - RESULTADO_SALDO_MENSAL.SALDO_ANTERIOR AS SALDO_ANTERIOR
										  
										  FROM
										  (
										 	  SELECT
										 	  UO.CD_EXERCICIO,
										 	  UO.ID_UNIDADE_ORCAMENTARIA,
										 	  UO.CD_UNIDADE_ORCAMENTARIA,
										 	  UO.DS_UNIDADE_ORCAMENTARIA,
												
												{{if .TipoPoder}}
												ORG.FLG_TP_PODER,
												{{end}}

										 	  CC.IDEN_CONTA_CONTABIL,
										 	  CC.CONTA_EXPLOSAO,
										 	  CC.CODG_CONTA_CONTABIL,
										 	  CC.NOME_CONTA_CONTABIL,

												{{range .Meses}}
										 	  SUM(NVL(SM.VALR_CRE_{{.}}, 0)) +
												{{end}}
												
										 	  SUM(NVL(SM.VALR_CRE_{{.MesReferenciaNome}}, 0)) AS VALOR_CREDITO,

												{{range .Meses}}
										 	  SUM(NVL(SM.VALR_DEB_{{.}}, 0)) +
												{{end}}

										 	  SUM(NVL(SM.VALR_DEB_{{.MesReferenciaNome}}, 0)) AS VALOR_DEBITO,

												{{if eq .MesReferencia 1}}
										 	  SUM(0) AS SALDO_ANTERIOR
												{{else}}
										 	  SUM(NVL(SM.VALR_{{.MesAnteriorReferenciaNome}}, 0)) AS SALDO_ANTERIOR
												{{end}}
										 	  
										 	  FROM
										 	  UNIDADE_ORCAMENTARIA UO 
										 	  LEFT JOIN
										 	  UNIDADE_GESTORA UG ON (UG.ID_UNIDADE_ORCAMENTARIA = UO.ID_UNIDADE_ORCAMENTARIA)
										 	  LEFT JOIN
										 	  ORGAO ORG ON ORG.ID_ORGAO = UO.ID_ORGAO 
										 	  LEFT JOIN
										 	  ACWTB0215 SM ON (SM.ID_UNIDADE_GESTORA = UG.ID_UNIDADE_GESTORA)
										 	  LEFT JOIN
										 	  ACWTB0032 CC ON (CC.IDEN_CONTA_CONTABIL = SM.IDEN_CONTA_CONTABIL)
										 	  LEFT JOIN
										 	  ITEM_DOMINIO DOMINIO_RP ON (CC.FLAG_CONTA_CONTABIL_RP = DOMINIO_RP.ID_ITEM_DOMINIO)
										 	  LEFT JOIN
										 	  ITEM_DOMINIO DOMINIO_ENCERRA ON (CC.FLAG_TIPO_ENCERRAMENTO = DOMINIO_ENCERRA.ID_ITEM_DOMINIO)

									      WHERE UO.CD_EXERCICIO = {{.AnoExercicio}}

												{{if .UnidadeGestora}}
												{{end}}
												
												{{if .UnidadeOrcamentaria}}
												{{end}}

												{{if .TipoAdministracao}}
												{{end}}

												{{if eq .MesContabil 1}}
										 	  AND SM.FLAG_MES_CONTABIL = 1428
												{{else if eq .MesContabil 2}}
										 	  AND SM.FLAG_MES_CONTABIL = 1429
												{{else if eq .MesContabil 3}}
										 	  AND SM.FLAG_MES_CONTABIL = 1430
												{{else if eq .MesContabil 4}}
										 	  AND SM.FLAG_MES_CONTABIL != 10120
												{{end}}

												{{if eq .TipoPoder 1}}
										 	  AND ORG.FLG_TP_PODER = {{.TipoPoder}}
										 	  AND UO.CD_UNIDADE_ORCAMENTARIA <> 08101
										 	  AND UO.CD_UNIDADE_ORCAMENTARIA <> 08601
												{{else if or (eq .TipoPoder 2) (eq .TipoPoder 3)}}
										 	  AND ORG.FLG_TP_PODER = {{.TipoPoder}}
												{{else if eq .TipoPoder 4}}
										 	  AND (UO.CD_UNIDADE_ORCAMENTARIA <> 08101 OR UO.CD_UNIDADE_ORCAMENTARIA <> 08601)
												{{end}}

												{{if .IndicativoContaContabilRP}}
												AND DOMINIO_RP.CD_ITEM_DOMINIO = {{.IndicativoContaContabilRP}}
												{{end}}

												{{if .TipoEncerramento}}
												AND DOMINIO_ENCERRA.CD_ITEM_DOMINIO = {{.TipoEncerramento}}
												{{end}}

												{{if eq .IndicativoSuperavitFinanceiro 1}}
												AND CC.FLAG_SUPERAVIT_FINANCEIRO = 856
												{{else if eq .IndicativoSuperavitFinanceiro 2}}
												AND CC.FLAG_SUPERAVIT_FINANCEIRO = 857
												{{end}}

												{{if eq .IndicativoComposicaoMSC 1}}
												AND CC.FLAG_COMPOSICAO_MSC = 856
												{{else if eq .IndicativoComposicaoMSC 2}}
												AND CC.FLAG_COMPOSICAO_MSC = 857
												{{end}}

												{{if .ConsolidadoRPPS}}
												AND UO.CD_UNIDADE_ORCAMENTARIA IN (15301, 15601, 15602, 15603)
												{{end}}
										 	  
										 	  GROUP BY
										 	  UO.CD_EXERCICIO,
										 	  UO.CD_UNIDADE_ORCAMENTARIA,
										 	  UO.DS_UNIDADE_ORCAMENTARIA,
										 	  UO.ID_UNIDADE_ORCAMENTARIA,

												{{if .TipoPoder}}
												ORG.FLG_TP_PODER,
												{{end}}

										 	  CC.IDEN_CONTA_CONTABIL,
										 	  CC.CONTA_EXPLOSAO,
										 	  CC.CODG_CONTA_CONTABIL,
										 	  CC.NOME_CONTA_CONTABIL
										  ) RESULTADO_SALDO_MENSAL
										  LEFT JOIN
										  UNIDADE_GESTORA UG ON (UG.ID_UNIDADE_ORCAMENTARIA = RESULTADO_SALDO_MENSAL.ID_UNIDADE_ORCAMENTARIA)
										  LEFT JOIN
										  ACWTB0197 SI ON SI.CD_EXERCICIO = RESULTADO_SALDO_MENSAL.CD_EXERCICIO
										  AND SI.ID_UNIDADE_GESTORA = UG.ID_UNIDADE_GESTORA
										  AND SI.IDEN_CONTA_CONTABIL = RESULTADO_SALDO_MENSAL.IDEN_CONTA_CONTABIL

										  GROUP BY
										  RESULTADO_SALDO_MENSAL.CD_EXERCICIO,
										  RESULTADO_SALDO_MENSAL.ID_UNIDADE_ORCAMENTARIA,
										  RESULTADO_SALDO_MENSAL.CD_UNIDADE_ORCAMENTARIA,
										  RESULTADO_SALDO_MENSAL.DS_UNIDADE_ORCAMENTARIA,
										  RESULTADO_SALDO_MENSAL.IDEN_CONTA_CONTABIL,
										  RESULTADO_SALDO_MENSAL.CONTA_EXPLOSAO,
										  RESULTADO_SALDO_MENSAL.CODG_CONTA_CONTABIL,
										  RESULTADO_SALDO_MENSAL.NOME_CONTA_CONTABIL,
										  RESULTADO_SALDO_MENSAL.SALDO_ANTERIOR,
										  RESULTADO_SALDO_MENSAL.VALOR_CREDITO,
										  RESULTADO_SALDO_MENSAL.VALOR_DEBITO
										) RESULTADO_SALDO_INICIAL
										LEFT JOIN
										UNIDADE_GESTORA UG ON (UG.ID_UNIDADE_ORCAMENTARIA = RESULTADO_SALDO_INICIAL.ID_UNIDADE_ORCAMENTARIA)

										GROUP BY
										RESULTADO_SALDO_INICIAL.CD_UNIDADE_ORCAMENTARIA,
										RESULTADO_SALDO_INICIAL.DS_UNIDADE_ORCAMENTARIA,
										RESULTADO_SALDO_INICIAL.IDEN_CONTA_CONTABIL,
										RESULTADO_SALDO_INICIAL.CONTA_EXPLOSAO,
										RESULTADO_SALDO_INICIAL.CODG_CONTA_CONTABIL,
										RESULTADO_SALDO_INICIAL.NOME_CONTA_CONTABIL,
										RESULTADO_SALDO_INICIAL.SALDO_ANTERIOR,
										RESULTADO_SALDO_INICIAL.VALOR_CREDITO,
										RESULTADO_SALDO_INICIAL.VALOR_DEBITO`

	tmpl, err := template.New("query").Parse(queryTemplate)

	if err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroMontagemTemplate
	}

	var sqlQuery strings.Builder
	err = tmpl.Execute(&sqlQuery, parametros)

	if err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroExecucaoTemplate
	}

	compactSqlQuery := strings.Join(strings.Fields(sqlQuery.String()), " ")
	log.Printf("RelatorioFIP215Handler: %s", compactSqlQuery)
	rows, err := s.db.Query(compactSqlQuery, godror.PrefetchCount(10000), godror.FetchArraySize(10000))

	if err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroConsultaBancoDados
	}

	defer rows.Close()

	for rows.Next() {
		var dado dado

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
