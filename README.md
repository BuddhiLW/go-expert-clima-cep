# 🌡️ CEP Temperatura API

Sistema distribuído em Go com OpenTelemetry + Zipkin que recebe um CEP brasileiro e retorna a temperatura atual da cidade em Celsius, Fahrenheit e Kelvin.

## 🏗️ Arquitetura

- **Serviço A**: Validação de entrada e proxy para Serviço B
- **Serviço B**: Orquestração (busca CEP → busca temperatura → conversão)
- **Zipkin**: Tracing distribuído e observabilidade
- **OpenTelemetry**: Instrumentação e coleta de métricas

## 🚀 API Live

**URL**: https://cep-temperatura-667491814881.southamerica-east1.run.app

### Exemplos de Uso

```bash
# Temperatura de São Paulo (via Serviço A)
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01310100"}'

# Health checks
curl http://localhost:8080/health  # Serviço A
curl http://localhost:8081/health  # Serviço B
```

## 📋 Requisitos

- Go 1.25+
- Chave da WeatherAPI (obtenha em https://www.weatherapi.com/)

## 🛠️ Instalação Local

### Opção 1: Docker Compose (Recomendado)

1. **Clone e configure**
```bash
git clone <repository-url>
cd cep-temperatura
./scripts/dev-setup.sh
```

2. **Inicie os serviços**
```bash
docker-compose up --build
```

3. **Teste a aplicação**
```bash
./scripts/test-services.sh
```

### Opção 2: Desenvolvimento Local

1. **Configure o ambiente**
```bash
cp env.example .env
# Edite .env com sua WEATHER_API_KEY
```

2. **Execute os serviços**
```bash
# Terminal 1 - Serviço B
go run cmd/service-b/main.go

# Terminal 2 - Serviço A  
go run cmd/service-a/main.go

# Terminal 3 - Zipkin
docker run -d -p 9411:9411 openzipkin/zipkin
```

## 🧪 Testes

```bash
# Todos os testes
go test ./...

# Teste da API
./test-api.sh
```

## 📡 Endpoints

### Serviço A (Porta 8080)

#### POST /cep

Valida CEP e encaminha para Serviço B.

**Request:**
```json
{
  "cep": "01310100"
}
```

**Resposta de sucesso (200):**
```json
{
  "city": "São Paulo",
  "temp_C": 21.4,
  "temp_F": 70.52,
  "temp_K": 294.4
}
```

**Códigos de erro:**
- `422` - CEP inválido (não tem 8 dígitos)
- `404` - CEP não encontrado
- `500` - Erro interno

#### GET /health

Health check do Serviço A.

### Serviço B (Porta 8081)

#### GET /temperature/:cep

Busca temperatura para CEP (usado internamente pelo Serviço A).

#### GET /health

Health check do Serviço B.

## 🔍 Observabilidade

### Zipkin UI
- **URL**: http://localhost:9411
- **Funcionalidades**: Traces distribuídos, latência, dependências

### Spans Implementados
- `validate-cep`: Validação de formato do CEP
- `fetch-location`: Busca de localização via ViaCEP
- `fetch-temperature`: Busca de temperatura via WeatherAPI
- `convert-temperatures`: Conversão entre escalas
- `call-service-b`: Chamada entre serviços

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