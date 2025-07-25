package tools

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/marcusadriano/mcp-moedas/pkg/moedas"
	"github.com/mark3labs/mcp-go/mcp"
)

var PromptCotacaoMoedas = mcp.NewPrompt("instrucao_cotacao_moeda",
	mcp.WithPromptDescription("Prompt para consultar a cotação de moedas no Bacen PTAX"),
)

var PromptCotacaoMoedasHandler = func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {

	moedasDisponiveis, err := moedas.Disponiveis(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar moedas disponíveis: %w", err)
	}
	moedasDisponiveisStr := ""
	for _, moeda := range moedasDisponiveis.Values {
		moedasDisponiveisStr += fmt.Sprintf("%s (%s), ", moeda.Simbolo, moeda.NomeFormatado)
	}

	return mcp.NewGetPromptResult(
		"Instrução para solicitar cotação de uma moeda disponível no PTAX do Bacen",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleAssistant,
				mcp.NewTextContent(fmt.Sprintf("Você é um assistente de cotação de moedas PTAX do Banco Central do Brasil. Essas são as moedas disponíveis: %s", moedasDisponiveisStr)),
			),
			mcp.NewPromptMessage(
				mcp.RoleAssistant,
				mcp.NewTextContent("Qual moeda deseja consultar?"),
			),
		},
	), nil
}

var CotacaoMoedasTool = mcp.NewTool("cotacao_moedas",
	mcp.WithDescription("Cotação de Moedas no Bacen PTAX"),
	mcp.WithString("simbolo",
		mcp.Required(),
		mcp.Description("Código da moeda a ser consultada, conforme tabela de códigos de moedas do Bacen. Exemplos: USD, EUR, CAD, etc"),
	),
	mcp.WithString("data",
		mcp.Description("Data desejada da cotação no formato MM-DD-YYYY. Se não for informada, será considerada a data atual."),
	),
)

func CotacaoMoedasHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	simbolo, err := request.RequireString("simbolo")
	if err != nil {
		return nil, errors.New("simbolo parameter is required")
	}

	data, err := request.RequireString("data")
	if err != nil {
		return cotacaoMaisRecente(ctx, simbolo)
	}

	return cotacaoDataEspecifica(ctx, simbolo, data)
}

func resultadoFinal(cotacao *moedas.RespostaTipoCotacaoMoeda) *mcp.CallToolResult {
	cotacaoMaisRecente := cotacao.Values[len(cotacao.Values)-1]

	resultado := fmt.Sprintf("R$ %.2f em %s", cotacaoMaisRecente.CotacaoCompra, cotacaoMaisRecente.DataHoraCotacao)
	return mcp.NewToolResultText(resultado)
}

func cotacaoMaisRecente(ctx context.Context, simbolo string) (*mcp.CallToolResult, error) {
	cotacao, err := moedas.ConsultarPorSiglaUltimaData(ctx, simbolo)
	if err != nil {
		return nil, err
	}

	if len(cotacao.Values) == 0 {
		return nil, fmt.Errorf("não foi possível encontrar a cotação para o símbolo %s", simbolo)
	}

	return resultadoFinal(cotacao), nil
}

func cotacaoDataEspecifica(ctx context.Context, simbolo string, dataStr string) (*mcp.CallToolResult, error) {

	data, err := time.Parse("01-02-2006", dataStr)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %s", dataStr)
	}

	cotacao, err := moedas.ConsultarPorSiglaEData(ctx, simbolo, data)
	if err != nil {
		return nil, err
	}

	if len(cotacao.Values) == 0 {
		return nil, fmt.Errorf("não foi possível encontrar a cotação para o símbolo %s na data %s", simbolo, dataStr)
	}

	return resultadoFinal(cotacao), nil
}
