build:
	go build -o timer-cli cmd/timer-cli/main.go

test:
	go test ./... -coverprofile=c.out

test_html:
	go tool cover -html=c.out
