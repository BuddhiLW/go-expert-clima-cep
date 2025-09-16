#!/bin/bash

# Script para testar HEAD requests nos health endpoints

set -e

echo "🔍 Testando HEAD requests nos health endpoints..."

# Aguardar serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 5

# Função para testar HEAD request
test_head_request() {
    local url=$1
    local name=$2
    
    echo "🔍 Testando HEAD $name..."
    echo "   URL: $url"
    
    response=$(curl -s -I "$url")
    http_code=$(echo "$response" | head -n1 | cut -d' ' -f2)
    
    if [ "$http_code" -eq "200" ]; then
        echo "   ✅ Status: $http_code"
        echo "   📄 Headers:"
        echo "$response" | head -n5
    else
        echo "   ❌ Status: $http_code"
        echo "   📄 Response:"
        echo "$response"
    fi
    echo ""
}

# Testar Serviço A
echo "=== SERVIÇO A (Porta 8080) ==="
test_head_request "http://localhost:8080/health" "Health Básico"
test_head_request "http://localhost:8080/ready" "Readiness Check"
test_head_request "http://localhost:8080/live" "Liveness Check"

echo "=== SERVIÇO B (Porta 8081) ==="
test_head_request "http://localhost:8081/health" "Health Básico"
test_head_request "http://localhost:8081/ready" "Readiness Check"
test_head_request "http://localhost:8081/live" "Liveness Check"

echo "✅ Testes de HEAD requests concluídos!"
echo ""
echo "💡 Agora o Docker Compose deve conseguir fazer health checks corretamente!"
