build:
	docker build -t marian/pinger .

test:
	docker run --rm -ti -v $(PWD):/config marian/pinger /config/config.yaml
