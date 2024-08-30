package testutil

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/app"
)

// Commit commits a block at a given time. Reminder: At the end of each
// Tendermint Consensus round the following methods are run
//  1. BeginBlock
//  2. DeliverTx
//  3. EndBlock
//  4. Commit
func Commit(ctx sdk.Context, app *app.EthermintApp, t time.Duration, vs *tmtypes.ValidatorSet) (sdk.Context, error) {
	header := ctx.BlockHeader()
	res, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:          header.Height,
		ProposerAddress: header.ProposerAddress,
	})
	if err != nil {
		return ctx, err
	}

	if vs != nil {
		nextVals, err := applyValSetChanges(vs, res.ValidatorUpdates)
		if err != nil {
			return ctx, err
		}
		header.ValidatorsHash = vs.Hash()
		header.NextValidatorsHash = nextVals.Hash()
	}

	if _, err := app.Commit(); err != nil {
		return ctx, err
	}

	header.Height++
	header.Time = header.Time.Add(t)
	header.AppHash = app.LastCommitID().Hash

	if _, err := app.PrepareProposal(&abci.RequestPrepareProposal{
		Height: header.Height,
		Time:   header.Time,
	}); err != nil {
		return ctx, err
	}

	return ctx.WithBlockHeader(header), nil
}

// applyValSetChanges takes in tmtypes.ValidatorSet and []abci.ValidatorUpdate and will return a new tmtypes.ValidatorSet which has the
// provided validator updates applied to the provided validator set.
func applyValSetChanges(valSet *tmtypes.ValidatorSet, valUpdates []abci.ValidatorUpdate) (*tmtypes.ValidatorSet, error) {
	updates, err := tmtypes.PB2TM.ValidatorUpdates(valUpdates)
	if err != nil {
		return nil, err
	}

	// must copy since validator set will mutate with UpdateWithChangeSet
	newVals := valSet.Copy()
	err = newVals.UpdateWithChangeSet(updates)
	if err != nil {
		return nil, err
	}

	return newVals, nil
}
