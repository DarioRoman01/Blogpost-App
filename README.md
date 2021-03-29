# Blogpost-app
a sample REST API for a blogpost application using Golang with Echo framework
and MongoDB to store the data.

![](https://github.com/Haizza1/my_first_repo/blob/master/golang.png)

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
* Docker and Docker-compose or golang installed on your system
* If you dont have Docker and want to use it with docker check this link [Install Docker](https://docs.docker.com/engine/install/)
* If you dont have go and want to use it with go check this link [Install Golang](https://golang.org/doc/install)

## Commands for Docker
```
$ export COMPOSE_FILE=docker-compose.yml

$ docker-compose build .

$ docker-compose up
```

## Commands for Go
```
$ go mod download

# No compile
$ go run main.go

# If you want to compile it
$ go build -o main.go

$ ./main
```
**The app runs in the port :1323 and Mongo in the port :27017 so make sure you have that ports available or change the value in the docker files.**



