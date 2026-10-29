build:
  go build -o bin/net-cli main.go

test:
  go test ./... -v
