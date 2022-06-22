# Simple Bank
![img_1.png](readme_resource/docker_logo.png)
![img_2.png](readme_resource/postgreSQL_logo.png)
## a go/docker/postgres project to implement a simple bank server
![img.png](readme_resource/bank_logo.png)

It will provide APIs for the frontend to do following things:
* Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
* Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
* Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.

## Setup local development
### install tool
* Docker
  ```shell
  sudo apt-get update sudo apt-get install docker-ce docker-ce-cli containerd.io
  ```
* [TablePuls](https://tableplus.com/linux)
* Migrate
  ```shell
  curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz
  ```
* Sqlc 
  ```shell
  go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
  ```
* Gomock
  ```shell
  go install github.com/golang/mock/mockgen@v1.6.0  
  ```
* Viper 
  ```shell
  go get github.com/spf13/viper
  ```

### Setup infrastructure

* Build docker compose
  ```
  reference docker-compose
  ```
* Unit test with GitHub Action
  ```
  reference .github/workflow test.yml
  ```
  * Push image to Alibaba Cloud Container Registry
  ```
  reference .GitHub/workflow deploy.yml
  ```
 

## How to generate code
* Generate schema SQL file with DBML:
  ```shell
  make db_schema`
  ```  
* Generate SQL CRUD with sqlc:
  ```shell
  make sqlc
  ```
* Generate DB mock with gomock:
  ```shell
  make mock
  ```
* Create a new db migration:
  ```shell
  migrate create -ext sql -dir db/migration -seq <migration_name>  
  ```
## How to run
* Run server:
  ```shell
  make server
  ```
* Run test
  ```shell
  make test
  ```

