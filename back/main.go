package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"
)

const prompt_text = `Você é um serviço de API que converte imagens de gráficos em dados JSON. Sua única função é analisar a imagem e retornar um array JSON válido. Não forneça explicações ou qualquer texto que não seja o JSON.

**Exemplo de Saída Esperada:**
[
  {
    "chart_type": "bar",
    "fake_data": [
      {"name": "Jan", "value": 150},
      {"name": "Fev", "value": 200}
    ]
  },
  {
    "chart_type": "pie",
    "fake_data": [
      {"name": "Produto A", "value": 45},
      {"name": "Produto B", "value": 55}
    ]
  }
]

Analise a imagem fornecida e gere a resposta estritamente nesse formato.`

type ChartDefinition struct {
	ChartType string `json:"chart_type"`
	FakeData  []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"fake_data"`
}

type GeminiDashboardResponse []ChartDefinition

func extractGrapich() {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := os.ReadFile("./teste.png")

	parts := []*genai.Part{
		genai.NewPartFromBytes(bytes, "image/png"),
		genai.NewPartFromText(prompt_text),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	// thinkingBudget := int32(0)

	result, err := client.Models.GenerateContent(

		ctx,
		"gemini-2.5-flash",
		contents,
		nil,
		// &genai.GenerateContentConfig{
		// 	ThinkingConfig: &genai.ThinkingConfig{
		// 		ThinkingBudget: &thinkingBudget, // Disables thinking
		// 	},
		// },
	)
	if err != nil {
		log.Fatal(err)
	}

	rawResponse := result.Text()
	fmt.Println("--- Resposta Bruta da API ---")

	cleanJSON := strings.TrimSpace(rawResponse)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")

	fmt.Println(cleanJSON)

	var dashboardData GeminiDashboardResponse

	if err := json.Unmarshal([]byte(cleanJSON), &dashboardData); err != nil {
		log.Fatalf("Erro ao converter JSON: %v", err)
	}

	fmt.Println("\n--- Dados Convertidos com Sucesso ---")
	for i, chart := range dashboardData {
		fmt.Printf("Gráfico %d: Tipo = %s\n", i+1, chart.ChartType)
		for _, data := range chart.FakeData {
			fmt.Printf("  - Nome: %s, Valor: %d\n", data.Name, data.Value)
		}
	}

}

func main() {
	extractGrapich()
}
