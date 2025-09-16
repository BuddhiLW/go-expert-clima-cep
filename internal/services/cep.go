package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"cep-temperatura/internal/models"
)

// CEPService interface para operações de CEP
type CEPService interface {
	ValidateCEP(cep string) bool
	GetLocation(cep string) (*models.CEPResponse, error)
}

type cepService struct {
	baseURL string
	client  *http.Client
}

// NewCEPService cria uma nova instância do serviço de CEP
func NewCEPService() CEPService {
	return &cepService{
		baseURL: "https://viacep.com.br/ws",
		client:  &http.Client{},
	}
}

// ValidateCEP valida se o CEP está no formato correto
func (s *cepService) ValidateCEP(cep string) bool {
	return validateCEP(cep)
}

// GetLocation busca a localização pelo CEP
func (s *cepService) GetLocation(cep string) (*models.CEPResponse, error) {
	if !s.ValidateCEP(cep) {
		return nil, fmt.Errorf("invalid zipcode")
	}

	formattedCEP := formatCEP(cep)
	url := fmt.Sprintf("%s/%s/json/", s.baseURL, formattedCEP)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar CEP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var cepResponse models.CEPResponse
	if err := json.Unmarshal(body, &cepResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	if cepResponse.Erro {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &cepResponse, nil
}

// validateCEP valida se o CEP está no formato correto (8 dígitos)
func validateCEP(cep string) bool {
	// Remove hífens e espaços
	cleanCEP := strings.ReplaceAll(cep, "-", "")
	cleanCEP = strings.ReplaceAll(cleanCEP, " ", "")

	// Verifica se tem exatamente 8 dígitos
	matched, _ := regexp.MatchString(`^\d{8}$`, cleanCEP)
	return matched
}

// formatCEP formata o CEP removendo hífens e espaços
func formatCEP(cep string) string {
	cleanCEP := strings.ReplaceAll(cep, "-", "")
	cleanCEP = strings.ReplaceAll(cleanCEP, " ", "")
	return cleanCEP
}
