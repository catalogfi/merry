package main

import (
	"github.com/catalogfi/merry"
	"github.com/catalogfi/merry/jsonrpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	lgr, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	m := merry.Default()
	panic(jsonrpc.NewServer(lgr, gin.New(), jsonrpc.NewMerryHandler(&m)).Start(2201))
}
