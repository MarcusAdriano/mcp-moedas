package tools

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
)

var CotacaoMoedasTool = mcp.NewTool("cotacao_moedas",
	mcp.WithDescription("Cotação de Moedas no Bacen PTAX"),
	mcp.WithString("moeda",
		mcp.Required(),
		mcp.Description("Código da moeda a ser consultada, conforme tabela de códigos de moedas do Bacen. Exemplos: USD, EUR, CAD, etc"),
	),
	mcp.WithString("data",
		mcp.Description("Data desejada da cotação no formato MM-DD-YYYY. Se não for informada, será considerada a data atual."),
	),
)

func CotacaoMoedasHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return nil, errors.New("CotacaoMoedasHandler not implemented yet")
}
