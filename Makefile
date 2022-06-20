pull_api: 
	docker-compose down
	docker pull maslow123/keuanganku-users
	docker pull maslow123/keuanganku-pos
	docker pull maslow123/keuanganku-transactions
	docker pull maslow123/keuanganku-balance
	docker pull maslow123/keuanganku-apigateway

infratest: pull_api
	docker-compose up -d --force-recreate testdb
	echo Starting for db...
	sleep 15
	docker-compose up migratedb

runapi: infratest
	docker-compose up -d --force-recreate userapi
	docker-compose up -d --force-recreate posapi
	docker-compose up -d --force-recreate transactionapi
	docker-compose up -d --force-recreate balanceapi
	docker-compose up -d --force-recreate api-gateway

test: runapi
	cd users && go test -v ./... -coverprofile cover.out
	cd pos && go test -v ./... -coverprofile cover.out
	cd transactions && go test -v ./... -coverprofile cover.out
	cd balance && go test -v ./... -coverprofile cover.out
	cd api-gateway && go test -v ./... -coverprofile cover.out
	
	docker-compose down

resetdb:
	docker-compose down
	make infratest
