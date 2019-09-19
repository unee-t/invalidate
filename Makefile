PROJECT = invalidate
STACK_NAME ?= $(PROJECT)
AWS_REGION = ap-southeast-1
DEPLOY_S3_PREFIX = invalidate

.PHONY: deps clean build

deps:
	go get -u .
	go mod tidy

test: build
	sam local invoke -e codepipeline-event.json

clean:
	rm -rf invalidate

build:
	GOOS=linux GOARCH=amd64 go build

logs:
	AWS_PROFILE=uneet-dev sam logs -n invalidate -s yesterday

tail:
	AWS_PROFILE=uneet-dev sam logs -n invalidate -t

validate: template.yaml
	AWS_PROFILE=uneet-dev sam validate --template template.yaml

dev: build
	AWS_PROFILE=uneet-dev sam package --template-file template.yaml --s3-bucket dev-media-unee-t --s3-prefix $(DEPLOY_S3_PREFIX) --output-template-file packaged.yaml
	AWS_PROFILE=uneet-dev sam deploy --template-file ./packaged.yaml --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM
