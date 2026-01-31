test:
	go test ./... -coverprofile=c.out

test_html:
	go tool cover -html=c.out
