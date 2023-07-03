# MessengerService
Hi! This is Elvis and this repository is dedicated to developing a prototype of a messenger service that enables users to sign up, log in, log out, and create servers where other users can join. The messenger service aims to provide a platform for seamless communication and collaboration among users.


## Install ##
All that you need is [Golang](https://golang.org/). Once you run the application, it will expose a [target port](./config.json) on the host.
```
messengerservice-app-1  | [GIN-debug] Listening and serving HTTP on 0.0.0.0:4200
```

## Quickstart - Docker

To get started with the project in docker, follow these steps:

1. Clone the repository to your local machine:

   ```shell
   git clone https://github.com/Elvis-Benites-N/MessengerServiceGo.git
   ```

2. Move to the project directory:

   ```shell
   cd MessengerServiceGo
   ```

3. Build the Docker containers:

   ```shell
   docker-compose build
   ```

4. Start the Docker containers:

   ```shell
   docker-compose up
   ```

   This command will start the necessary services defined in the `docker-compose.yml` file.

5. Access the application:

   Once the Docker containers are up and running, you can access the application in your Postman software at `http://localhost:4200`.

That's it! You have successfully set up the project using Docker. If you encounter any issues or need further instructions, please refer to the project documentation or reach out to the project maintainers for assistance.

Make sure you have Docker and `docker-compose` installed on your machine before following these steps.

**Note:** If you need to customize any configurations or environment variables for your Docker setup, please refer to the `docker-compose.yml` file and make the necessary changes before running the `docker-compose build` and `docker-compose up` commands.

Feel free to modify this guide as needed to match the specific setup and requirements of your project.

I hope this helps you get started with Docker in your project! Let me know if you have any further questions.


## Quickstart - local

1. Clone the repository:
   ```shell
   git clone https://github.com/Elvis-Benites-N/MessengerServiceGo.git
   ```

2. Navigate to the project directory:
   ```shell
   cd MessengerServiceGo
   ```

3. Rename the `.env.example` file to `.env`:
   ```shell
   mv .env.example .env
   ```

4. Open the `.env` file and update the database configuration variables (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `DB_SSLMODE`) with your local PostgreSQL database settings.
 ```shell
   DB_HOST=database_host
   DB_USER=database_user
   DB_PASSWORD=database_password
   DB_NAME=database_name
   DB_PORT=database_port
   DB_SSLMODE=disable
   ```

5. Install the required dependencies:
   ```shell
   go get ./...
   ```

6. Run the project:
   ```shell
   go run local/main.go
   ```

   This will start the project and connect to the PostgreSQL database using the provided configuration.

7. Open Postman or any API testing tool and test the APIs using the appropriate endpoints and request methods.

Remember to ensure that you have PostgreSQL 13.11 (or an older version) and Go 1.20 (or an older version) installed on your local machine before running the project.

Note: It's recommended to have a clean and separate development environment for running the project locally, as it may have different dependencies and configuration compared to the production environment.

## USER - API
Send HTTP requests to user:
  * `POST /signup`: create a new user
  * `POST /login`: login with existing credentials
  * `GET /logout`: logout of session

E.g. A basic chat /signup POST request could look as follows:
```json
{
	"username" : "elvisbenites",
    "email" : "ebnbenites@gmail.com",
    "password" : "admin123"
}
```
If successful, the server will respond with HTTP code 200 and the newly created user:
```json
{
    "id": "1",
    "username": "elvisbenites",
    "email": "ebnbenites@gmail.com"
}
```
E.g. A basic chat /login POST request could look as follows:
```json
{
    "email" : "ebnbenites@gmail.com",
    "password" : "admin123"
}
```
If successful, the server will respond with HTTP code 200 and the user will be able to login:
```json
{
    "id": "1",
    "username": "elvisbenites"
}
```
E.g. If succesful a basic chat /logout GET request could look as follows, the server will respond with HTTP code 200 and the user will be able to logout:
```json
{
    "message": "Logout successful"
}
```

## SERVER - API
Send HTTP requests to `/ws` for server API:
  * `POST /createServer`: create a new chat server
  * `GET /joinServer/:serverId`: join to a chat server by its ID
  * `GET /getServers`: list all the servers created
  * `GET /getClients/:serverId`: list all the users connected to an server

E.g. A basic chat /createServer POST request could look as follows. if this process is successful, the server will respond with HTTP code 200 and shows the data entered:
```json
{
	"id": "1",
    "name": "GoDevelopers"
}
```
E.g. A basic chat /joinServer/:serverId request could look as follows in Postman using WebSocket:
 ```shell
    ws://localhost:4200/ws/joinServer/1?userId=2&username=user2
   ```

E.g. If succesful a basic chat /getServers GET request could look as follows, the server will respond with HTTP code 200:
```json
[
    {
        "id": "1",
        "name": "GoDeverlopers"
    }
]
```

E.g. If succesful a basic chat /getClients/:serverId GET request could look as follows, the server will respond with HTTP code 200:
```json
[
    {
        "id": "2",
        "username": "user2"
    },
    {
        "id": "1",
        "username": "elvisbenites"
    }
]
```

## EXPECTED RESULT

### Architecture in Postman
The expected results in conjunction with the application of the Golang language are shown below.

![Initial Postman Architecture](/assets/Postman-Workspace.png)

### Dashboard in Postman

![Initial Postman Architecture](/assets/Postman-Dashboard.png)
