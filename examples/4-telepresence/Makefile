docker_image ?= moorara/example

build:
	go build -o app

swap:
	elepresence --swap-deployment example --run ./app

docker:
	docker image build --tag $(docker_image):latest .
	docker image push $(docker_image):latest

k8s-deploy:
	kubectl create -f kubernetes.yaml

k8s-delete:
	kubectl delete deployment example
	kubectl delete secret example

.PHONY: build swap docker k8s-deploy k8s-delete
