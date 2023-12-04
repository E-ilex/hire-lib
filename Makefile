ifdef port
app_port = $(port)
else
app_port := 3000
endif

createdb:
	mkdir data && touch data/qa.sqlite3

removedb:
	rm -r data

migrateup:
	sql-migrate up -config=db/dbconfig.yml

migratedown:
	sql-migrate down -config=db/dbconfig.yml

migratestatus:
	sql-migrate status -config=db/dbconfig.yml

initdb:
	sqlite3 -init db/scripts/create.sql data/qa.sqlite3 .quit

deletedb:
	sqlite3 data/qa.sqlite3 < db/scripts/delete.sql

build:
	docker build -t toggl-server:latest .

run:
	docker run --name toggl-server -p $(app_port):$(app_port) toggl-server:latest -v /data/qa.sqlite3:/data/qa.sqlite3

start:
	docker start -i toggl-server

.PHONY: createdb migrateup migratedown migratestatus initdb deletedb build run start
