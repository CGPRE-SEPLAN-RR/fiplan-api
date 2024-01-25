package model

type DadoRelatorioFIP215 struct {
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

type RelatorioFIP215 struct {
	Dados []DadoRelatorioFIP215
} // @name RelatorioFIP215
