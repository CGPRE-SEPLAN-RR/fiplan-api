package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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
