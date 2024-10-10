keys:
	@echo "Generating some random keys that AWS will use to authenticate to our xks-proxy service ..."
	@echo "ACCESS_KEY_ID (sigv4_access_key_id)         = AKIA$(shell openssl rand 10 | base32)"
	@echo "SECRET_ACCESS_KEY (sigv4_secret_access_key) = $(shell openssl rand 40 | base64)"
	@echo "Update the config/settings.toml file and start the xks-proxy service for change to take place!"

xks-proxy:
	docker build -t xks-proxy:latest .

infra: xks-proxy
	docker-compose up -d

keyswitch:
	cd KeySwitch && go build -o ../keyswitch .
