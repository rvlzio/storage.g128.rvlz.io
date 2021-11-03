build:
	docker compose build

create-database:
	docker compose run --rm --no-deps management-tool python manager.py create-database ${name}

drop-database:
	docker compose run --rm --no-deps management-tool python manager.py drop-database ${name}

create-database-user:
	docker compose run --rm --no-deps management-tool python manager.py create-database-user ${database} ${username}

drop-database-user:
	docker compose run --rm --no-deps management-tool python manager.py drop-database-user ${database} ${username}

create-migration:
	docker compose run --rm --no-deps migrations migrate create -ext sql -dir versions -seq ${name}

run-migrations:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions up'

reverse-migrations:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions down'

run-one-migration:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions up 1'

reverse-one-migration:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions down 1'

migration-version:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions version'

drop-migrations:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions drop'

go-to-migration:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions goto ${version}'

force-migration:
	docker compose run --rm --no-deps migrations bash -c 'migrate -database $${DATABASE_URI} -path ./versions force ${version}'
