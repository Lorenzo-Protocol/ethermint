package evmv1

import (
	"fmt"

	"cosmossdk.io/x/tx/signing"
	protov2 "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// supportedTxs holds the Ethereum transaction types
var supportedTxs = map[string]TxDataV2{
	"/ethermint.evm.v1.DynamicFeeTx": &DynamicFeeTx{},
	"/ethermint.evm.v1.AccessListTx": &AccessListTx{},
	"/ethermint.evm.v1.LegacyTx":     &LegacyTx{},
}

// getSender extracts the sender address from the signature values using the latest signer for the given chainID.
func getSender(txData TxDataV2) (common.Address, error) {
	signer := ethtypes.LatestSignerForChainID(txData.GetChainID())
	from, err := signer.Sender(ethtypes.NewTx(txData.AsEthereumData()))
	if err != nil {
		return common.Address{}, err
	}
	return from, nil
}

// CustomGetSigner returns a map of GetSigners functions for the different Ethereum transaction types.
//
// The returned map contains a single entry with the protoreflect.FullName of MsgEthereumTx as the key.
// The value is a function that takes an arbitrary protov2.Message and returns the sender address as a
// []byte. The function returns an error if the message is not a MsgEthereumTx or if the sender address
// cannot be extracted from the transaction data.
//
// The function is used by the SDK's tx.Config to determine the signers for a transaction.
func CustomGetSigner() map[protoreflect.FullName]signing.GetSignersFunc {
	return map[protoreflect.FullName]signing.GetSignersFunc{
		protov2.MessageName(&MsgEthereumTx{}): func(msg protov2.Message) ([][]byte, error) {
			msgEthTx, ok := msg.(*MsgEthereumTx)
			if !ok {
				return nil, fmt.Errorf("invalid type, expected MsgEthereumTx and got %T", msg)
			}

			txData, found := supportedTxs[msgEthTx.Data.TypeUrl]
			if !found {
				return nil, fmt.Errorf("invalid TypeUrl %s", msgEthTx.Data.TypeUrl)
			}

			if err := msgEthTx.Data.UnmarshalTo(txData); err != nil {
				return nil, err
			}

			sender, err := getSender(txData)
			if err != nil {
				return nil, err
			}

			return [][]byte{sender.Bytes()}, nil
		},
	}
}
