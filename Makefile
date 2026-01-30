build:
	cd frontend && yarn build
	ENV=prod go build -buildvcs=false -o ./bin/go-vite ./main.go

dev:
	cd frontend && yarn dev & sleep 3 && DEV_MODE=true go run ./cmd/server
