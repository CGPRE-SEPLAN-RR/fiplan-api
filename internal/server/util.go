package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

/*** Erro ***/
type Erro echo.HTTPError // @name Erro

var ErroLeituraSQL *echo.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, "Ocorreu um erro ao ler o arquivo SQL.")
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

/*** Leitura de Queries ***/
func LerSQL(path string) (string, error) {
	content, err := os.ReadFile(fmt.Sprintf("internal/database/queries/%s", path))

	if err != nil {
		return "", err
	}

	return string(content), nil
}
/*** Leitura de Queries ***/
