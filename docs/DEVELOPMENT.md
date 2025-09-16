# 🛠️ Guia de Desenvolvimento

Este documento explica como executar e desenvolver o sistema CEP Temperatura com OpenTelemetry + Zipkin.

## 🏗️ Arquitetura do Sistema

```
┌─────────────┐    HTTP     ┌─────────────┐    HTTP     ┌─────────────┐
│   Cliente   │ ──────────► │  Serviço A  │ ──────────► │  Serviço B  │
│             │             │ (Validação) │             │(Orquestração)│
└─────────────┘             └─────────────┘             └─────────────┘
                                    │                           │
                                    ▼                           ▼
                            ┌─────────────┐             ┌─────────────┐
                            │   Zipkin    │             │   ViaCEP    │
                            │ (Tracing)   │             │   WeatherAPI│
                            └─────────────┘             └─────────────┘
```

## 🚀 Execução Rápida

### Opção 1: Script Automatizado
```bash
./scripts/run-dev.sh
```

### Opção 2: Docker Compose
```bash
docker-compose up --build
```

### Opção 3: Manual
```bash
# Terminal 1 - Zipkin
docker run -d -p 9411:9411 openzipkin/zipkin

# Terminal 2 - Serviço B
go run cmd/service-b/main.go

# Terminal 3 - Serviço A
go run cmd/service-a/main.go
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `WEATHER_API_KEY` | Chave da WeatherAPI | Obrigatória |
| `SERVICE_B_URL` | URL do Serviço B | `http://localhost:8081` |
| `ZIPKIN_ENDPOINT` | Endpoint do Zipkin | `http://localhost:9411/api/v2/spans` |
| `PORT` | Porta do serviço | `8080` (A) / `8081` (B) |

### Arquivo .env
```bash
cp env.example .env
# Edite com suas configurações
```

## 🧪 Testes

### Teste Manual
```bash
# CEP válido
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01310100"}'

# CEP inválido
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"1234567"}'

# Health checks
curl http://localhost:8080/health
curl http://localhost:8081/health
```

### Teste Automatizado
```bash
./scripts/test-services.sh
```

### Health Checks

#### Health Básico
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
```

#### Health Detalhado (com dependências)
```bash
curl http://localhost:8080/health/detailed
curl http://localhost:8081/health/detailed
```

#### Readiness Check (Kubernetes)
```bash
curl http://localhost:8080/ready
curl http://localhost:8081/ready
```

#### Liveness Check (Kubernetes)
```bash
curl http://localhost:8080/live
curl http://localhost:8081/live
```

#### Teste Automatizado de Health
```bash
./scripts/test-health.sh
```

## 🔍 Observabilidade

### Zipkin UI
- **URL**: http://localhost:9411
- **Funcionalidades**:
  - Visualização de traces distribuídos
  - Análise de latência
  - Mapa de dependências
  - Detalhes de spans

### Spans Implementados

#### Serviço A
- `validate-cep`: Validação de formato do CEP
- `call-service-b`: Chamada HTTP para Serviço B

#### Serviço B
- `validate-cep`: Validação de formato do CEP
- `fetch-location`: Busca de localização via ViaCEP
- `fetch-temperature`: Busca de temperatura via WeatherAPI
- `convert-temperatures`: Conversão entre escalas

### Atributos dos Spans
- `cep`: CEP sendo processado
- `city`: Cidade encontrada
- `state`: Estado encontrado
- `temperature.celsius`: Temperatura em Celsius
- `temperature.fahrenheit`: Temperatura em Fahrenheit
- `temperature.kelvin`: Temperatura em Kelvin
- `weather.duration_ms`: Tempo de resposta da WeatherAPI
- `http.duration_ms`: Tempo de resposta HTTP
- `http.status_code`: Código de status HTTP

## 🐛 Debugging

### Logs
```bash
# Ver logs do Docker Compose
docker-compose logs -f

# Ver logs de um serviço específico
docker-compose logs -f service-a
docker-compose logs -f service-b
```

### Verificar Conectividade
```bash
# Verificar se Zipkin está rodando
curl http://localhost:9411/api/v2/services

# Verificar se Serviço B está acessível
curl http://localhost:8081/health

# Verificar se Serviço A está acessível
curl http://localhost:8080/health
```

### Problemas Comuns

1. **Erro 500 - WeatherAPI**
   - Verificar se `WEATHER_API_KEY` está configurada
   - Verificar conectividade com api.weatherapi.com

2. **Erro de conexão entre serviços**
   - Verificar se `SERVICE_B_URL` está correto
   - Verificar se Serviço B está rodando

3. **Traces não aparecem no Zipkin**
   - Verificar se `ZIPKIN_ENDPOINT` está correto
   - Verificar se Zipkin está rodando
   - Aguardar alguns segundos para propagação

## 📊 Monitoramento

### Métricas Importantes
- **Latência**: Tempo total de processamento
- **Throughput**: Requisições por segundo
- **Error Rate**: Taxa de erro por serviço
- **Dependencies**: Tempo de resposta de APIs externas

### Alertas Sugeridos
- Latência > 5 segundos
- Error rate > 5%
- Falha na conectividade com APIs externas

## 🔄 Desenvolvimento

### Estrutura do Projeto
```
cmd/
├── service-a/main.go    # Serviço A
└── service-b/main.go    # Serviço B

internal/
├── handlers/            # HTTP handlers
├── services/            # Lógica de negócio
├── models/              # Estruturas de dados
├── config/              # Configuração
└── telemetry/           # OpenTelemetry

scripts/
├── dev-setup.sh         # Setup inicial
├── run-dev.sh           # Execução local
└── test-services.sh     # Testes automatizados
```

### Adicionando Novos Spans
```go
// Criar span
ctx, span := tracer.Start(ctx, "nome-do-span")
defer span.End()

// Adicionar atributos
span.SetAttributes(
    attribute.String("key", "value"),
    attribute.Int64("number", 123),
)

// Registrar erro
span.RecordError(err)
```

### Adicionando Novos Serviços
1. Criar novo `cmd/service-x/main.go`
2. Adicionar ao `docker-compose.yml`
3. Criar `Dockerfile.service-x`
4. Atualizar configurações
5. Adicionar testes

## 📚 Referências

- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/getting-started/)
- [Zipkin](https://zipkin.io/)
- [Gin Framework](https://gin-gonic.com/)
- [ViaCEP API](https://viacep.com.br/)
- [WeatherAPI](https://www.weatherapi.com/)
