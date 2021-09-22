## Efishery technical test: REST API for auth service and fetch commodity service
### How to start
#### 1. Start the dev environment(mysql) using this command
```
make start-dev
```
or use docker-compose command
```
docker-compose up -d
```
notes: docker and docker-compose requred to run the script, https://docs.docker.com/compose/install/



#### 2. Create table user and role on database
```
docker exec -i   efishery_db_1 mysql -uroot -proot test_db < auth.sql
```
or you can usq the DDL in the file.
(add sudo before the command if docker not yet managed as non-root user)

#### 3. run the app
```
go run main.go http
```

#### notes:

- test the api can be done using the postman collection in this directory
