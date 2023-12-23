# Unattended Programming Test: The Test Signer

| Method | URL Pattern       | Handler                | Action                                                                         |
|--------|-------------------|------------------------|--------------------------------------------------------------------------------|
| GET    | /ping             | pingHandler            |                                                                                |
| POST   | /signature        | createSignatureHandler | Accepts a user JWT, questions and answers, and creates and returns a signature |
| POST   | /signature/verify | verifySignatureHandler | Accepts a user JWT and signature and returns                                   |

## Running the server locally

1. Make sure that you have docker-compose and make installed on local machine
2. Run `make start` to start the server for the first time, this will pull all the required docker containers and run database migrations
3. Run `make tests` to run the application tests

## Decisions

* In order to implement proper lock mechanism Redis is added as a dependency to the project.