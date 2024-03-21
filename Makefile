keys:
	echo "AKIA$(openssl rand 10 | base32)"
	echo $(openssl rand 40 | base64)

xks-proxy:
	docker build -t xks-proxy:latest .

infra: xks-proxy
	docker-compose up -d
