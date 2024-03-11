CCM

DEBUG  RUN LOCAL

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/cloud-controller-manager ./cmd/main.go
./bin/cloud-controller-manager --cloud-provider=autok8s --v=5 --kubeconfig ~/.kube/config  --cloud-config=./test

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go run ./cmd/main.go --cloud-provider=autok8s --v=5 --kubeconfig ~/.kube/config  --cloud-config=./test
