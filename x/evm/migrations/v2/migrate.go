package v2

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/evmos/ethermint/x/evm/types"
)

// MigrateStore sets the default AllowUnprotectedTxs parameter.
func MigrateStore(ctx sdk.Context, legacySubspace types.Subspace) error {
	paramstore, ok := legacySubspace.(*paramtypes.Subspace)
	if !ok {
		return fmt.Errorf("invalid legacySubspace type: %T", paramstore)
	}
	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore = &ps
	}
	// add RejectUnprotected
	paramstore.Set(ctx, types.ParamStoreKeyAllowUnprotectedTxs, types.DefaultAllowUnprotectedTxs)
	return nil
}
