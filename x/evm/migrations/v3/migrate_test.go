package v3_test

import (
	"fmt"
	"testing"

	v3 "github.com/evmos/ethermint/x/evm/migrations/v3"

	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/evmos/ethermint/encoding"

	"github.com/evmos/ethermint/app"
	v3types "github.com/evmos/ethermint/x/evm/migrations/v3/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func TestMigrateStore(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	evmKey := storetypes.NewKVStoreKey(types.StoreKey)
	tEvmKey := storetypes.NewTransientStoreKey(fmt.Sprintf("%s_test", types.StoreKey))
	ctx := testutil.DefaultContext(evmKey, tEvmKey)
	paramstore := paramtypes.NewSubspace(
		encCfg.Codec, encCfg.Amino, evmKey, tEvmKey, "evm",
	).WithKeyTable(v3types.ParamKeyTable())

	params := v3types.DefaultParams()
	paramstore.SetParamSet(ctx, &params)

	require.Panics(t, func() {
		var preMigrationConfig types.ChainConfig
		paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &preMigrationConfig)
	})
	var preMigrationConfig v3types.ChainConfig
	paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &preMigrationConfig)
	require.NotNil(t, preMigrationConfig.MergeForkBlock)

	paramstore = paramtypes.NewSubspace(
		encCfg.Codec, encCfg.Amino, evmKey, tEvmKey, "evm",
	).WithKeyTable(types.ParamKeyTable())
	err := v3.MigrateStore(ctx, &paramstore)
	require.NoError(t, err)

	updatedDefaultConfig := types.DefaultChainConfig()

	var postMigrationConfig types.ChainConfig
	paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &postMigrationConfig)
	require.Equal(t, postMigrationConfig.GrayGlacierBlock, updatedDefaultConfig.GrayGlacierBlock)
	require.Equal(
		t,
		postMigrationConfig.MergeNetsplitBlock,
		updatedDefaultConfig.MergeNetsplitBlock,
	)
	require.Panics(t, func() {
		var preMigrationConfig v3types.ChainConfig
		paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &preMigrationConfig)
	})
}
