package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hrntknr/tendermint-test/app"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"
)

func main() {
	_app := app.NewApplication()
	srv, err := server.NewServer("tcp://0.0.0.0:46658", "socket", _app)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	srv.SetLogger(logger.With("module", "abci-server"))

	if err != nil {
		panic(err)
	}
	err = srv.Start()
	if err != nil {
		panic(err)
	}
	defer srv.Stop()

	for {
		time.Sleep(10 * time.Second)
		_app.DeliverTx([]byte(time.Now().String()))
		_app.Commit()
		res := _app.Query(types.RequestQuery{Path: "storage"})
		fmt.Printf("%s\n", string(res.Value))
		fmt.Printf("%d\n", res.Height)
	}
}
