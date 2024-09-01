# Test Task Go Json encryptor/decryptor
Service that generates encrypted JSON file and decrypt JSON it and store. It depends on current mode.

## Features
- Generate random JSON
- Encrypt/decrypt using Fernet
- Store data from JSON using scraping keys

## How to run
To run a service first you need to:

### 1) Configure .env
```
# SERVER
SERVER_MODE="ENCRYPTION"
SERVER_HTTP_ADDR="0.0.0.0:8080"

# DATA
DATA_FOLDER="path-to-folder"
DATA_CRYPTO_KEY="secret-key"

# DATABASE
DATA_DATABASE_HOST=database-host
DATA_DATABASE_PORT=5432
DATA_DATABASE_NAME=database-name
DATA_DATABASE_USER=database-user
DATA_DATABASE_PASSWORD=database-password
```

### 2) Run
* Run docker service
* Run and build docker containers
```bash
$ docker compose --env-file .env up -d --build
```
* Call API
```bash
$ curl localhost:8080/run
```
Each API call generate file or collect data from it. 
It depends on current server mode

## Test
* Tests located in **internal/tests**
* Look into **usecase_*_test.go** files
* Before running test you should run docker service
* To run test
```bash
make test
```

## Check encryption and decryption
### Encryption
* In **.env** file change **SERVER_MODE**
```
SERVER_MODE="ENCRYPTION"
...
```
* Rerun docker container
```bash
$ docker compose --env-file .env up -d
```
* Call API
```bash
$ curl localhost:8080/run
```
* See response
```
Successfully generated JSON file
```
* Check generated JSON file in **data** folder in the root of the project

Content of the folder should be something like data/data-1234567890.json

### Decryption
* In **.env** file change **SERVER_MODE**
```
SERVER_MODE="DECRYPTION"
...
```
* Rerun docker container
```bash
$ docker compose --env-file .env up -d
```
* Call API
```bash
$ curl localhost:8080/run
```
* See response
```
Successfully collected dat
```
* Check collected JSON file in **data/collected** folder in the root of the project

Content of the folder should be something like data/collected/data-1234567890.json