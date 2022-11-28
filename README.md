# Social Network Project

> This project is an inspiration from the social business network Linkedin. Built on a monolithic architecture, it uses rest API, cache, sql and message-broker.
<p align="center">
  <img src="/assets/arch-diagram.png" height="440">
</p>

----

### Diagram of Database

<p align="center">
  <img src="/assets/db-diagram.png" height="400">
</p>

----


### Built With

* <img height="60" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" alt="gin-gonic"/>
* <img height="50" src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-plain.svg" alt="PostegreSQL"/>
* <img height="50" src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-plain.svg" alt="redis"/>
* <img height="50" src="https://assets.zabbix.com/img/brands/rabbitmq.svg" alt="redis"/>

### Pre-requisites
* go
  ```sh
  sudo snap install go --classic
  ```
* install [docker](https://docs.docker.com/install/)
* install [docker-compose](https://docs.docker.com/compose/install/)

  ```sh
  docker-compose -v
  docker version
  ```

### Installation
1. Clone the repository
   ```sh
   git clone https://github.com/victorvcruz/social_network_project.git
   ```
2. Install go packages
   ```sh
   go mod download
   ```
3. Run docker-compose
   ```sh
   sudo docker-compose up -d
   ```
  
## Usage
   
The project is a CRUD application of four operations

* Account
* Post
* Comment
* Interaction

To create operations define your **API_PORT** in `.env` file for use by endpoints

### To start execution
* run
   ```sh
   go run main.go
   ```


### Account Operations
- The `http://localhost:8080/accounts` endpoint is used for creating new accounts
- The `http://localhost:8080/accounts/auth` endpoint is used for creating new token
- The `http://localhost:8080/accounts/follows` endpoint is used to follow other accounts
- The `http://localhost:8080/accounts/following` endpoint is used to find which accounts are following
- The `http://localhost:8080/accounts/follower` endpoint is used to find which accounts follow it

#### :one: Request:

```console
curl -iX POST \
  --url 'http://localhost:8080/accounts' \
  --header 'Content-Type: application/json' \
  --data ' {
    "username": "jonh",
    "name": "Jonh Deep",
    "description": "my name is Jonh",
    "email": "jonh.deep@gmail.com",
    "password": "jonh1x3"
  }'
```
> This endpoint contains get, update and delete which need auth-token

#### :two: Request:

```console
curl -iX POST \
  --url 'http://localhost:8080/accounts/auth' \
  --header 'Content-Type: application/json' \
  --data ' {
    "email": "jonh.deep@gmail.com",
    "password": "jonh1x3"
  }'
```
> Token created: eyJhbGciOiJIUzI1NiIsInR5.example

#### :three: Request:

```console
curl -iX POST \
  --url 'http://localhost:8080/accounts/follows' \
  --header 'Content-Type: application/json' \ 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5.example' \
  --data ' {
    "id": "8216385e-730b-40a7-8fbd-a37889feac7d"
  }'
```
> Insert id to follow

### Post Operations
- The `http://localhost:8080/posts` endpoint is used for creating new posts
- The `http://localhost:8080/accounts/follows` endpoint is used for find post by accounts are following

#### :four: Request:

```console
curl -iX POST \
  --url 'http://localhost:8080/posts' \
  --header 'Content-Type: application/json' \ 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5.example' \
  --data ' {
    "content": "yes, my name is Jonh?"
  }'
```
> This endpoint contains get, update and delete

#### :five: Request:
```console
curl -X GET \
  --url 'http://localhost:8080/accounts/follows/posts?page=1' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5.example'
```
> Insert search pages

### Comment Operations

- The `http://localhost:8080/comments/:id` endpoint is used for creating new comments

#### :six: Request:

```console
curl -iX POST \
  --url 'http://localhost:8080/comments/52ba9bd3-e7e2-47fc-8ef4-99a24b32f888?comment_id=5e4a643c-befc-4854-bbe5-c7bbbb67ca2f' \
  --header 'Content-Type: application/json' \ 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5.example' \
  --data ' {
    "content": "no, your name not is Jonh."
  }'
```
> Insert post_id in params of endpoint and comment_id in query. This endpoint contains get, update and delete

### Interaction Operations

- The `http://localhost:8080/interaction` endpoint is used for creating new interaction

#### :seven: Request:

```console
curl -iX POST \
  --url 'http://localhost:8080/interaction' \
  --header 'Content-Type: application/json' \ 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5.example' \
  --data ' {
    "comment_id": "b166452d-9194-4b67-935b-4138428eb3d8",
    "type": "LIKE"
  }'
```
> This endpoint contains update and delete

## To stop execution
* run
  ```sh
  sudo docker-compose down
  ```
