pull_api: 
	docker-compose down
	# docker pull maslow123/keuanganku

infratest: pull_api
	docker-compose up -d --force-recreate testdb
	echo Starting for db...
	# sleep 15
	docker-compose up migratedb

runapi:
	docker-compose up -d --force-recreate api
	docker-compose up -d --force-recreate api-gateway