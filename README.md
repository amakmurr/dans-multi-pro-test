## Description
This project purposed for code challenge in PT Dans Multi Pro for position Backend Developer (Golang)

## Pre-requisite
- Golang 1.19.x
- Docker

## How to run application
We will use docker to running application

### Run postgres
Execute the below command and make sure postgres server has been run
```shell
docker-compose up postgres
```

### Run application
In other terminal session, execute the below command to run the application
```shell
docker-compose up --build app
```

### Access application
Visit swagger on http://localhost:8000/static/swagger

You can login using the below credential
```shell
username: test-user
password: password
```

Get the access token through `/login` endpoint, and using it as `Bearer token` to access other endpoints.