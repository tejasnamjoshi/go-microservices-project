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
 - This creates 2 docker containers
   - Mysql
   - NATS
 - Once the containers are up and running, open your MYSQL client.
   - The hostname will be 127.0.0.1 and port will be 3306
   - The username and password will be root / root
 - Once inside the MYSQL client, import / run the db.sql file available in the root directory
   - This will setup a new database and will create the required tables.
 - Now it is time to run the services.
 - To build all the services, run the following command `make buildAll`
   - This will create a build directory and create executables in it.
   - It will also copy over the .env file
 - Navigate to the build directory using `cd build`
 - Now there are 2 ways to run the services
   - Open 3 terminals and run 1 service in each
     - A single service can be run like so `./auth` and so on.
   - This method allows you to run multiple services with a single command
     - `./nats & ./auth & ./todo`
     - The drawback here is that the first 2 services run in the background i.e. they cannot be stopped with a simple `ctrl + c` command.
     - To stop all services, first press `ctrl + c`. This will stop the todo service. 
     - Next, you will see 2 Numbers just after you run the above command. These will be a 4 digit number.
     - Run the following commands
       - `kill -9 <number1> & kill -9 <number2>`
     - This will kill the nats and auth service.