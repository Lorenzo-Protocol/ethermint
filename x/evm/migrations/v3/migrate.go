package v3

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/evmos/ethermint/x/evm/types"
)

// MigrateStore sets the default for GrayGlacierBlock and MergeNetsplitBlock in ChainConfig parameter.
func MigrateStore(ctx sdk.Context, legacySubspace types.Subspace) error {
	paramstore, ok := legacySubspace.(*paramtypes.Subspace)
	if !ok {
		return fmt.Errorf("invalid legacySubspace type: %T", paramstore)
	}

	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore = &ps
	}
	prevConfig := &types.ChainConfig{}
	paramstore.GetIfExists(ctx, types.ParamStoreKeyChainConfig, prevConfig)

	defaultConfig := types.DefaultChainConfig()

	prevConfig.GrayGlacierBlock = defaultConfig.GrayGlacierBlock
	prevConfig.MergeNetsplitBlock = defaultConfig.MergeNetsplitBlock

	paramstore.Set(ctx, types.ParamStoreKeyChainConfig, prevConfig)
	return nil
}
