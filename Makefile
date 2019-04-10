build:
	docker build -t quay.io/marian/pinger .

test:
	docker run --rm -ti -v $(PWD):/config quay.io/marian/pinger /config/config.yaml
