package moedas

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	urlBase         = "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata"
	cotacaoMoedaDia = "/CotacaoMoedaDia(moeda=@moeda,dataCotacao=@dataCotacao)"
	moedas          = "/Moedas"
	formatoJson     = "json"
)

type RespostaTipoMoeda struct {
	Values []TipoMoeda `json:"value"`
}

type TipoMoeda struct {
	Simbolo       string `json:"simbolo"`
	NomeFormatado string `json:"nomeFormatado"`
	TipoMoeda     string `json:"tipoMoeda"`
}

type TipoCotacaoMoeda struct {
	ParidadeCompra  float64 `json:"paridadeCompra"`
	ParidadeVenda   float64 `json:"paridadeVenda"`
	CotacaoCompra   float64 `json:"cotacaoCompra"`
	CotacaoVenda    float64 `json:"cotacaoVenda"`
	DataHoraCotacao string  `json:"dataHoraCotacao"`
	TipoBoletim     string  `json:"tipoBoletim"`
}

type RespostaTipoCotacaoMoeda struct {
	Values []TipoCotacaoMoeda `json:"value"`
}

func Disponiveis(ctx context.Context) (*RespostaTipoMoeda, error) {

	url, _ := url.Parse(urlBase + moedas)
	query := url.Query()
	query.Add("$format", formatoJson)
	query.Add("$top", "50")
	query.Add("$skip", "0")
	query.Add("$orderby", "simbolo")

	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	req.Header.Set("Accept", "application/json;odata.metadata=minimal")

	log.Printf("Sending request to %s...\n", url.String())

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}
	var resposta RespostaTipoMoeda
	if err := json.NewDecoder(resp.Body).Decode(&resposta); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	return &resposta, nil
}

func ConsultarPorSigla(ctx context.Context, sigla string) (*RespostaTipoCotacaoMoeda, error) {
	data := time.Now()
	return consultarPorSigla(ctx, sigla, data)
}

func ConsultarPorSiglaEData(ctx context.Context, sigla string, data time.Time) (*RespostaTipoCotacaoMoeda, error) {
	return consultarPorSigla(ctx, sigla, data)
}

func ConsultarPorSiglaUltimaData(ctx context.Context, sigla string) (*RespostaTipoCotacaoMoeda, error) {
	data := time.Now()
	maxTentativas := 5
	i := 0
	for i < maxTentativas {
		if data.Weekday() == time.Saturday || data.Weekday() == time.Sunday {
			data = data.Add(-24 * time.Hour)
			continue
		}
		i++

		res, err := consultarPorSigla(ctx, sigla, data)
		if err != nil {
			return nil, err
		}
		if len(res.Values) > 0 {
			return res, nil
		}
		data = data.Add(-24 * time.Hour)
	}
	return nil, fmt.Errorf("no data found for the last 5 days (excluding weekends)")
}

func consultarPorSigla(ctx context.Context, sigla string, data time.Time) (*RespostaTipoCotacaoMoeda, error) {

	dataFormatada := data.Format("01-02-2006")

	url, _ := url.Parse(urlBase + cotacaoMoedaDia)
	query := url.Query()
	query.Add("$format", formatoJson)
	query.Add("@moeda", fmt.Sprintf("'%s'", sigla))
	query.Add("@dataCotacao", fmt.Sprintf("'%s'", dataFormatada))

	url.RawQuery = query.Encode()
	bcUrl := url.String()
	log.Printf("Sending request to %s...\n", bcUrl)

	req, err := http.NewRequestWithContext(ctx, "GET", bcUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}
	var resposta RespostaTipoCotacaoMoeda
	if err := json.NewDecoder(resp.Body).Decode(&resposta); err != nil {
		log.Fatalf("Error decoding response: %v", err)
		return nil, err
	}
	return &resposta, nil
}
