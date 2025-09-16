#!/bin/bash

set -e

echo "🎬 Demo da API CEP Temperatura"
echo "================================"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para fazer requisição e mostrar resultado
make_request() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    
    echo -e "\n${BLUE}📡 $description${NC}"
    echo "URL: $method $url"
    if [ -n "$data" ]; then
        echo "Data: $data"
    fi
    
    if [ "$method" = "POST" ]; then
        response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" "$url")
    fi
    
    # Separar body e status code
    body=$(echo "$response" | head -n -1)
    status=$(echo "$response" | tail -n 1)
    
    if [ "$status" -ge 200 ] && [ "$status" -lt 300 ]; then
        echo -e "${GREEN}✅ Status: $status${NC}"
    else
        echo -e "${RED}❌ Status: $status${NC}"
    fi
    
    echo "Response:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
}

# Verificar se os serviços estão rodando
echo -e "\n${YELLOW}🔍 Verificando se os serviços estão rodando...${NC}"

if ! curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${RED}❌ Serviço A não está rodando em http://localhost:8080${NC}"
    echo "Execute: sudo docker compose up --build -d"
    exit 1
fi

if ! curl -s http://localhost:8081/health > /dev/null; then
    echo -e "${RED}❌ Serviço B não está rodando em http://localhost:8081${NC}"
    echo "Execute: sudo docker compose up --build -d"
    exit 1
fi

echo -e "${GREEN}✅ Serviços estão rodando!${NC}"

# 1. Health Checks
echo -e "\n${YELLOW}🏥 1. Testando Health Checks${NC}"
make_request "GET" "http://localhost:8080/health" "" "Health Check - Serviço A"
make_request "GET" "http://localhost:8081/health" "" "Health Check - Serviço B"
make_request "GET" "http://localhost:8080/health/detailed" "" "Health Check Detalhado - Serviço A"
make_request "GET" "http://localhost:8081/health/detailed" "" "Health Check Detalhado - Serviço B"

# 2. CEP Válido
echo -e "\n${YELLOW}🌡️ 2. Testando CEP Válido (São Paulo)${NC}"
make_request "POST" "http://localhost:8080/cep" '{"cep":"01310100"}' "CEP: 01310-100 (São Paulo)"

# 3. CEP Válido - Rio de Janeiro
echo -e "\n${YELLOW}🌡️ 3. Testando CEP Válido (Rio de Janeiro)${NC}"
make_request "POST" "http://localhost:8080/cep" '{"cep":"20040020"}' "CEP: 20040-020 (Rio de Janeiro)"

# 4. CEP Inválido - Formato
echo -e "\n${YELLOW}❌ 4. Testando CEP Inválido (Formato)${NC}"
make_request "POST" "http://localhost:8080/cep" '{"cep":"1234567"}' "CEP: 1234567 (7 dígitos - inválido)"

# 5. CEP Inválido - Caracteres
echo -e "\n${YELLOW}❌ 5. Testando CEP Inválido (Caracteres)${NC}"
make_request "POST" "http://localhost:8080/cep" '{"cep":"abc12345"}' "CEP: abc12345 (com letras - inválido)"

# 6. CEP Não Encontrado
echo -e "\n${YELLOW}🔍 6. Testando CEP Não Encontrado${NC}"
make_request "POST" "http://localhost:8080/cep" '{"cep":"99999999"}' "CEP: 99999999 (não existe)"

# 7. Teste de HEAD requests
echo -e "\n${YELLOW}🔍 7. Testando HEAD Requests${NC}"
echo "HEAD /ready - Serviço A:"
curl -s -I http://localhost:8080/ready | head -n 1
echo "HEAD /ready - Serviço B:"
curl -s -I http://localhost:8081/ready | head -n 1

# 8. Informações sobre Tracing
echo -e "\n${YELLOW}📊 8. Informações sobre Tracing${NC}"
echo "Para visualizar os traces:"
echo "1. Acesse: http://localhost:9411"
echo "2. Clique em 'Run Query'"
echo "3. Clique em um trace para ver detalhes"
echo "4. Analise a timeline e spans individuais"

echo -e "\n${GREEN}🎉 Demo concluída!${NC}"
echo "Verifique os traces em: http://localhost:9411"
