keys:
	echo "AKIA$(openssl rand 10 | base32)"
	echo $(openssl rand 40 | base64)

xks-proxy:
	docker build -t xks-proxy:latest .

infra: xks-proxy
	docker-compose up -d

test:
	export XKS_PROXY_HOST="localhost:8080"
	export URI_PREFIX="aws-redteam-kit"
	export SIGV4_ACCESS_KEY_ID="AKIAHPRO52CGOA4VJPU4"
	export SIGV4_SECRET_ACCESS_KEY="yTBJ7hqzT7TTtqCHGqFjJQDTJaZolkCz5i5h5TUSwHNUK1glZ6rMpQ=="
	export KEY_ID="thekey"
	export REGION="ap-south-1"
	cd aws-kms-xksproxy-test-client/ && ./test-xks-proxy