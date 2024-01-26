package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
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
	var unidadeOrcamentaria uint16
	var mesReferencia uint8
	var mesContabil uint8
	var tipoPoder uint8
	var tipoAdministracao uint8
	var tipoEncerramento uint8
	var consolidadoRPPS uint8
	var indicativoComposicaoMSC uint8
	var indicativoContaContabilRP uint8
	var indicativoSuperavitFinanceiro uint8

	var parametros struct {
		anoExercicio                  string
		unidadeGestora                string
		unidadeOrcamentaria           string
		mesReferencia                 string
		mesContabil                   string
		tipoPoder                     string
		tipoAdministracao             string
		tipoEncerramento              string
		consolidadoRPPS               string
		indicativoComposicaoMSC       string
		indicativoContaContabilRP     string
		indicativoSuperavitFinanceiro string
	}
	/*** Parâmetros ***/

	/*** Dados Estáticos ***/
	mesContabilParaFlag := map[uint8]string{
		1: "1428",
		2: "1429",
		3: "1430",
		4: "10120",
	}

	indicativoSuperavitFinanceiroParaFlag := map[uint8]string{
		1: "856",
		2: "857",
	}

	indicativoComposicaoMSCParaFlag := map[uint8]string{
		1: "856",
		2: "857",
	}

	mesParaTexto := map[uint8]string{
		1:  "JANEIRO",
		2:  "FEVEREIRO",
		3:  "MARCO",
		4:  "ABRIL",
		5:  "MAIO",
		6:  "JUNHO",
		7:  "JULHO",
		8:  "AGOSTO",
		9:  "SETEMBRO",
		10: "OUTUBRO",
		11: "NOVEMBRO",
		12: "DEZEMBRO",
	}
	/*** Dados Estáticos ***/

	/*** Validação dos Parâmetros ***/
	valueBinder := echo.QueryParamsBinder(c)

	var errors []string
	
	if err := valueBinder.MustUint16("ano_exercicio", &anoExercicio).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o ano de exercício no parâmetro 'ano_exercicio'.")
	}

	if err := Validate.Var(anoExercicio, fmt.Sprintf("gte=2010,lte=%d", time.Now().Year())); err != nil {
		errors =append(errors, fmt.Sprintf("Por favor, forneça um ano de exercício válido entre 2010 e %d para o parâmetro 'ano_exercicio'.", time.Now().Year()))
	} else {
		parametros.anoExercicio = fmt.Sprint(anoExercicio)
	}

	if err := valueBinder.Uint16("unidade_gestora", &unidadeGestora).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça a unidade gestora no parâmetro 'unidade_gestora'.")
	} else {
		parametros.unidadeGestora = fmt.Sprint(unidadeGestora)
	}

	if err := valueBinder.Uint16("unidade_orcamentaria", &unidadeOrcamentaria).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça a unidade orçamentária no parâmetro 'unidade_orcamentaria'.")
	} else {
		parametros.unidadeOrcamentaria = fmt.Sprint(unidadeOrcamentaria)
	}

	if err := valueBinder.MustUint8("mes_referencia", &mesReferencia).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o mês de referência no parâmetro 'mes_referencia'.")
	}

	if err := Validate.Var(mesReferencia, "gte=1,lte=12"); err != nil {
		errors = append(errors, "Por favor, forneça um mês de referência válido entre 1 e 12 para o parâmetro 'mes_referencia'.")
	} else {
		parametros.mesReferencia = mesParaTexto[mesReferencia]
	}

	if err := valueBinder.MustUint8("mes_contabil", &mesContabil).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o mês contábil no parâmetro 'mes_contabil'.")
	}

	if err := Validate.Var(mesContabil, "gte=1,lte=4"); err != nil {
		errors = append(errors, "Por favor, forneça um mês contábil válido entre 1 e 4 para o parâmetro 'mes_contabil'.")
	} else {
		parametros.mesContabil = mesContabilParaFlag[mesContabil]
	}

	if err := valueBinder.Uint8("tipo_poder", &tipoPoder).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o tipo de poder no parâmetro 'tipo_poder'.")
	}

	if err := Validate.Var(tipoPoder, "gte=0,lte=5"); err != nil {
		errors = append(errors, "Por favor, forneça um tipo de poder válido entre 1 e 5 para o parâmetro 'tipo_poder'.")
	} else {
		if tipoPoder != 0 {
			parametros.tipoPoder = string(tipoPoder)
		} else {
			parametros.tipoPoder = ""
		}
	}

	if err := valueBinder.Uint8("tipo_administracao", &tipoAdministracao).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o tipo de administração no parâmetro 'tipo_administracao'.")
	}

	if err := Validate.Var(tipoAdministracao, "gte=0,lte=3"); err != nil {
		errors = append(errors, "Por favor, forneça um tipo de administração válido entre 1 e 3 para o parâmetro 'tipo_administracao'.")
	} else {
		if tipoAdministracao != 0 {
			parametros.tipoAdministracao = string(tipoAdministracao)
		} else {
			parametros.tipoAdministracao = ""
		}
	}

	if err := valueBinder.Uint8("tipo_encerramento", &tipoEncerramento).BindError(); err != nil {
		errors = append(errors,"Por favor, forneça o tipo de encerramento no parâmetro 'tipo_encerramento'.")
	}

	if err := Validate.Var(tipoEncerramento, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um tipo de encerramento válido entre 1 e 2 para o parâmetro 'tipo_encerramento'.")
	} else {
		if tipoEncerramento != 0 {
			parametros.tipoEncerramento = string(tipoEncerramento)
		} else {
			parametros.tipoEncerramento = ""
		}
	}

	if err := valueBinder.Uint8("consolidado_rpps", &consolidadoRPPS).BindError(); err != nil {
		errors = append(errors,"Por favor, forneça o consolidado do RPPS no parâmetro 'consolidado_rpps'.")
	}

	if err := Validate.Var(consolidadoRPPS, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um consolidado do RPPS válido entre 1 e 2 para o parâmetro 'consolidado_rpps'.")
	} else {
		if consolidadoRPPS != 0 {
			parametros.consolidadoRPPS = string(consolidadoRPPS)
		} else {
			parametros.consolidadoRPPS = ""
		}
	}

	if err := valueBinder.Uint8("indicativo_conta_contabil_rp", &indicativoContaContabilRP).BindError(); err != nil {
		errors = append(errors,"Por favor, forneça o indicativo de conta contábil de RP no parâmetro 'indicativo_conta_contabil_rp'.")
	}

	if err := Validate.Var(indicativoContaContabilRP, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um indicativo de conta contábil de RP válido entre 1 e 2 para o parâmetro 'indicativo_conta_contabil_rp'.")
	} else {
		if indicativoContaContabilRP != 0 {
			parametros.indicativoContaContabilRP = string(indicativoContaContabilRP)
		} else {
			parametros.indicativoContaContabilRP = ""
		}
	}

	if err := valueBinder.Uint8("indicativo_superavit_fincanceiro", &indicativoSuperavitFinanceiro).BindError(); err != nil {
		errors = append(errors,"Por favor, forneça o indicativo de superávit financeiro no parâmetro 'indicativo_superavit_fincanceiro'.")
	}

	if err := Validate.Var(indicativoSuperavitFinanceiro, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um indicativo de superávit financeiro válido entre 1 e 2 para o parâmetro 'indicativo_superavit_fincanceiro'.")
	} else {
		if indicativoSuperavitFinanceiro != 0 {
			parametros.indicativoSuperavitFinanceiro = indicativoSuperavitFinanceiroParaFlag[indicativoSuperavitFinanceiro]
		} else {
			parametros.indicativoSuperavitFinanceiro = ""
		}
	}

	if err := valueBinder.Uint8("indicativo_composicao_msc", &indicativoComposicaoMSC).BindError(); err != nil {
		errors = append(errors,"Por favor, forneça o indicativo de composição da MSC no parâmetro 'indicativoComposicaoMSC'.")
	}

	if err := Validate.Var(indicativoComposicaoMSC, "gte=0,lte=2"); err != nil {
		errors = append(errors, "Por favor, forneça um indicativo de composição da MSC válido entre 1 e 2 para o parâmetro 'indicativo_composicao_msc'.")
	} else {
		if indicativoComposicaoMSC != 0 {
			parametros.indicativoComposicaoMSC = indicativoComposicaoMSCParaFlag[indicativoComposicaoMSC]
		} else {
			parametros.indicativoComposicaoMSC = ""
		}
	}

	if len(errors) > 0 {
		return ErroValidacaoParametro(errors)
	}
	/*** Validação dos Parâmetros ***/

	/*** Consulta no Banco de Dados ***/
	var relatorio model.RelatorioFIP215

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
										 	  CC.IDEN_CONTA_CONTABIL,
										 	  CC.CONTA_EXPLOSAO,
										 	  CC.CODG_CONTA_CONTABIL,
										 	  CC.NOME_CONTA_CONTABIL,
										 	  SUM(NVL(SM.VALR_CRE_DEZEMBRO, 0)) AS VALOR_CREDITO,
										 	  SUM(NVL(SM.VALR_DEB_DEZEMBRO, 0)) AS VALOR_DEBITO,
										 	  SUM(NVL(SM.VALR_NOVEMBRO, 0)) AS SALDO_ANTERIOR
										 	  
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
										 	  
										 	  WHERE SM.FLAG_MES_CONTABIL != 10120
										 	  
										 	  GROUP BY
										 	  UO.CD_EXERCICIO,
										 	  UO.CD_UNIDADE_ORCAMENTARIA,
										 	  UO.DS_UNIDADE_ORCAMENTARIA,
										 	  UO.ID_UNIDADE_ORCAMENTARIA,
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

											{{if .mesContabil}}
											AND SM.FLAG_MES_CONTABIL = {{.mesContabil}}
											{{end}}

											{{if .indicativoContaContabilRP}}
											AND DOMINIO_RP.CD_ITEM_DOMINIO = {{.indicativoContaContabilRP}}
											{{end}}

											{{if .indicativoSuperavitFinanceiro}}
											AND CC.FLAG_SUPERAVIT_FINANCEIRO = {{.indicativoSuperavitFinanceiro}}
											{{end}}

											{{if .indicativoComposicaoMSC}}
											AND CC.FLAG_COMPOSICAO_MSC = {{.indicativoComposicaoMSC}}
											{{end}}
										  
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

	rows, err := s.db.Query(sqlQuery.String())

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
