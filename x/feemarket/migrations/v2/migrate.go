package v2

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	v1types "github.com/evmos/ethermint/x/feemarket/migrations/v1"
	v2types "github.com/evmos/ethermint/x/feemarket/migrations/v2/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

// KeyPrefixBaseFeeV1 is the base fee key prefix used in version 1
var KeyPrefixBaseFeeV1 = []byte{2}

// MigrateStore migrates the BaseFee value from the store to the params for
// In-Place Store migration logic.
func MigrateStore(
	ctx sdk.Context,
	legacySubspace feemarkettypes.Subspace,
	storeKey storetypes.StoreKey,
) error {
	paramstore, ok := legacySubspace.(*paramtypes.Subspace)
	if !ok {
		return fmt.Errorf("invalid legacySubspace type: %T", paramstore)
	}

	baseFee := v2types.DefaultParams().BaseFee

	store := ctx.KVStore(storeKey)

	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(v2types.ParamKeyTable())
		paramstore = &ps
	}

	switch {
	case store.Has(KeyPrefixBaseFeeV1):
		bz := store.Get(KeyPrefixBaseFeeV1)
		baseFee = math.NewIntFromBigInt(new(big.Int).SetBytes(bz))
	case paramstore.Has(ctx, v2types.ParamStoreKeyNoBaseFee):
		paramstore.GetIfExists(ctx, v2types.ParamStoreKeyBaseFee, &baseFee)
	}

	var (
		noBaseFee                                bool
		baseFeeChangeDenom, elasticityMultiplier uint32
		enableHeight                             int64
	)

	paramstore.GetIfExists(ctx, v2types.ParamStoreKeyNoBaseFee, &noBaseFee)
	paramstore.GetIfExists(ctx, v2types.ParamStoreKeyBaseFeeChangeDenominator, &baseFeeChangeDenom)
	paramstore.GetIfExists(ctx, v2types.ParamStoreKeyElasticityMultiplier, &elasticityMultiplier)
	paramstore.GetIfExists(ctx, v2types.ParamStoreKeyEnableHeight, &enableHeight)

	params := v2types.Params{
		NoBaseFee:                noBaseFee,
		BaseFeeChangeDenominator: baseFeeChangeDenom,
		ElasticityMultiplier:     elasticityMultiplier,
		BaseFee:                  baseFee,
		EnableHeight:             enableHeight,
	}

	paramstore.SetParamSet(ctx, &params)
	store.Delete(KeyPrefixBaseFeeV1)
	return nil
}

// MigrateJSON accepts exported v0.9 x/feemarket genesis state and migrates it to
// v0.10 x/feemarket genesis state. The migration includes:
// - Migrate BaseFee to Params
func MigrateJSON(oldState v1types.GenesisState) v2types.GenesisState {
	return v2types.GenesisState{
		Params: v2types.Params{
			NoBaseFee:                oldState.Params.NoBaseFee,
			BaseFeeChangeDenominator: oldState.Params.BaseFeeChangeDenominator,
			ElasticityMultiplier:     oldState.Params.ElasticityMultiplier,
			EnableHeight:             oldState.Params.EnableHeight,
			BaseFee:                  oldState.BaseFee,
		},
		BlockGas: oldState.BlockGas,
	}
}
