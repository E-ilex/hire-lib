### Scope of HW

The solution should meet the basic and the first two bonus requirements. Unfortunately, unit tests were not added, due to time restriction.

### Setup DB

How to create sqlite db

1. Using migration tool 
    1. Install migration tool 
`go install github.com/rubenv/sql-migrate/...@latest`
    2. Create db
`make createdb`
    3. Migrate up
`make migrateup`
    4. (Optional) Migrate down 
    `make migratedown`



2. Manually using sql file
    1. Create db
`make createdb`
    2. Init db
`make initdb`
    3. (Optional) Delete db
`make deletedb`

### Run Service

How to run service

1. Locally
    1. run main.go

2. ~~Using Docker~~ 
    1. I faced some problems with volume mounting after migrating from docker-compose to Dockerfile

### Curl

1. Create question

```curl --location --request POST 'http://localhost:3000/v1/question' \
--header 'Content-Type: application/json' \
--data-raw '{
  "body": "Where does the sun rise?",
  "options": [
    {
      "body": "East",
      "correct": true
    },
    {
      "body": "West",
      "correct": false
    }
  ]
}'
```

2. Select questions

```
curl --location --request GET 'http://localhost:3000/v1/question?limit=5&last_id=2'
```

3. Update question

```
curl --location --request PUT 'http://localhost:3000/v1/question/1' \
--header 'Content-Type: application/json' \
--data-raw '{
  "body": "Where does the sun set?",
  "options": [
    {
      "body": "East",
      "correct": false
    },
    {
      "body": "West",
      "correct": true
    }
  ]
}'
```

4. Delete question

```
curl --location --request DELETE 'http://localhost:3000/v1/question/1'
```