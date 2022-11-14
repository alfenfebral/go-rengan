
# Go-rengan
Go codebase deep fry with oil
## Stack
- Chi (net/http)
- Wire (Dependency Injection)
- MongoDB
- Opentelemetry
- Uptrace
- RabbitMQ
## Run
Run Uptrace with docker compose
```bash
  docker compose up -d
```
Start the server using go run
```bash
  go run cmds/app/main.go
```
Start the server using [air](https://github.com/cosmtrek/air)
```bash
  make run
```
## Unit Test
Run Unit testing
```bash
  make test
```
Run Coverage
```bash
  make test/cover
```