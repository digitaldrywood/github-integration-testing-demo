#!/bin/bash

# Setup script for local integration testing

set -e

echo "🚀 Setting up local testing environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    echo "Please install Docker from https://docs.docker.com/get-docker/"
    exit 1
fi

# Start MinIO for S3 testing
echo -e "${YELLOW}Starting MinIO (S3-compatible storage)...${NC}"
docker run -d \
    --name minio-test \
    -p 9000:9000 \
    -p 9001:9001 \
    -e MINIO_ROOT_USER=minioadmin \
    -e MINIO_ROOT_PASSWORD=minioadmin \
    minio/minio server /data --console-address ":9001"

# Wait for MinIO to be ready
echo "Waiting for MinIO to be ready..."
for i in {1..30}; do
    if curl -s http://localhost:9000/minio/health/ready > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MinIO is ready${NC}"
        break
    fi
    sleep 1
done

# Create test bucket
echo "Creating test bucket..."
docker run --rm \
    --network host \
    -e AWS_ACCESS_KEY_ID=minioadmin \
    -e AWS_SECRET_ACCESS_KEY=minioadmin \
    amazon/aws-cli \
    --endpoint-url http://localhost:9000 \
    s3 mb s3://test-bucket

# Start PostgreSQL for database testing
echo -e "${YELLOW}Starting PostgreSQL...${NC}"
docker run -d \
    --name postgres-test \
    -p 5432:5432 \
    -e POSTGRES_USER=testuser \
    -e POSTGRES_PASSWORD=testpass \
    -e POSTGRES_DB=testdb \
    postgres:15

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec postgres-test pg_isready -U testuser > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL is ready${NC}"
        break
    fi
    sleep 1
done

# Export environment variables
echo -e "${YELLOW}Setting environment variables...${NC}"
cat > .env.test << EOF
# S3 Configuration
export AWS_ACCESS_KEY_ID=minioadmin
export AWS_SECRET_ACCESS_KEY=minioadmin
export AWS_REGION=us-east-1
export TEST_S3_BUCKET=test-bucket
export S3_ENDPOINT=http://localhost:9000

# Database Configuration
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=testuser
export DB_PASSWORD=testpass
export DB_NAME=testdb

# API Configuration (mock)
export API_ENDPOINT=http://localhost:8080
export API_KEY=test-api-key
EOF

echo -e "${GREEN}✅ Local testing environment is ready!${NC}"
echo ""
echo "To use the environment variables:"
echo "  source .env.test"
echo ""
echo "To run tests:"
echo "  # Unit tests"
echo "  go test ./..."
echo ""
echo "  # Integration tests with S3"
echo "  go test -tags=integration ./tests -s3"
echo ""
echo "  # All integration tests"
echo "  go test -tags=integration ./tests -s3 -database -api"
echo ""
echo "To stop and cleanup:"
echo "  docker stop minio-test postgres-test"
echo "  docker rm minio-test postgres-test"