#!/bin/bash
set -e

echo "🧪 Running integration tests..."

# Start services
echo "Starting docker-compose services..."
docker-compose up -d

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
timeout=60
elapsed=0
while [ $elapsed -lt $timeout ]; do
    if docker-compose ps | grep -q "healthy"; then
        echo "Services are healthy!"
        break
    fi
    sleep 2
    elapsed=$((elapsed + 2))
done

if [ $elapsed -ge $timeout ]; then
    echo "❌ Services failed to become healthy"
    docker-compose logs
    docker-compose down
    exit 1
fi

# Test health endpoint
echo "Testing /health endpoint..."
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$response" != "200" ]; then
    echo "❌ Health check failed with status $response"
    docker-compose down
    exit 1
fi
echo "✅ Health check passed"

# Generate test token
echo "Generating test token..."
token_response=$(curl -s -X POST http://localhost:8080/auth/token \
    -H "Content-Type: application/json" \
    -d '{"user_id": "test-user", "role": "admin"}')
token=$(echo $token_response | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$token" ]; then
    echo "❌ Failed to generate token"
    docker-compose down
    exit 1
fi
echo "✅ Token generated"

# Test CRUD operations
echo "Testing service creation..."
create_response=$(curl -s -X POST http://localhost:8080/api/v1/services \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $token" \
    -d '{
        "name": "test-service",
        "region": "us-east-1",
        "image": "nginx",
        "version": "1.0.0"
    }')

service_id=$(echo $create_response | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
if [ -z "$service_id" ]; then
    echo "❌ Failed to create service"
    echo "Response: $create_response"
    docker-compose down
    exit 1
fi
echo "✅ Service created with ID: $service_id"

# Test list services
echo "Testing list services..."
list_response=$(curl -s -X GET "http://localhost:8080/api/v1/services?limit=10" \
    -H "Authorization: Bearer $token")
if ! echo "$list_response" | grep -q "services"; then
    echo "❌ Failed to list services"
    docker-compose down
    exit 1
fi
echo "✅ List services passed"

# Test get service
echo "Testing get service..."
get_response=$(curl -s -X GET "http://localhost:8080/api/v1/services/$service_id" \
    -H "Authorization: Bearer $token")
if ! echo "$get_response" | grep -q "$service_id"; then
    echo "❌ Failed to get service"
    docker-compose down
    exit 1
fi
echo "✅ Get service passed"

# Test update service
echo "Testing update service..."
update_response=$(curl -s -X PATCH "http://localhost:8080/api/v1/services/$service_id" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $token" \
    -d '{"status": "running"}')
if ! echo "$update_response" | grep -q "running"; then
    echo "❌ Failed to update service"
    docker-compose down
    exit 1
fi
echo "✅ Update service passed"

# Test delete service
echo "Testing delete service..."
delete_response=$(curl -s -o /dev/null -w "%{http_code}" \
    -X DELETE "http://localhost:8080/api/v1/services/$service_id" \
    -H "Authorization: Bearer $token")
if [ "$delete_response" != "200" ]; then
    echo "❌ Failed to delete service"
    docker-compose down
    exit 1
fi
echo "✅ Delete service passed"

# Cleanup
echo "Cleaning up..."
docker-compose down

echo "✅ All integration tests passed!"
