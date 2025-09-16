# 🐳 Deploy com Docker - Guia Completo

Este guia explica como construir e fazer deploy da aplicação usando Docker com variáveis de ambiente.

## 📋 Pré-requisitos

- Docker instalado
- Arquivo `.env` configurado
- Google Cloud SDK (para deploy)

## 🔧 Configuração Inicial

### 1. Configurar Ambiente

```bash
# Configurar automaticamente
./setup-env.sh

# Ou manualmente
cp env.example .env
# Edite o arquivo .env com suas configurações
```

### 2. Verificar Configuração

```bash
# Verificar se o arquivo .env está correto
cat .env
```

## 🏗️ Construção da Imagem

### Opção 1: Script Automatizado (Recomendado)

```bash
# Build com script
./build.sh

# Build com tag específica
./build.sh cep-temperatura:v1.0.0

# Build com arquivo de ambiente específico
./build.sh cep-temperatura:latest .env.prod
```

### Opção 2: Build Manual

```bash
# Carregar variáveis de ambiente
export $(grep -v '^#' .env | xargs)

# Build da imagem
docker build \
  --build-arg WEATHER_API_KEY="$WEATHER_API_KEY" \
  --build-arg PORT="${PORT:-8080}" \
  --build-arg HOST="${HOST:-0.0.0.0}" \
  -t cep-temperatura:latest .
```

## 🧪 Teste Local

### Teste Automatizado

```bash
# Teste completo da imagem
./test-docker.sh

# Teste com tag específica
./test-docker.sh cep-temperatura:v1.0.0
```

### Teste Manual

```bash
# Executar container
docker run -d --name cep-temperatura -p 8080:8080 cep-temperatura:latest

# Testar endpoints
curl http://localhost:8080/health
curl http://localhost:8080/temperature/01310100

# Parar container
docker stop cep-temperatura
docker rm cep-temperatura
```

## 🚀 Deploy no Google Cloud Run

### Opção 1: Deploy com Docker Local

```bash
# Deploy completo
./deploy.sh SEU_PROJECT_ID

# Deploy com arquivo de ambiente específico
ENV_FILE=.env.prod ./deploy.sh SEU_PROJECT_ID
```

### Opção 2: Deploy com Cloud Build (Recomendado)

```bash
# Deploy com Cloud Build
./deploy-cloudbuild.sh SEU_PROJECT_ID

# Deploy com arquivo de ambiente específico
./deploy-cloudbuild.sh SEU_PROJECT_ID .env.prod
```

### Opção 3: Deploy Manual

```bash
# 1. Configurar projeto
gcloud config set project SEU_PROJECT_ID

# 2. Carregar variáveis
export $(grep -v '^#' .env | xargs)

# 3. Build e push
docker build \
  --build-arg WEATHER_API_KEY="$WEATHER_API_KEY" \
  --build-arg PORT="${PORT:-8080}" \
  --build-arg HOST="${HOST:-0.0.0.0}" \
  -t gcr.io/SEU_PROJECT_ID/cep-temperatura:latest .

docker push gcr.io/SEU_PROJECT_ID/cep-temperatura:latest

# 4. Deploy
gcloud run deploy cep-temperatura \
  --image gcr.io/SEU_PROJECT_ID/cep-temperatura:latest \
  --platform managed \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY="$WEATHER_API_KEY",PORT="${PORT:-8080}",HOST="${HOST:-0.0.0.0}"
```

## 🔍 Verificação do Deploy

### 1. Verificar Status

```bash
# Status do serviço
gcloud run services describe cep-temperatura \
  --platform managed \
  --region southamerica-east1

# Logs do serviço
gcloud logs tail --service=cep-temperatura
```

### 2. Testar Endpoints

```bash
# Obter URL do serviço
SERVICE_URL=$(gcloud run services describe cep-temperatura \
  --platform managed \
  --region southamerica-east1 \
  --format 'value(status.url)')

# Testar health check
curl "$SERVICE_URL/health"

# Testar CEP
curl "$SERVICE_URL/temperature/01310100"
```

## 🐛 Troubleshooting

### Problemas Comuns

1. **Erro de build: "WEATHER_API_KEY not found"**
   ```bash
   # Verificar se o arquivo .env existe e tem a variável
   cat .env | grep WEATHER_API_KEY
   ```

2. **Erro de push: "permission denied"**
   ```bash
   # Fazer login no Google Container Registry
   gcloud auth configure-docker
   ```

3. **Erro de deploy: "service not found"**
   ```bash
   # Verificar se o projeto está correto
   gcloud config get-value project
   ```

4. **Container não inicia**
   ```bash
   # Verificar logs do container
   docker logs cep-temperatura-test
   ```

### Debug da Imagem

```bash
# Executar container em modo interativo
docker run -it --rm cep-temperatura:latest sh

# Verificar variáveis de ambiente dentro do container
docker run --rm cep-temperatura:latest env
```

## 📊 Monitoramento

### Métricas do Container

```bash
# Estatísticas do container
docker stats cep-temperatura-test

# Logs em tempo real
docker logs -f cep-temperatura-test
```

### Métricas do Cloud Run

- Acesse o [Google Cloud Console](https://console.cloud.google.com/)
- Vá para Cloud Run > cep-temperatura
- Visualize métricas de CPU, memória, requisições

## 🔄 Atualizações

### Atualizar Imagem

```bash
# 1. Atualizar código
git pull

# 2. Rebuild da imagem
./build.sh cep-temperatura:v2.0.0

# 3. Deploy da nova versão
./deploy-cloudbuild.sh SEU_PROJECT_ID
```

### Rollback

```bash
# Listar revisões
gcloud run revisions list --service=cep-temperatura

# Fazer rollback
gcloud run services update-traffic cep-temperatura \
  --to-revisions=REVISION_NAME=100
```

## 💡 Dicas

1. **Use tags semânticas** para versionamento
2. **Teste localmente** antes do deploy
3. **Monitore logs** após o deploy
4. **Use Cloud Build** para deploys automáticos
5. **Configure CI/CD** com GitHub Actions

---

🎉 **Deploy concluído com sucesso!** Sua aplicação está rodando no Google Cloud Run.
