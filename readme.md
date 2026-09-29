`docker build -t lets-go . &&
mkdir -p data &&
DB_NAME="$PWD/data/local_sales.db" RECORD_COUNT=100 go run ./cmd/setup &&
docker run --rm -p 8080:8080 -v "$PWD/data:/data" lets-go`