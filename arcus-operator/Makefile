.PHONY: build

IMAGE=jam2in/arcus-operator:latest

all: build push

build:
		operator-sdk build $(IMAGE)

push:
		docker push $(IMAGE)
		kubectl delete deployment --all
		kubectl delete pod --all
		kubectl delete configmap --all
		kubectl delete service --all
		kubectl delete arcus --all

operator:
		kubectl apply -f deploy/operator.yaml

arcus:
		kubectl apply -f deploy/crds/jam2in_v1_arcus_cr.yaml

log:
		kubectl logs -f $(shell kubectl get pods | awk 'FNR==2{print $$1}')

generate:
		operator-sdk generate k8s
