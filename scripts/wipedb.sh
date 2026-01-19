if [ -f .env ]; then
    source .env 
fi

go run ./cmd/mock wipeall