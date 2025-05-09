// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package common

import (
	stdbytes "bytes"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
)

/* -------------------------------------------------------------------------- */
/*                                    Root                                    */
/* -------------------------------------------------------------------------- */

type (
	// Bytes32 defines the commonly used 32-byte array.
	Bytes32 = bytes.B32

	// Domain as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Domain = bytes.B32

	// DomainType as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	DomainType = bytes.B4

	// Hash32 as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Hash32 = bytes.B32

	// Version as per the Ethereum 2.0 specification.
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Version = bytes.B4

	// ForkDigest as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	ForkDigest = bytes.B4
)

// Root represents a 32-byte Merkle root.
// We use this type to represent roots that come from the consensus layer.
type Root [RootSize]byte

const RootSize = 32

// NewRootFromHex creates a new root from a hex string.
//
// Errors if:
// - input is not prefixed with "0x".
// - input is not valid hex of 32 bytes.
func NewRootFromHex(input string) (Root, error) {
	val, err := hex.ToBytes(input)
	if err != nil {
		return Root{}, err
	}
	if len(val) != RootSize {
		return Root{}, bytes.ErrIncorrectLength
	}
	return Root(val), nil
}

// NewRootFromBytes creates a new root from a byte slice.
func NewRootFromBytes(input []byte) Root {
	var root Root
	copy(root[:], input)
	return root
}

// Equals returns true if the two roots are equal.
func (r Root) Equals(other Root) bool {
	return stdbytes.Equal(r[:], other[:])
}

// Hex converts a root to a hex string.
func (r Root) Hex() string { return hex.EncodeBytes(r[:]) }

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (r Root) String() string {
	return r.Hex()
}

// MarshalText returns the hex representation of r.
func (r Root) MarshalText() ([]byte, error) {
	return []byte(r.Hex()), nil
}

// UnmarshalText parses a root in hex syntax.
func (r *Root) UnmarshalText(input []byte) error {
	var err error
	*r, err = NewRootFromHex(string(input))
	return err
}

// MarshalJSON returns the JSON representation of r.
func (r Root) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Hex())
}

// UnmarshalJSON parses a root in hex syntax.
//
// NOTE: Enforces the input to include any extra character in the first and last position.
// Technically this is used to remove the quote `"`. For example, the input may look like:
// []byte(`"0x6969696969696969696969696969696969696969696969696969696969696969"`)
func (r *Root) UnmarshalJSON(input []byte) error {
	if len(input) <= 1 {
		return errors.Wrapf(
			bytes.ErrIncorrectLength, "input length (%d) is too small", len(input),
		)
	}
	return r.UnmarshalText(input[1 : len(input)-1])
}
