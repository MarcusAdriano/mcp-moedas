package main

import (
	"flag"
	"log"

	tools "github.com/marcusadriano/mcp-moedas/pkg/mcptools"
	"github.com/mark3labs/mcp-go/server"
)

const (
	appDescription = "Cotação de Moedas no Bacen PTAX"
	appVersion     = "0.0.1"
)

var (
	runWithSee = false
)

func main() {

	flag.BoolVar(&runWithSee, "see", false, "Run with see")
	flag.Parse()

	s := server.NewMCPServer(
		appDescription,
		appVersion,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// tools
	s.AddTool(tools.CotacaoMoedasTool, tools.CotacaoMoedasHandler)
	s.AddTool(tools.MoedasDisponiveisTool, tools.MoedasDisponiveisHandler)

	// prompts
	s.AddPrompt(tools.PromptCotacaoMoedas, tools.PromptCotacaoMoedasHandler)

	if runWithSee {
		serveSSE(s)
	}
	serveStdio(s)
}

func serveSSE(s *server.MCPServer) {
	sseServer := server.NewSSEServer(s,
		server.WithSSEEndpoint("/sse"),
	)
	log.Println("Starting MCP server in SSE mode...")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("SSE server error: %v\n", err)
	}
}

func serveStdio(s *server.MCPServer) {
	log.Println("Starting MCP server in STD/IO mode...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
