package main

import (
	"log"

	apiMoedas "github.com/marcusadriano/mcp-moedas/pkg/moedas"
)

func main() {
	moedas, err := apiMoedas.Disponiveis()
	if err != nil {
		log.Fatalf("Error fetching available currencies: %v", err)
	}

	log.Printf("Total available currencies: %d", len(moedas.Values))

	for _, moeda := range moedas.Values {
		log.Printf("Currency: %s, Symbol: %s, Type: %s", moeda.NomeFormatado, moeda.Simbolo, moeda.TipoMoeda)
	}

	cotacao, err := apiMoedas.ConsultarPorSiglaUltimaData("USD")
	if err != nil {
		log.Fatalf("Error fetching currency by symbol: %v", err)
	}
	log.Printf("%+v", cotacao)
}
