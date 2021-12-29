# Todo App using golang microservices

## Introduction

This app contains 2 microservices 
 - Auth
 - Todo

As the names suggest, Auth is used for authentication whereas todo deals with CRUD todo operations.
For inter service communication, NATS is used. It also uses mysql for storing data.
MYSQL and NATS are both setup in a docker environment.

## Setup Instructions

 - Make sure you have Docker installed on your machine.
 - Make sure .env is set to the correct values.
 - Open a terminal in this directory and run the following command `docker-compose up -d`
 - This creates the following docker containers
   - Mysql
   - NATS
   - nats-service
   - auth-service
   - todo-service
   - next-todo
 - Once the containers are up and running, you can start using the frontend app or postman.
 - The `schema.sql` file from the mysql initializes the database with the correct database and tables.

## Running the Frontend NextJS app
 - Head over to the following URL on your browser `http://localhost:3000`

## POSTMAN
 - There are 2 json postman files
   - Environment file
   - Collection file.
 - Open Postman and import these files.

## Libraries and Frameworks used
 - go-chi for routing
 - Validator for validations
 - NATS for inter-service communication
 - jwt-go for JWT based authentication.
 - SQLX for DB interaction (mysql)
 - NextJS for the Frontend
 - Docker for containerizing all services ( BE + FE )
