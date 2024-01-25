SELECT
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
RESULTADO_SALDO_INICIAL.VALOR_DEBITO
