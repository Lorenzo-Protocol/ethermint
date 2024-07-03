package v2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/evmos/ethermint/encoding"

	"github.com/evmos/ethermint/app"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	v1 "github.com/evmos/ethermint/x/feemarket/migrations/v1"
	v2 "github.com/evmos/ethermint/x/feemarket/migrations/v2"
	v2types "github.com/evmos/ethermint/x/feemarket/migrations/v2/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

func TestMigrateStore(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	feemarketKey := storetypes.NewKVStoreKey(feemarkettypes.StoreKey)
	tFeeMarketKey := storetypes.NewTransientStoreKey(fmt.Sprintf("%s_test", feemarkettypes.StoreKey))
	ctx := testutil.DefaultContext(feemarketKey, tFeeMarketKey)
	paramstore := paramtypes.NewSubspace(
		encCfg.Codec, encCfg.Amino, feemarketKey, tFeeMarketKey, "feemarket",
	)

	paramstore = paramstore.WithKeyTable(v2types.ParamKeyTable())
	require.True(t, paramstore.HasKeyTable())

	params := v2types.DefaultParams()
	paramstore.SetParamSet(ctx, &params)

	// check that the fee market is not nil
	err := v2.MigrateStore(ctx, &paramstore, feemarketKey)
	require.NoError(t, err)
	require.False(t, ctx.KVStore(feemarketKey).Has(v2.KeyPrefixBaseFeeV1))

	fmKeeper := feemarketkeeper.NewKeeper(
		encCfg.Codec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		feemarketKey,
		tFeeMarketKey,
		paramstore,
	)

	expParams := fmKeeper.GetParams(ctx)
	require.False(t, expParams.BaseFee.IsNil())

	baseFee := fmKeeper.GetBaseFee(ctx)
	require.NotNil(t, baseFee)

	require.Equal(t, baseFee.Int64(), params.BaseFee.Int64())
}

func TestMigrateJSON(t *testing.T) {
	rawJson := `{
		"base_fee": "669921875",
		"block_gas": "0",
		"params": {
			"base_fee_change_denominator": 8,
			"elasticity_multiplier": 2,
			"enable_height": "0",
			"initial_base_fee": "1000000000",
			"no_base_fee": false
		}
  }`
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	var genState v1.GenesisState
	err := encCfg.Codec.UnmarshalJSON([]byte(rawJson), &genState)
	require.NoError(t, err)

	migratedGenState := v2.MigrateJSON(genState)

	require.Equal(t, int64(669921875), migratedGenState.Params.BaseFee.Int64())
}
