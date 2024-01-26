package test

import (
	"encoding/json"
	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/server"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRelatorioFIP215Handler(t *testing.T) {
	e := echo.New()
	s := &server.Server{}

	t.Run("valida ano de exercício com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida ano de exercício com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida ano de exercício menor que 2010", TestRelatorioFIP215Handler)
	t.Run("valida ano de exercício igual a 2010", TestRelatorioFIP215Handler)
	t.Run("valida ano de exercício igual a 2015", TestRelatorioFIP215Handler)
	t.Run("valida ano de exercício igual ao ano atual", TestRelatorioFIP215Handler)
	t.Run("valida ano de exercício maior que o ano atual", TestRelatorioFIP215Handler)

	t.Run("valida mês de referência com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida mês de referência com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida mês de referência menor que 1", TestRelatorioFIP215Handler)
	t.Run("valida mês de referência igual a 1", TestRelatorioFIP215Handler)
	t.Run("valida mês de referência igual a 6", TestRelatorioFIP215Handler)
	t.Run("valida mês de referência igual a 12", TestRelatorioFIP215Handler)
	t.Run("valida mês de referência maior que 12", TestRelatorioFIP215Handler)

	t.Run("valida unidade gestora com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida unidade gestora com tipo não inteiro", TestRelatorioFIP215Handler)

	t.Run("valida unidade orçamentária com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida unidade orçamentária com tipo não inteiro", TestRelatorioFIP215Handler)

	t.Run("valida mês contábil com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida mês contábil com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida mês contábil menor que 1", TestRelatorioFIP215Handler)
	t.Run("valida mês contábil igual a 1", TestRelatorioFIP215Handler)
	t.Run("valida mês contábil igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida mês contábil igual a 4", TestRelatorioFIP215Handler)
	t.Run("valida mês contábil maior que 4", TestRelatorioFIP215Handler)

	t.Run("valida tipo de poder com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida tipo de poder com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida tipo de poder menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida tipo de poder igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida tipo de poder igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida tipo de poder igual a 5", TestRelatorioFIP215Handler)
	t.Run("valida tipo de poder maior que 5", TestRelatorioFIP215Handler)

	t.Run("valida tipo de administração com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida tipo de administração com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida tipo de administração menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida tipo de administração igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida tipo de administração igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida tipo de administração igual a 3", TestRelatorioFIP215Handler)
	t.Run("valida tipo de administração maior que 3", TestRelatorioFIP215Handler)

	t.Run("valida tipo de encerramento com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida tipo de encerramento com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida tipo de encerramento menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida tipo de encerramento igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida tipo de encerramento igual a 1", TestRelatorioFIP215Handler)
	t.Run("valida tipo de encerramento igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida tipo de encerramento maior que 2", TestRelatorioFIP215Handler)

	t.Run("valida indicativo de superávit financeiro com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de superávit financeiro com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de superávit financeiro menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de superávit financeiro igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de superávit financeiro igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de superávit financeiro maior que 2", TestRelatorioFIP215Handler)

	t.Run("valida indicativo de conta contábil com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de conta contábil com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de conta contábil menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de conta contábil igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de conta contábil igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de conta contábil maior que 2", TestRelatorioFIP215Handler)

	t.Run("valida indicativo de composição da MSC com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de composição da MSC com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de composição da MSC menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de composição da MSC igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de composição da MSC igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida indicativo de composição da MSC maior que 2", TestRelatorioFIP215Handler)

	t.Run("valida consolidado RPPS com tipo inteiro", TestRelatorioFIP215Handler)
	t.Run("valida consolidado RPPS com tipo não inteiro", TestRelatorioFIP215Handler)
	t.Run("valida consolidado RPPS menor que 0", TestRelatorioFIP215Handler)
	t.Run("valida consolidado RPPS igual a 0", TestRelatorioFIP215Handler)
	t.Run("valida consolidado RPPS igual a 2", TestRelatorioFIP215Handler)
	t.Run("valida consolidado RPPS maior que 2", TestRelatorioFIP215Handler)

	req := httptest.NewRequest(http.MethodGet, "/relatorio/fip_125", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	if err := s.RelatorioFIP215Handler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}

	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}

	expected := map[string][]string{"erros": {"a"}}
	var actual map[string][]string

	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Errorf("handler() error decoding response body: %v", err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("handler() wrong response body. expected = %v, actual = %v", expected, actual)
		return
	}
}
