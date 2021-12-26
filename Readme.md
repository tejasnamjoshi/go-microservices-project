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
 
 ## MYSQL Setup
 **NOTE**
 This is to be done only for the first time.

 - Once the containers are up and running, open your MYSQL client.
   - The hostname will be 127.0.0.1 and port will be 3306
   - The username and password will be root / root
 - Once inside the MYSQL client, import / run the db.sql file available in the root directory
   - This will setup a new database and will create the required tables.


## POSTMAN
 - There are 2 json postman files
   - Environment file
   - Collection file.
 - Open Postman and import these files.
 
## Running the Frontend NextJS app
**NOTE**
If you are doing this for the first time, run `yarn`.

 - If you want to run the frontend app, open a terminal and navigate to the next-todo-tailwind directory.
 - Run the following command `yarn dev`.
 - Head over to the following URL on your browser `http://localhost:3000`