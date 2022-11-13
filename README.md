
# Go-rengan
Go codebase with oil
## Stack
- Chi (net/http)
- MongoDB
- Opentelemetry
- Uptrace
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