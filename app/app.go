package app

import (
	"encoding/binary"

	"github.com/hrntknr/tendermint-test/code"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
)

type Application struct {
	types.BaseApplication

	Count uint64
}

func NewApplication() *Application {
	return &Application{}

}

func (app *Application) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{Data: cmn.Fmt("{\"count\":%d}", app.Count)}
}

func (app *Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (app *Application) DeliverTx(value []byte) types.ResponseDeliverTx {
	app.Count++
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *Application) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *Application) Commit() (resp types.ResponseCommit) {
	if app.Count == 0 {
		return types.ResponseCommit{}
	}
	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, app.Count)
	return types.ResponseCommit{Data: hash}
}

func (app *Application) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	switch reqQuery.Path {
	case "count":
		return types.ResponseQuery{Value: []byte(cmn.Fmt("%d", app.Count))}
	default:
		return types.ResponseQuery{Log: cmn.Fmt("Invalid query path. Expected count, got %s", reqQuery.Path)}
	}
}
