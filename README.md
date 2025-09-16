# 🌡️ Sistema de Temperatura por CEP

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## 🚀 Funcionalidades

- ✅ Validação de CEP (8 dígitos)
- ✅ Consulta de localização via ViaCEP API
- ✅ Consulta de temperatura via WeatherAPI
- ✅ Conversão de temperaturas (Celsius, Fahrenheit, Kelvin)
- ✅ Tratamento de erros adequado
- ✅ Testes automatizados
- ✅ Containerização com Docker
- ✅ Deploy no Google Cloud Run

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture** e **SOLID**:

```
internal/
├── handlers/     # Camada de apresentação (HTTP handlers)
├── services/     # Camada de domínio (regras de negócio)
└── models/       # Entidades e DTOs
```

### Princípios Aplicados

- **SRP**: Cada serviço tem uma responsabilidade específica
- **OCP**: Fácil extensão sem modificação
- **LSP**: Interfaces bem definidas
- **ISP**: Interfaces específicas para cada cliente
- **DIP**: Dependência de abstrações, não implementações

## 📋 Requisitos

- Go 1.25+
- Docker (opcional)
- Chave da WeatherAPI

## 🛠️ Instalação e Execução

### Execução Local

1. **Clone o repositório**
```bash
git clone <repository-url>
cd cep-temperatura
```

2. **Configure o ambiente**
```bash
./setup-env.sh
```

3. **Execute a aplicação**
```bash
./main
```

**Ou manualmente:**

1. **Instale as dependências**
```bash
go mod tidy
```

2. **Crie o arquivo .env**
```bash
cp env.example .env
# Edite o arquivo .env com suas configurações
```

3. **Execute a aplicação**
```bash
go run cmd/main.go
```

### Execução com Docker

1. **Build e execução**
```bash
docker-compose up --build
```

2. **Apenas build**
```bash
docker build -t cep-temperatura .
```

3. **Execução do container**
```bash
docker run -p 8080:8080 cep-temperatura
```

## 🧪 Testes

### Executar todos os testes
```bash
go test ./...
```

### Executar testes com verbose
```bash
go test -v ./...
```

### Executar testes de um pacote específico
```bash
go test ./internal/services
```

### Testar API manualmente
```bash
./test-api.sh
```

## 📡 Endpoints da API

### GET /temperature/:cep

Retorna a temperatura atual para a localidade do CEP informado.

**Parâmetros:**
- `cep`: CEP válido com 8 dígitos numéricos

**Respostas:**

#### ✅ Sucesso (200 OK)
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### ❌ CEP Inválido (422 Unprocessable Entity)
```json
{
  "message": "invalid zipcode"
}
```

#### ❌ CEP Não Encontrado (404 Not Found)
```json
{
  "message": "can not find zipcode"
}
```

### GET /health

Endpoint de verificação de saúde da API.

**Resposta:**
```json
{
  "status": "ok"
}
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão | Arquivo |
|----------|-----------|--------|---------|
| `PORT` | Porta do servidor | `8080` | `.env` |
| `HOST` | Host do servidor | `0.0.0.0` | `.env` |
| `WEATHER_API_KEY` | Chave da WeatherAPI | `b5d4215a52bf4e2da2f144209251609` | `.env` |
| `WEATHER_BASE_URL` | URL base da WeatherAPI | `http://api.weatherapi.com/v1` | `.env` |

### Arquivos de Configuração

- **`.env`** - Variáveis de ambiente (criado automaticamente pelo `setup-env.sh`)
- **`configs/config.yaml`** - Configuração padrão
- **`configs/config.dev.yaml`** - Configuração de desenvolvimento
- **`configs/config.prod.yaml`** - Configuração de produção

### APIs Externas

- **ViaCEP**: https://viacep.com.br/ (gratuita)
- **WeatherAPI**: https://www.weatherapi.com/ (requer chave)

## 🚀 Deploy no Google Cloud Run

### Pré-requisitos

1. Google Cloud SDK instalado
2. Projeto no Google Cloud configurado
3. Docker instalado

### Passos para Deploy

1. **Configure o projeto**
```bash
gcloud config set project SEU_PROJECT_ID
```

2. **Build e push da imagem**
```bash
# Build da imagem
docker build -t gcr.io/SEU_PROJECT_ID/cep-temperatura .

# Push para Google Container Registry
docker push gcr.io/SEU_PROJECT_ID/cep-temperatura
```

3. **Deploy no Cloud Run**
```bash
gcloud run deploy cep-temperatura \
  --image gcr.io/SEU_PROJECT_ID/cep-temperatura \
  --platform managed \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=b5d4215a52bf4e2da2f144209251609
```

### Deploy Automatizado

O projeto inclui um workflow GitHub Actions para deploy automático (ver `.github/workflows/deploy.yml`).

## 📊 Exemplos de Uso

### cURL

```bash
# CEP válido
curl "https://sua-api.run.app/temperature/01310100"

# Health check
curl "https://sua-api.run.app/health"
```

### JavaScript

```javascript
const response = await fetch('https://sua-api.run.app/temperature/01310100');
const data = await response.json();
console.log(data);
// { temp_C: 28.5, temp_F: 83.3, temp_K: 301.5 }
```

### Python

```python
import requests

response = requests.get('https://sua-api.run.app/temperature/01310100')
data = response.json()
print(data)
# {'temp_C': 28.5, 'temp_F': 83.3, 'temp_K': 301.5}
```

## 🧮 Fórmulas de Conversão

- **Celsius para Fahrenheit**: `F = C × 1.8 + 32`
- **Celsius para Kelvin**: `K = C + 273`

## 🐛 Troubleshooting

### Problemas Comuns

1. **Erro 500 - Erro interno do servidor**
   - Verifique se a chave da WeatherAPI está correta
   - Verifique a conectividade com as APIs externas

2. **Erro 422 - CEP inválido**
   - Certifique-se de que o CEP tem exatamente 8 dígitos
   - Remova hífens e espaços

3. **Erro 404 - CEP não encontrado**
   - Verifique se o CEP existe no ViaCEP
   - Teste com CEPs conhecidos (ex: 01310-100)

### Logs

Para ver logs detalhados:
```bash
# Local
go run main.go

# Docker
docker-compose logs -f

# Cloud Run
gcloud logs read --service=cep-temperatura --limit=50
```

## 📈 Monitoramento

### Health Check

```bash
curl https://sua-api.run.app/health
```

### Métricas

- Tempo de resposta
- Taxa de sucesso/erro
- Uso de CPU e memória (Cloud Run)

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👥 Autores

- **Ramanujan** - *Desenvolvimento inicial* - [GitHub](https://github.com/ramanujan)

## 🙏 Agradecimentos

- [ViaCEP](https://viacep.com.br/) - API gratuita de CEPs
- [WeatherAPI](https://www.weatherapi.com/) - API de clima
- [Gin](https://gin-gonic.com/) - Framework web em Go
- [Google Cloud Run](https://cloud.google.com/run) - Plataforma de deploy

---

⭐ **Se este projeto foi útil para você, considere dar uma estrela!** ⭐