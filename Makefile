BINARY_NAME = main.out
DOCKER_USERNAME = atriiy
APPLICATION_NAME = pool
GIT_HASH ?= $(shell git log --format="%h" -n 1)


build:
	docker build --tag ${DOCKER_USERNAME}/${APPLICATION_NAME}:${GIT_HASH} .

push:
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:${GIT_HASH}

release:
	docker pull ${DOCKER_USERNAME}/${APPLICATION_NAME}:${GIT_HASH}
	docker tag  ${DOCKER_USERNAME}/${APPLICATION_NAME}:${GIT_HASH} ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest

start:
	docker compose up

test:
	go test -v ./...

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

clean:
	go clean
	rm -rf ./out
	rm -rf ./cache
	rm -f *.out

clean-docker:
	docker rm dockerdev
	docker rm dockerdev-db
	docker rmi pool-crawler
