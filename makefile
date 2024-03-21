.PHONY: build

build:
	sam build

run:
	DOCKER_HOST=unix://${HOME}/.docker/run/docker.sock sam local start-api

build-run:
	sam build && DOCKER_HOST=unix://${HOME}/.docker/run/docker.sock sam local start-api

deploy:
	sam build && sam deploy

teardown:
	sam delete --stack-name ${STACK_NAME}
