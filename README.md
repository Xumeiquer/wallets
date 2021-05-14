# WALLETS

This is a small tool for tracking wallets for indexed funds.

## Download

WALLETS is avaliable to download [here](https://github.com/Xumeiquer/wallets/releases).

## Build

1. Download and install [Go](https://golang.org/)
1. Download and install [Task](https://github.com/go-task/task)
1. Clone this repository
1. Run `task clean && task`

## Running WALLETS

1. Run `./dist/server serve`
1. Navigate to [127.0.0.1:9090](http://127.0.0.1:9090)

WALLETS uses a small database to store your wallet information. That database is a file called `wallets.db` by default and it is placed next to the server also by default. You can customized it by using the server options.

### Server options

The server has some options that helps you to customize the execution.

```txt
./dist/server --help
Wallet tracking webapp & server :: v0.0.1

Usage:
  serve [flags]

Flags:
      --bind string       ip to bind the server (default "127.0.0.1")
      --database string   path where resides the database (default "wallet.db")
      --debug             debug
  -h, --help              help for serve
      --port int          port to listen connection from (default 9090)
```

