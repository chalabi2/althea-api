# Althea API

Open-source backend for efficiently serving Canto data built using [Redis](https://github.com/redis/redis), [Fiber](https://github.com/gofiber/fiber) and [Go](https://github.com/golang/go). Built to minimize load on nodes to allow applications to scale better.

## Dependencies

- `golang 1.18` or above
- `redis 7.0` ([install here](https://redis.io/docs/getting-started/installation/))

## Quickstart

```bash
# clone repo
git clone git@github.com:Plex-Engineer/canto-api.git

# create .env file and set variables:
nano .env
CANTO_MAINNET_RPC_URL = https://nodes.chandrastation.com/testnet/evm/althea/
CANTO_BACKUP_RPC_URL = https://nodes.chandrastation.com/testnet/evm/althea/
PORT = :3003
DB_HOST = localhost
DB_PORT = 6379
CANTO_MAINNET_GRPC_URL = <grpc url>
MULTICALL_ADDRESS=0xe9cBc7b381aA17C7574671e445830E3b90648368
QUERY_INTERVAL = 3

# build binary
cd canto-api
go build

# run redis
redis-server

# run binary
./canto-api
```

## Docker

Use docker compose:

`docker compose up -d`
