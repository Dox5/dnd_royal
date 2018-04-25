package model

import (
  "crypto/rand"
  "math"
  "math/big"
)

type Identifier uint64

func MakeId() Identifier {
    max := big.Int{}

    max.SetUint64(math.MaxUint64)
    id, err := rand.Int(rand.Reader, &max)

    // TODO: Pass error back to initiator!
    if err != nil {
        panic("Arrgh")
    }

    return Identifier(id.Uint64())
}
