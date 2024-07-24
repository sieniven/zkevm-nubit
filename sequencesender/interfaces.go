package sequencesender

import (
	"context"

	ethmanTypes "github.com/sieniven/zkevm-nubit/etherman/types"
)

type dataAbilitier interface {
	PostSequence(ctx context.Context, sequences []ethmanTypes.Sequence) ([]byte, []byte, error)
}
