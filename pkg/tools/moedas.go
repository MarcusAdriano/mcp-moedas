package tools

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/marcusadriano/mcp-moedas/pkg/moedas"
	"github.com/mark3labs/mcp-go/mcp"
)

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

	simbolo, ok := request.Params.Arguments["simbolo"].(string)
	if !ok {
		return nil, errors.New("simbolo parameter is required")
	}

	data, ok := request.Params.Arguments["data"].(string)
	if !ok {
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
	cotacao, err := moedas.ConsultarPorSiglaUltimaData(simbolo)
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

	cotacao, err := moedas.ConsultarPorSiglaEData(simbolo, data)
	if err != nil {
		return nil, err
	}

	if len(cotacao.Values) == 0 {
		return nil, fmt.Errorf("não foi possível encontrar a cotação para o símbolo %s na data %s", simbolo, dataStr)
	}

	return resultadoFinal(cotacao), nil
}
