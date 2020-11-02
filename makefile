start:
	docker stack deploy -c docker-compose.yaml -c build/mysql/docker-compose.yaml gcl

build:
	sh scripts/init.sh
