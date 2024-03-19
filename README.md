## This repository is the basic layout for creating a telegram bot in golang using [kbgod/illuminate](https://github.com/kbgod/illuminate)

### Quick start

Before you start, you need to create a bot in the telegram and get the token. After that, you need to create postgres database.

**Init env file**
```shell
cp .env.example .env
```

**Install dependencies**
```shell
go mod tidy
```

**Run**
```shell
go run cmd/bot/main.go
```

**Test FSM**
```
/set_my_name
```