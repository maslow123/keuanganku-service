## Tools that need to be installed
- [Docker](https://www.docker.com/)
- [Make](https://community.chocolatey.org/packages/make)
- [Go](https://go.dev/)

## Things that must be considered
***Make sure no postgres service is running in the background.**

## How to run the application?
- Clone this repository
- Open your terminal / cmd, and type the command on below:
    ```
    make runapi
    ```
- Make sure the service is running properly, as follows:
    ```
    $ docker ps
    CONTAINER ID   IMAGE                                      COMMAND                  CREATED          STATUS          PORTS                      NAMES
    189722a2eba6   maslow123/keuanganku-apigateway:latest     "./main"                 17 minutes ago   Up 17 minutes   0.0.0.0:8000->8000/tcp     api-gateway
    acb20813364b   maslow123/keuanganku-balance:latest        "./main"                 5 hours ago      Up 5 hours      0.0.0.0:50054->50054/tcp   balanceapi
    320dd23cd40f   maslow123/keuanganku-transactions:latest   "./main"                 5 hours ago      Up 5 hours      0.0.0.0:50053->50053/tcp   transactionapi
    f92c965140c9   maslow123/keuanganku-pos:latest            "./main"                 5 hours ago      Up 5 hours      0.0.0.0:50052->50052/tcp   posapi
    36576c823046   maslow123/keuanganku-users:latest          "./main"                 5 hours ago      Up 5 hours      0.0.0.0:50051->50051/tcp   userapi
    ba4ca701b577   postgres:latest                            "docker-entrypoint.s…"   5 hours ago      Up 5 hours      0.0.0.0:5433->5432/tcp     testdb
    ```
- If all services are running well, then import `Keuanganku.postman_collection.json` into POSTMAN
- Finish

## How to run unit testing of each service?
- You just need to enter the command
    ```make test```
- Finish.
## Shut down all services
- Make sure you're on root folder (keuanganku-service), and enter the command
```$ docker-compose down```
- Finish.
