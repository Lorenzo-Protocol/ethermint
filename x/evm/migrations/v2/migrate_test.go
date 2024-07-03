package v2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/evmos/ethermint/encoding"

	"github.com/evmos/ethermint/app"
	v2 "github.com/evmos/ethermint/x/evm/migrations/v2"
	v2types "github.com/evmos/ethermint/x/evm/migrations/v2/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func TestMigrateStore(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	feemarketKey := storetypes.NewKVStoreKey(types.StoreKey)
	tFeeMarketKey := storetypes.NewTransientStoreKey(fmt.Sprintf("%s_test", types.StoreKey))
	ctx := testutil.DefaultContext(feemarketKey, tFeeMarketKey)
	paramstore := paramtypes.NewSubspace(
		encCfg.Codec, encCfg.Amino, feemarketKey, tFeeMarketKey, "evm",
	).WithKeyTable(v2types.ParamKeyTable())

	params := v2types.DefaultParams()
	paramstore.SetParamSet(ctx, &params)

	require.Panics(t, func() {
		var result bool
		paramstore.Get(ctx, types.ParamStoreKeyAllowUnprotectedTxs, &result)
	})

	paramstore = paramtypes.NewSubspace(
		encCfg.Codec, encCfg.Amino, feemarketKey, tFeeMarketKey, "evm",
	).WithKeyTable(types.ParamKeyTable())
	err := v2.MigrateStore(ctx, &paramstore)
	require.NoError(t, err)

	var result bool
	paramstore.Get(ctx, types.ParamStoreKeyAllowUnprotectedTxs, &result)
	require.Equal(t, types.DefaultAllowUnprotectedTxs, result)
}
