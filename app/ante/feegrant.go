package ante

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type EthFeeGrantValidator struct {
	feegrantKeeper FeegrantKeeper
	evmKeeper      EVMKeeper
}

// NewEthFeeGrantValidator creates a new EthFeeGrantValidator
func NewEthFeeGrantValidator(evmKeeper EVMKeeper, fk FeegrantKeeper) EthFeeGrantValidator {
	return EthFeeGrantValidator{
		feegrantKeeper: fk,
		evmKeeper:      evmKeeper,
	}
}

func (ev EthFeeGrantValidator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	params := ev.evmKeeper.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(ev.evmKeeper.ChainID())
	blockNum := big.NewInt(ctx.BlockHeight())
	signer := ethtypes.MakeSigner(ethCfg, blockNum)
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type %T, expected %T", tx, (*evmtypes.MsgEthereumTx)(nil))
		}
		ethTx := msgEthTx.AsTransaction()
		sender, err := signer.Sender(ethTx)
		if err != nil {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrorInvalidSigner,
				"couldn't retrieve sender address ('%s') from the ethereum transaction: %s",
				msgEthTx.From,
				err.Error(),
			)
		}
		txData, err := evmtypes.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, sdkerrors.Wrap(err, "failed to unpack tx data")
		}
		feeGrantee := sender.Bytes()
		feeGranteeCosmosAddr := sdk.AccAddress(feeGrantee)
		feePayer := msgEthTx.GetFeePayer()
		feeAmt := txData.Fee()
		if feeAmt.Sign() == 0 {
			return ctx, sdkerrors.Wrap(err, "failed to fee amount")
		}

		fees := sdk.Coins{sdk.NewCoin(params.EvmDenom, sdk.NewIntFromBigInt(feeAmt))}

		msgs := []sdk.Msg{msg}

		if feePayer != nil {
			err := ev.feegrantKeeper.UseGrantedFees(ctx, feePayer, feeGrantee, fees, msgs)
			if err != nil {
				return ctx, sdkerrors.Wrapf(err,
					"%s(%s) not allowed to pay fees from %s", sender.Hex(), feeGranteeCosmosAddr, feePayer)
			}
		}
	}
	return next(ctx, tx, simulate)
}
