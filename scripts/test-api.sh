#!/bin/bash

# Script para testar a API de temperatura por CEP
# Uso: ./test-api.sh [URL_BASE]
# Exemplo: ./test-api.sh http://localhost:8080

BASE_URL=${1:-"http://localhost:8080"}

echo "🧪 Testando API de Temperatura por CEP"
echo "📍 URL Base: $BASE_URL"
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_status=$3
    local description=$4
    
    echo -n "Testando $description... "
    
    response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" --connect-timeout 10)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ PASSOU${NC} (HTTP $http_code)"
        if [ "$expected_status" = "200" ]; then
            echo "   Resposta: $body"
        fi
    else
        echo -e "${RED}❌ FALHOU${NC} (Esperado: HTTP $expected_status, Recebido: HTTP $http_code)"
        echo "   Resposta: $body"
    fi
    echo ""
}

# Teste 1: Health check
test_endpoint "GET" "/health" "200" "Health Check"

# Teste 2: CEP válido
test_endpoint "GET" "/temperature/01310100" "200" "CEP Válido (01310-100)"

# Teste 3: CEP válido sem hífen
test_endpoint "GET" "/temperature/20040020" "200" "CEP Válido sem hífen (20040-020)"

# Teste 4: CEP inválido - formato incorreto
test_endpoint "GET" "/temperature/123" "422" "CEP Inválido - formato incorreto"

# Teste 5: CEP inválido - com letras
test_endpoint "GET" "/temperature/1234567a" "422" "CEP Inválido - com letras"

# Teste 6: CEP não encontrado
test_endpoint "GET" "/temperature/99999999" "404" "CEP Não Encontrado"

echo "🏁 Testes concluídos!"
