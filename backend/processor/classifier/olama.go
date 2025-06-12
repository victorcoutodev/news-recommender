package classifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type OllamaClassifier struct {
	Model   string
	BaseURL string
}

func NewOllamaClassifier(model string) *OllamaClassifier {
	if model == "" {
		model = "mistral"
	}

	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://host.docker.internal:11434" // padrão local
	}

	return &OllamaClassifier{
		Model:   model,
		BaseURL: baseURL,
	}
}

func (o *OllamaClassifier) Classify(text string) (string, error) {
	prompt := fmt.Sprintf(`Classifique o texto em uma das categorias a seguir:
tecnologia, política, esportes, economia, saúde, entretenimento, educação, segurança.
Responda APENAS com a categoria, sem explicação. Texto: "%s"`, text)

	url := fmt.Sprintf("%s/api/generate", o.BaseURL)

	payload := map[string]interface{}{
		"model":  o.Model,
		"prompt": prompt,
		"stream": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("erro ao criar payload JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("erro na requisição ao Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("erro no servidor Ollama: status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta JSON: %v", err)
	}

	cat := normalizeCategory(result.Response)
	return cat, nil
}

func normalizeCategory(raw string) string {
	cat := strings.ToLower(strings.TrimSpace(raw))
	switch cat {
	case "tecnologia", "tech":
		return "tecnologia"
	case "política", "politica":
		return "política"
	case "esportes", "esporte":
		return "esportes"
	case "economia":
		return "economia"
	case "saúde", "saude":
		return "saúde"
	case "entretenimento":
		return "entretenimento"
	case "educação", "educacao":
		return "educação"
	case "segurança", "seguranca":
		return "segurança"
	default:
		return "outra"
	}
}
