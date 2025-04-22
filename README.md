### MCP Moedas

Estudando [MCP](https://modelcontextprotocol.io/introduction) na prática utilizando a API de moedas do [Banco Central do Brasil](https://dadosabertos.bcb.gov.br/dataset/dolar-americano-usd-todos-os-boletins-diarios/resource/472f1280-edb1-4cc2-ba36-7bfdd04c1fef?inner_span=True).

### Pré-requisitos

- [Claude Desktop](https://claude.ai/download)
- [Go 1.24+](https://go.dev)
- [Docker](https://www.docker.com/get-started)

### Executando via Golang

1. Clone o repositório:

```bash
git clone https://github.com/MarcusAdriano/mcp-moedas
```

2. Compilar e "instalar" o projeto:

Mac:
```bash
go build -o $(go env GOPATH)/bin/mcp-moedas cmd/mcp/main.go
```

Windows:
```bash
go build -o %GOPATH%\bin\mcp-moedas.exe cmd\mcp\main.go
```

3. Configurando o Claude.ai:

Edite o arquivo JSON com os mcp servers e adicione um novo:

```json
{
    "mcpServers": {
        "cotacao-moedas-server": {
            "command": "{{path-go-bin}}/mcp-moedas",
            "args": []
        }
    }
}
```

Para saber o caminho do binário, execute o comando:

```bash
go env GOPATH
```