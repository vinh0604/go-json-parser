up:
	go mod download

test:
  go test -v ./...

build:
  go build -o ./bin/ ./cmd/json-parser/*.go

aider:
	ANTHROPIC_API_KEY=$(cat .anthropic_key) aider --sonnet