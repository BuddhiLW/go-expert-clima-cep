#!/bin/bash

# Script para testar endpoints de health

set -e

echo "🏥 Testando Health Endpoints..."

# Aguardar serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 5

# Função para testar endpoint
test_endpoint() {
    local url=$1
    local name=$2
    local expected_status=${3:-200}
    
    echo "🔍 Testando $name..."
    echo "   URL: $url"
    
    response=$(curl -s -w "\n%{http_code}" "$url")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo "   ✅ Status: $http_code (esperado: $expected_status)"
        echo "   📄 Response:"
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo "   ❌ Status: $http_code (esperado: $expected_status)"
        echo "   📄 Response:"
        echo "$body"
    fi
    echo ""
}

# Testar Serviço A
echo "=== SERVIÇO A (Porta 8080) ==="
test_endpoint "http://localhost:8080/health" "Health Básico"
test_endpoint "http://localhost:8080/health/detailed" "Health Detalhado"
test_endpoint "http://localhost:8080/ready" "Readiness Check"
test_endpoint "http://localhost:8080/live" "Liveness Check"

echo "=== SERVIÇO B (Porta 8081) ==="
test_endpoint "http://localhost:8081/health" "Health Básico"
test_endpoint "http://localhost:8081/health/detailed" "Health Detalhado"
test_endpoint "http://localhost:8081/ready" "Readiness Check"
test_endpoint "http://localhost:8081/live" "Liveness Check"

echo "✅ Testes de health concluídos!"
echo ""
echo "💡 Dicas:"
echo "   - /health: Status básico do serviço"
echo "   - /health/detailed: Verifica dependências (Serviço A verifica Serviço B)"
echo "   - /ready: Usado pelo Kubernetes para readiness probe"
echo "   - /live: Usado pelo Kubernetes para liveness probe"
