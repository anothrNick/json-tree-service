rebuild: install
	docker-compose build --no-cache

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

remove:
	docker-compose rm -f

clean:
	docker-compose down --rmi all -v --remove-orphans
	rm json-tree-service
