# Althea API

Open-source backend for efficiently serving ALTHEA data built using [Redis](https://github.com/redis/redis), [Fiber](https://github.com/gofiber/fiber) and [Go](https://github.com/golang/go). Built to minimize load on nodes to allow applications to scale better.

## Dependencies

- `golang 1.18` or above
- `redis 7.0` ([install here](https://redis.io/docs/getting-started/installation/))

## Quickstart

```bash
# clone repo
git clone https://github.com/chalabi2/althea-api

# create .env file and set variables:
nano .env
ALTHEA_MAINNET_RPC_URL = https://nodes.chandrastation.com/testnet/evm/althea/
ALTHEA_BACKUP_RPC_URL = https://nodes.chandrastation.com/testnet/evm/althea/
PORT = :3003
DB_HOST = localhost
DB_PORT = 6379
ALTHEA_MAINNET_GRPC_URL = <grpc url>
MULTICALL_ADDRESS=0xe9cBc7b381aA17C7574671e445830E3b90648368
QUERY_INTERVAL = 3

# build binary
cd althea-api
go build

# run redis
redis-server

# run binary
./althea-api
```

## Docker

Use docker compose:

`docker compose up -d`
