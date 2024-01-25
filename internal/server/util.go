package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/*** Erro ***/
type Erro echo.HTTPError // @name Erro

var ErroConsultaBancoDados *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao consultar o banco de dados.")
var ErroConsultaLinhaBancoDados *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao consultar uma linha no banco de dados.")
var ErroRedeOuResultadoBancoDados *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro de rede ou problema no resultado do banco de dados.")

/*** Erro ***/

/*** Data ***/
const (
	JANEIRO = iota + 1
	FEVEREIRO
	MARCO
	ABRIL
	MAIO
	JUNHO
	JULHO
	AGOSTO
	SETEMBRO
	OUTUBRO
	NOVEMBRO
	DEZEMBRO
)

var DataTexto map[int]string = map[int]string{
	JANEIRO:   "JANEIRO",
	FEVEREIRO: "FEVEREIRO",
	MARCO:     "MARCO",
	ABRIL:     "ABRIL",
	MAIO:      "MAIO",
	JUNHO:     "JUNHO",
	JULHO:     "JULHO",
	AGOSTO:    "AGOSTO",
	SETEMBRO:  "SETEMBRO",
	OUTUBRO:   "OUTUBRO",
	NOVEMBRO:  "NOVEMBRO",
	DEZEMBRO:  "DEZEMBRO",
}
/*** Data ***/
