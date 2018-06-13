package app

import (
	"encoding/binary"
	"encoding/json"

	"github.com/hrntknr/tendermint-test/code"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

var (
	stateKey = []byte("stateKey")
)

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

// var _ types.Application = (*Application)(nil)

type Application struct {
	types.BaseApplication

	state State
}

func NewApplication() *Application {
	state := loadState(dbm.NewMemDB())
	return &Application{state: state}
}

func (app *Application) Info(req types.RequestInfo) types.ResponseInfo {
	value := app.state.db.Get([]byte("storage"))
	return types.ResponseInfo{Data: cmn.Fmt("{\"storage\":%s}", value)}
}

func (app *Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (app *Application) DeliverTx(value []byte) types.ResponseDeliverTx {
	app.state.db.Set([]byte("storage"), value)
	app.state.Size++
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *Application) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *Application) Commit() (resp types.ResponseCommit) {
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height++

	saveState(app.state)
	return types.ResponseCommit{Data: appHash}
}

func (app *Application) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	switch reqQuery.Path {
	case "storage":
		value := app.state.db.Get([]byte("storage"))
		return types.ResponseQuery{Value: value}
	default:
		return types.ResponseQuery{Log: cmn.Fmt("Invalid query path. Expected storage, got %v", reqQuery.Path)}
	}
}
