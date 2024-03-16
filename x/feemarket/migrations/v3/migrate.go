package v3

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	v2types "github.com/evmos/ethermint/x/feemarket/migrations/v2/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

// MigrateStore adds the MinGasPrice param with a value of 0
// and MinGasMultiplier to 0,5
func MigrateStore(ctx sdk.Context, legacySubspace feemarkettypes.Subspace) error {
	paramstore, ok := legacySubspace.(*paramtypes.Subspace)
	if !ok {
		return fmt.Errorf("invalid legacySubspace type: %T", paramstore)
	}

	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(feemarkettypes.ParamKeyTable())
		paramstore = &ps
	}

	// add MinGasPrice
	paramstore.Set(ctx, feemarkettypes.ParamStoreKeyMinGasPrice, feemarkettypes.DefaultMinGasPrice)
	// add MinGasMultiplier
	paramstore.Set(
		ctx,
		feemarkettypes.ParamStoreKeyMinGasMultiplier,
		feemarkettypes.DefaultMinGasMultiplier,
	)
	return nil
}

// MigrateJSON accepts exported v0.10 x/feemarket genesis state and migrates it to
// v0.11 x/feemarket genesis state. The migration includes:
// - add MinGasPrice param
// - add MinGasMultiplier param
func MigrateJSON(oldState v2types.GenesisState) feemarkettypes.GenesisState {
	return feemarkettypes.GenesisState{
		Params: feemarkettypes.Params{
			NoBaseFee:                oldState.Params.NoBaseFee,
			BaseFeeChangeDenominator: oldState.Params.BaseFeeChangeDenominator,
			ElasticityMultiplier:     oldState.Params.ElasticityMultiplier,
			EnableHeight:             oldState.Params.EnableHeight,
			BaseFee:                  oldState.Params.BaseFee,
			MinGasPrice:              feemarkettypes.DefaultMinGasPrice,
			MinGasMultiplier:         feemarkettypes.DefaultMinGasMultiplier,
		},
		BlockGas: oldState.BlockGas,
	}
}
