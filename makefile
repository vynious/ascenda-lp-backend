.PHONY: build

S3_BUCKET=aws-sam-cli-managed-default-samclisourcebucket-czurtb144gdf/itsa-api
STACK_NAME=itsa-api

build:
	sam build

run:
	DOCKER_HOST=unix://${HOME}/.docker/run/docker.sock sam local start-api

build-run:
	sam build && DOCKER_HOST=unix://${HOME}/.docker/run/docker.sock sam local start-api

deploy:
	sam build && sam deploy

deploy-package:
	sam package \
	--template-file template.yaml \
	--output-template-file package.yml \
	--s3-bucket ${S3_BUCKET}

	sam deploy \
	--template-file package.yml \
	--stack-name ${STACK_NAME} \
	--capabilities CAPABILITY_IAM

package:
	sam package \
	--template-file template.yaml \
	--output-template-file package.yml \
	--s3-bucket ${S3_BUCKET}

teardown:
	sam delete --stack-name ${STACK_NAME}
