serve:
	go run ./cmd/server \
		--port 8081 \
		--https-certificate-file ./certs/cert.pem \
		--https-certificate-key-file ./certs/cert-key.pem \
		--apiserver-host https://localhost:6443 \
		--apiserver-cacert ./certs/k8s-ca.pem \
		--oidc-client-id kubernetes \
		--oidc-client-secret DzWSrlGXQhmae5gxF3cls6kWf0eESENl \
		--oidc-issuer-url https://10.0.1.200/realms/dejavu

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./bin/kubectl-login  ./cmd/client
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./bin/kubectl-server ./cmd/server
