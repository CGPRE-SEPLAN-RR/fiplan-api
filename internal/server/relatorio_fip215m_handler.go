package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type dadoRelatorioFIP215M struct {
	CodigoContaSICONFI string  `json:"codigo_unidade_orcamentaria"`
	ValorCredito       float64 `json:"valor_credito"`
	ValorDebito        float64 `json:"valor_debito"`
	SaldoAbertura      float64 `json:"saldo_atual"`
} // @name DadoRelatorioFIP215M

type relatorioFIP215M struct {
	Dados []dadoRelatorioFIP215M
} // @name RelatorioFIP215M

// RelatorioFIP215MHandler godoc
//
// @Summary     FIP215M - Emitir Matriz de Saldos Contábeis - MSC SICONFI
// @Description Emite a Matriz de Saldos Contábeis - MSC SICONFI
// @Tags        Relatório
// @Accept      json
// @Produce     json
// @Param       ano_exercicio                    query    uint16 true  "Ano de Exercício"
// @Param       mes_referencia                   query    uint8  true  "Mês de Referência"                                Enums(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
// @Param       mes_contabil                     query    uint8  true  "Mês Contábil (1-Execução / 2-Apuração / 3-Todos)" Enums(1, 2, 3)
// @Param       codigo_poder_orgao               query    uint16 true  "Código do Poder/Órgão SICONFI"
// @Success     200                              {object} relatorioFIP215
// @Failure     400                              {object} Erro
// @Failure     500                              {object} Erro
// @Router      /relatorio/fip_215 [get]
func (s *Server) RelatorioFIP215MHandler(c echo.Context) error {
	/*** Parâmetros ***/
	parametros := struct {
		// FIPLAN
		AnoExercicio     uint16
		MesReferencia    uint8
		MesContabil      uint8
		CodigoPoderOrgao uint32

		// Adicionai
		MesReferenciaNome string
		NomePoderOrgao    string
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

	if err := valueBinder.MustUint8("mes_referencia", &parametros.MesReferencia).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o mês de referência no parâmetro 'mes_referencia'.")
	}

	if err := Validate.Var(parametros.MesReferencia, "gte=1,lte=12"); err != nil {
		errors = append(errors, "Por favor, forneça um mês de referência válido entre 1 e 12 para o parâmetro 'mes_referencia'.")
	} else {
		parametros.MesReferenciaNome = MesParaNome[parametros.MesReferencia]
	}

	if err := valueBinder.MustUint8("mes_contabil", &parametros.MesContabil).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o mês contábil no parâmetro 'mes_contabil'.")
	}

	if err := Validate.Var(parametros.MesContabil, "gte=1,lte=3"); err != nil {
		errors = append(errors, "Por favor, forneça um mês contábil válido entre 1 e 3 para o parâmetro 'mes_contabil'.")
	}

	if err := valueBinder.MustUint32("codigo_poder_orgao", &parametros.CodigoPoderOrgao).BindError(); err != nil {
		errors = append(errors, "Por favor, forneça o código do poder/órgão SICONFI no parâmetro 'codigo_poder_orgao'.")
	}

	if parametros.MesReferencia != 12 && parametros.MesContabil != 1 {
		errors = append(errors, "A opção de emissão do relatório para o mês contábil 2 (apuração) ou 3 (todos) só está disponível para o mês de referência 12 (dezembro).")
	}

	if len(errors) > 0 {
		return ErroValidacaoParametro(errors)
	}
	/*** Validação dos Parâmetros ***/

	/*** Consulta no Banco de Dados ***/
	var sqlQuery strings.Builder

	var msc relatorioFIP215M

	if parametros.CodigoPoderOrgao != 0 {
		queryTemplate := `SELECT UNIQUE NOME_PODER_ORGAO_SICONFI
		
																FROM ACWTB0803

																WHERE CODG_PODER_ORGAO_SICONFI = {{.codigo_poder_orgao}}
																AND CD_EXERCICIO = {{.AnoExercicio}}`

		tmpl, err := template.New("queryMSC").Parse(queryTemplate)

		if err != nil {
			log.Printf("RelatorioFIP215MHandler: %v", err)
			return ErroMontagemTemplate
		}

		err = tmpl.Execute(&sqlQuery, parametros)

		if err != nil {
			log.Printf("RelatorioFIP215MHandler: %v", err)
			return ErroExecucaoTemplate
		}

		compactSqlQuery := strings.Join(strings.Fields(sqlQuery.String()), " ")
		log.Printf("RelatorioFIP215MHandler: %s", compactSqlQuery)
		row := s.db.QueryRow(compactSqlQuery)

		var nomePoderOrgao string

		if err := row.Scan(
			&nomePoderOrgao,
		); err != nil {
			log.Printf("RelatorioFIP215Handler: %v", err)
			return ErroConsultaLinhaBancoDados
		}

		parametros.NomePoderOrgao = fmt.Sprintf("%d - %s", parametros.CodigoPoderOrgao, nomePoderOrgao)

		sqlQuery.Reset()
	} else {
		parametros.NomePoderOrgao = "CONSOLIDADO DO ESTADO"
	}

	queryTemplate := `SELECT CCS.CODG_CONTA_SICONFI,
	                  NVL(SUM(SA.VALR_CRE_{{.MesReferenciaNome}}),0) CREDITO,
										NVL(SUM(SA.VALR_DEB_{{.MesReferenciaNome}}),0) DEBITO,

										{{if eq .MesReferencia 1}}
			              NVL(SUM(SA.SALDO_ABERTURA),0) SALDO_ABERTURA
										{{else}}
			              NVL(SUM(SA.SALDO_ABERTURA),0) + NVL(SUM(SA.VALR_{{.MesReferenciaNome}}),0) SALDO_ABERTURA
										{{end}}
		 
		                FROM
										ACWTB0801 CCS,
										ACWTA8000 SA

                    WHERE SA.CODG_CONTA_SICONFI=CCS.CODG_CONTA_SICONFI
                    AND SA.CD_EXERCICIO=CCS.CD_EXERCICIO
                    AND SA.CD_EXERCICIO={{.AnoExercicio}}

										{{if eq .MesContabil 1}}
                    AND SA.FLAG_MES_CONTABIL=1428
										{{else if eq .MesContabil 2}}
                    AND SA.FLAG_MES_CONTABIL=1429
										{{end}}

										{{if .CodigoPoderOrgao}}
                    AND SA.IC1='{{.NomePoderOrgao}}'
										{{end}}

                    GROUP BY CCS.CODG_CONTA_SICONFI`

	tmplContaExplosao, err := template.New("queryMSC").Parse(queryTemplate)

	if err != nil {
		log.Printf("RelatorioFIP215MHandler: %v", err)
		return ErroMontagemTemplate
	}

	err = tmplContaExplosao.Execute(&sqlQuery, parametros.AnoExercicio)

	if err != nil {
		log.Printf("RelatorioFIP215MHandler: %v", err)
		return ErroExecucaoTemplate
	}

	compactSqlQuery := strings.Join(strings.Fields(sqlQuery.String()), " ")
	log.Printf("RelatorioFIP215MHandler: %s", compactSqlQuery)
	rows, err := s.db.Query(compactSqlQuery)

	sqlQuery.Reset()

	if err != nil {
		log.Printf("RelatorioFIP215MHandler: %v", err)
		return ErroConsultaBancoDados
	}

	defer rows.Close()

	for rows.Next() {
		var dado dadoRelatorioFIP215M

		if err := rows.Scan(
			&dado.CodigoContaSICONFI,
			&dado.ValorCredito,
			&dado.ValorDebito,
			&dado.SaldoAbertura,
		); err != nil {
			log.Printf("RelatorioFIP215Handler: %v", err)
			return ErroConsultaLinhaBancoDados
		}

		msc.Dados = append(msc.Dados, dado)
	}

	if err := rows.Err(); err != nil {
		log.Printf("RelatorioFIP215Handler: %v", err)
		return ErroRedeOuResultadoBancoDados
	}

	/*** Consulta no Banco de Dados ***/

	/*** Lógica Adicional ***/
	/*** Lógica Adicional ***/

	return c.JSON(http.StatusOK, msc)
}
