// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package core

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/QEasyWeb3/QEasyChain/common"
	"github.com/QEasyWeb3/QEasyChain/common/math"
)

var _ = (*validatorInfoMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (v ValidatorInfo) MarshalJSON() ([]byte, error) {
	type ValidatorInfo struct {
		Signer           common.Address        `json:"signer"         gencodec:"required"`
		Owner            common.Address        `json:"owner"         gencodec:"required"`
		Rate             *math.HexOrDecimal256 `json:"rate,omitempty"`
		Stake            *math.HexOrDecimal256 `json:"stake,omitempty"`
		AcceptDelegation bool                  `json:"acceptDelegation,omitempty"`
	}
	var enc ValidatorInfo
	enc.Signer = v.Signer
	enc.Owner = v.Owner
	enc.Rate = (*math.HexOrDecimal256)(v.Rate)
	enc.Stake = (*math.HexOrDecimal256)(v.Stake)
	enc.AcceptDelegation = v.AcceptDelegation
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (v *ValidatorInfo) UnmarshalJSON(input []byte) error {
	type ValidatorInfo struct {
		Signer           *common.Address       `json:"signer"         gencodec:"required"`
		Owner            *common.Address       `json:"owner"         gencodec:"required"`
		Rate             *math.HexOrDecimal256 `json:"rate,omitempty"`
		Stake            *math.HexOrDecimal256 `json:"stake,omitempty"`
		AcceptDelegation *bool                 `json:"acceptDelegation,omitempty"`
	}
	var dec ValidatorInfo
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Signer == nil {
		return errors.New("missing required field 'signer' for ValidatorInfo")
	}
	v.Signer = *dec.Signer
	if dec.Owner == nil {
		return errors.New("missing required field 'owner' for ValidatorInfo")
	}
	v.Owner = *dec.Owner
	if dec.Rate != nil {
		v.Rate = (*big.Int)(dec.Rate)
	}
	if dec.Stake != nil {
		v.Stake = (*big.Int)(dec.Stake)
	}
	if dec.AcceptDelegation != nil {
		v.AcceptDelegation = *dec.AcceptDelegation
	}
	return nil
}
