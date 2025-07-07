package tools

import (
	"context"

	"github.com/marcusadriano/mcp-moedas/pkg/moedas"
	"github.com/mark3labs/mcp-go/mcp"
)

var MoedasDisponiveisTool = mcp.NewTool("moedas_disponiveis",
	mcp.WithDescription("Lista de moedas disponíveis no Bacen PTAX"),
	mcp.WithNumber("limite",
		mcp.DefaultNumber(10),
		mcp.Description("Número máximo de moedas a serem retornadas. Padrão é 10."),
	),
)

func MoedasDisponiveisHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	moedas, err := moedas.Disponiveis(ctx)
	if err != nil {
		return nil, err
	}

	resultado := "Moedas disponíveis no Bacen PTAX:\n"
	for _, moeda := range moedas.Values {
		resultado += moeda.Simbolo + " - " + moeda.NomeFormatado + "\n"
	}

	return mcp.NewToolResultText(resultado), nil
}
