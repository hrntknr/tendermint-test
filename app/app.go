package app

import (
	"github.com/hrntknr/tendermint-test/code"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
)

type Application struct {
	types.BaseApplication

	storage []byte
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{Data: cmn.Fmt("{\"storage\":%s}", app.storage)}
}

func (app *Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (app *Application) DeliverTx(tx []byte) types.ResponseDeliverTx {
	app.storage = tx
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *Application) CheckTx(tx []byte) types.ResponseCheckTx {
	app.storage = tx
	return types.ResponseCheckTx{Code: code.CodeTypeOK}

}

func (app *Application) Commit() (resp types.ResponseCommit) {
	return types.ResponseCommit{Data: app.storage}
}

func (app *Application) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	switch reqQuery.Path {
	case "storage":
		return types.ResponseQuery{Value: app.storage}
	default:
		return types.ResponseQuery{Log: cmn.Fmt("Invalid query path. Expected storage, got %v", reqQuery.Path)}
	}
}
