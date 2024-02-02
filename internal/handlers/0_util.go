package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

/*** Validação ***/
var Validate = validator.New()

/*** Validação ***/

/*** Erro ***/
type Erro echo.HTTPError // @name Erro

var ErroMontagemTemplate *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao montar o template.")
var ErroExecucaoTemplate *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao executar o template.")
var ErroConsultaBancoDados *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao consultar o banco de dados.")
var ErroConsultaLinhaBancoDados *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao consultar uma linha no banco de dados.")
var ErroRedeOuResultadoBancoDados *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro de rede ou problema no resultado do banco de dados.")

func ErroValidacaoParametro(mensagem []string) *echo.HTTPError {
	return echo.NewHTTPError(
		http.StatusBadRequest,
		map[string][]string{
			"erros": mensagem,
		},
	)
}

/*** Erro ***/

/*** Dados Estáticos ***/
var MesParaNome map[int]string = map[int]string{
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
