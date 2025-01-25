## Local setup

1. Clone the repository and navigate to it
```
git clone https://github.com/Glenn-Chiang/forum-api
cd forum-api
```

2. Set up a PostgreSQL database with Docker
```
docker-compose up -d
```

3. Run the server
```
export ENV=development
go run cmd/main.go
```
