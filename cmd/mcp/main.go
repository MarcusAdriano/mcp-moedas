package main

import (
	"github.com/marcusadriano/mcp-moedas/pkg/tools"
	"github.com/mark3labs/mcp-go/server"
)

const (
	appDescription = "Cotação de Moedas no Bacen PTAX"
	appVersion     = "0.0.1"
)

func main() {
	s := server.NewMCPServer(
		appDescription,
		appVersion,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddTool(tools.CotacaoMoedasTool, tools.CotacaoMoedasHandler)
}
