build:
	cross-env GOOS=linux GOARCH=amd64 go build -o bin/bot ./cmd/bot/main.go
deploy:
	scp bin/bot root@fruitcoin:/root/fruitcoin/