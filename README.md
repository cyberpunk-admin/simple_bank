# Simple Bank
![img_1.png](img_1.png)
![img_2.png](img_2.png)
## a go/docker/postgres project to implement a simple bank server
![img.png](img.png)

It will provide APIs for the frontend to do following things:
* Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
* Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
* Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.

### install docker
` sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io`

* #### pull postgresDB image
`docker pull postgres:13-alpine`

### sqlc generates fully-type safe idiomatic Go code from SQL.
`sudo snap install sqlc`

### golong-migrate for databsase migration

`$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz
`
### gomock for unit test
`go install github.com/golang/mock/mockgen@v1.6.0`

### Viper for configuration manage
`go get github.com/spf13/viper`

## Setup infrastructure
reference Makefile

## Build docker compose
reference docker-compose

## unit test with GitHub Action
reference .github/workflow test.yml

## Push image to Alibaba Cloud Container Registry
reference .GitHub/workflow deploy.yml

