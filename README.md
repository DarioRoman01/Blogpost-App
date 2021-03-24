# Contacts
a sample REST API for a blogpost application using Golang with Echo framework
and MongoDB to store the data.

## What is this Project?
this is a backend for a blogpost app. The users has all the CRUD operations with hers posts,
follow others users and comment the posts. Im currently develop a chat service app

## Why Echo?
For this project in particular i wanted a fast service and Golang with Echo comply with that.
Also the scalability that give me Echo, Binding and other featuers that made me choose Echo over other 
frameworks like Gin, Buffallo, Gorilla, Revel etc. 

## Why MongoDB?
The main reason is for learning MongoDB, also the relations are not very complex, MongoDB give me
flexibility in the documents and the queries that is great for this kind of service. And What better
way to learn MongoDB than using it in a project.


# Usage

## Requirements 
* Docker and Docker-compose

## Usage
```
$ export COMPOSE_FILE=docker-compose.yml

$ docker-compose build .

$ docker-compose up
```
**The app runs in the port :1323 and Mongo in the port :27017 so make sure you have that ports available or change the value in the files.**



