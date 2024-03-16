package types

import (
	"math/big"
	"sync"
)

var (
	chainIDonce   sync.Once
	chainIDParser ChainIDParser
)

// ChainIDParser alias the function ParseChainID
type ChainIDParser func(chainID string) (*big.Int, error)

// InjectChainIDParser inject the externally implemented ParseChainID function, which can only be executed once
func InjectChainIDParser(parser ChainIDParser) {
	chainIDonce.Do(func() {
		chainIDParser = parser
	})
}
