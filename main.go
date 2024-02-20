package main

import (
	"context"

	"althea-api/config"
	// cqe "althea-api/queryengine/contracts"
	nqe "althea-api/queryengine/native"
	re "althea-api/requestengine"
)

func main() {
	config.NewConfig()
	ctx := context.Background()
	// go cqe.Run(ctx) // run contract query engine
	go nqe.Run(ctx) // run native query engine
	re.Run(ctx)     // run request engine
}
