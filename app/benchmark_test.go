package app

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/evmos/ethermint/encoding"
)

func BenchmarkEthermintApp_ExportAppStateAndValidators(b *testing.B) {
	db := dbm.NewMemDB()
	app := NewEthermintApp(
		log.NewNopLogger(), 
		db, 
		nil, 
		true, 
		map[int64]bool{}, 
		DefaultNodeHome, 
		0, 
		encoding.MakeConfig(ModuleBasics), 
		simtestutil.EmptyAppOptions{},
		baseapp.SetChainID("ethermint_9000-1"),
	)

	genesisState := NewTestGenesisState(app.AppCodec())
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		b.Fatal(err)
	}

	// Initialize the chain
	app.InitChain(
		&abci.RequestInitChain{
			ChainId:       "ethermint_9000-1",
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: 1,
	})
	app.Commit()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Making a new app object with the db, so that initchain hasn't been called
		app2 := NewEthermintApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encoding.MakeConfig(ModuleBasics), simtestutil.EmptyAppOptions{})
		if _, err := app2.ExportAppStateAndValidators(false, []string{}, []string{}); err != nil {
			b.Fatal(err)
		}
	}
}
