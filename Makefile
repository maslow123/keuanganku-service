pull_api: 
	docker-compose down
	docker pull maslow123/keuanganku-users
	docker pull maslow123/keuanganku-apigateway

infratest: 
	docker-compose up -d --force-recreate testdb
	echo Starting for db...
	# sleep 15
	docker-compose up migratedb

runapi:
	docker-compose up -d --force-recreate userapi
	docker-compose up -d --force-recreate api-gateway

resetdb:
	docker-compose down
	make infratest
