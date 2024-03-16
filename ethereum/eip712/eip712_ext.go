package eip712

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type Codec struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Amino             *codec.LegacyAmino
}

// InjectCodec set the encoding config to the singleton codecs (Amino and Protobuf).
// The process of unmarshaling SignDoc bytes into a SignDoc object requires having a codec
// populated with all relevant message types. As a result, we must call this method on app
// initialization with the app's encoding config.
func InjectCodec(cdc Codec) {
	aminoCodec = cdc.Amino
	protoCodec = codec.NewProtoCodec(cdc.InterfaceRegistry)
}
