serve:
	go run ./cmd/server \
		--port 8081 \
		--https-certificate-file ./certs/cert.crt \
		--https-certificate-key-file ./certs/cert.key \
		--apiserver-host https://localhost:6443 \
		--oidc-client-id kubernetes \
		--oidc-client-secret J6WcVVWR2Dz0x4A0bIZRGdZpO1Kt8J5t \
		--oidc-issuer-url http://localhost:8080/realms/k8s

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./bin/kubectl-login  ./cmd/client
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./bin/kubectl-server ./cmd/server
