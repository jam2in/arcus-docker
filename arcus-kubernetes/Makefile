# step0) make clean
# step1) make zk
# step2) make zknode
# step3) make memcached

zk:
	kubectl apply -f arcus-zookeeper.yaml

zknode:
	kubectl apply -f arcus-zookeeper-node.yaml

memcached:
	kubectl apply -f arcus-memcached.yaml

clean:
	kubectl delete -f arcus-zookeeper-node.yaml; kubectl delete -f arcus-memcached.yaml; kubectl delete -f arcus-zookeeper.yaml 
