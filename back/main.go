package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// Importe "google.golang.org/api/option" foi REMOVIDO
	"github.com/joho/godotenv"
	"google.golang.org/genai" // Apenas este import do genai é necessário
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

type ImageRequest struct {
	ImageData string `json:"image_data" binding:"required"`
}

func extractGrapich(c *gin.Context) {
    var req ImageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
        return
    }

    decodedBytes, err := base64.StdEncoding.DecodeString(req.ImageData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode base64 image"})
		fmt.Println(err)
        return
    }

	fileName := fmt.Sprintf("uploads/img_%d.png", time.Now().Unix())

// Garante que a pasta existe
if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar diretório"})
    return
}

// Escreve o arquivo
if err := os.WriteFile(fileName, decodedBytes, 0644); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar a imagem"})
    return
}

fmt.Println("Imagem salva em:", fileName)

    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "API key não definida"})
		fmt.Println(err)
        return
    }


    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao inicializar cliente"})
		fmt.Println(err)
        return
    }

    parts := []*genai.Part{
        genai.NewPartFromBytes(decodedBytes, "image/png"),
        genai.NewPartFromText(prompt_text),
    }

    contents := []*genai.Content{
        genai.NewContentFromParts(parts, genai.RoleUser),
    }

    result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha na requisição ao Gemini"})
		fmt.Println(err)
        return
    }

    // ✅ Proteção contra resposta vazia
    if len(result.Candidates) == 0 || 
       result.Candidates[0].Content == nil || 
       len(result.Candidates[0].Content.Parts) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Resposta vazia do Gemini"})
		fmt.Println(err)
        return
    }

    rawResponse := result.Text()
    cleanJSON := strings.TrimSpace(rawResponse)
    cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
    cleanJSON = strings.TrimSuffix(cleanJSON, "```")

    var dashboardData GeminiDashboardResponse
    if err := json.Unmarshal([]byte(cleanJSON), &dashboardData); err != nil {
        // Em vez de panic, devolve erro pro cliente
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao converter JSON", "raw": cleanJSON})
		fmt.Println((err))
        return
    }

    c.JSON(http.StatusOK, result)
}

func main() {
	// Carrega as variáveis do arquivo .env no início da execução
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	router := gin.Default();
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	fmt.Println("Backend Started on PORT 8080")

	router.POST("/generateImage", extractGrapich)


	router.Run()

	// extractGrapich()
}