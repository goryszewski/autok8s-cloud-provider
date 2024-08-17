update:
	@echo "[MAKE] Update client"
	go get -u github.com/goryszewski/libvirtApi-client

updateall:
	@echo "[MAKE] Update all modules"
	go get -u all
	go mod tidy

build:
	@echo "[MAKE] Build app"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/cloud-controller-manager ./cmd/main.go
