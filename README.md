# 🌡️ CEP Temperatura API

Sistema em Go que recebe um CEP brasileiro e retorna a temperatura atual da cidade em Celsius, Fahrenheit e Kelvin.

## 🚀 API Live

**URL**: https://cep-temperatura-667491814881.southamerica-east1.run.app

### Exemplos de Uso

```bash
# Temperatura de São Paulo
curl "https://cep-temperatura-667491814881.southamerica-east1.run.app/temperature/01310100"

# Health check
curl "https://cep-temperatura-667491814881.southamerica-east1.run.app/health"
```

## 📋 Requisitos

- Go 1.25+
- Chave da WeatherAPI (obtenha em https://www.weatherapi.com/)

## 🛠️ Instalação Local

1. **Clone e configure**
```bash
# Clone 
git clone <repository-url>
cd cep-temperatura

# Setup env
cp env.example .env
# preencha com sua API

```

2. **Execute**
```bash
# rodar docker
./scripts/test-docker.sh

# ou,
go run cmd/main.go
```

## 🧪 Testes

```bash
# Todos os testes
go test ./...

# Teste da API
./test-api.sh
```

## 📡 Endpoints

### GET /temperature/:cep

Retorna temperatura para o CEP informado.

**Exemplo de resposta:**
```json
{
  "temp_C": 21.4,
  "temp_F": 70.52,
  "temp_K": 294.4
}
```

**Códigos de erro:**
- `422` - CEP inválido (não tem 8 dígitos)
- `404` - CEP não encontrado

### GET /health

Verificação de saúde da API.

## 🏗️ Arquitetura

```
internal/
├── handlers/     # HTTP handlers
├── services/     # Lógica de negócio
└── models/       # Estruturas de dados
```

## 🚀 Deploy

### Deploy Automático (Cloud Build)

```bash
./deploy-cloudbuild.sh SEU_PROJECT_ID
```

### Deploy Manual

```bash
# Build e push
sudo docker build -t gcr.io/SEU_PROJECT_ID/cep-temperatura .
sudo docker push gcr.io/SEU_PROJECT_ID/cep-temperatura

# Deploy
gcloud run deploy cep-temperatura \
  --image gcr.io/SEU_PROJECT_ID/cep-temperatura \
  --platform managed \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=SUA_CHAVE
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `PORT` | Porta do servidor | `8080` |
| `HOST` | Host do servidor | `0.0.0.0` |
| `WEATHER_API_KEY` | Chave da WeatherAPI | Obrigatória (obtenha em weatherapi.com) |

### APIs Externas

- **ViaCEP**: https://viacep.com.br/ (gratuita)
- **WeatherAPI**: https://www.weatherapi.com/ (requer chave)

## 🐛 Troubleshooting

1. **Erro 500** - Verifique a chave da WeatherAPI
2. **Erro 422** - CEP deve ter exatamente 8 dígitos
3. **Erro 404** - CEP não existe no ViaCEP

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.