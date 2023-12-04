

### DATABASE

How to create sqlite db

1) Using migration tool 

Steps:
1. Install migration tool 
`go install github.com/rubenv/sql-migrate/...@latest`
2. Create db
`make createdb`
3. Migrate up
`make migrateup`
4. (Optional) Migrate down
`make migratedown`

2) Manually using sql file

Steps:
1. Create db
`make createdb`
2. Init db
`make initdb`
3. (Optional) Delete db
`make deletedb`

### NOTE
Hard Dependencies

How to deploy sqlite