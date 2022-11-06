// Copyright 2021 The Cube Authors
// This file is part of the Cube library.
//
// The Cube library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Cube library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Cube library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/QEasyWeb3/QEasyChain/accounts/abi"
	"github.com/QEasyWeb3/QEasyChain/common"
	"github.com/QEasyWeb3/QEasyChain/contracts/system"
	"github.com/QEasyWeb3/QEasyChain/core/state"
	"github.com/QEasyWeb3/QEasyChain/core/types"
	"github.com/QEasyWeb3/QEasyChain/core/vm"
	"github.com/QEasyWeb3/QEasyChain/crypto"
	"github.com/QEasyWeb3/QEasyChain/log"
)

const (
	extraVanity = 32                     // Fixed number of extra-data prefix bytes reserved for validator vanity
	extraSeal   = crypto.SignatureLength // Fixed number of extra-data suffix bytes reserved for validator seal
)

// genesisInit is tools to init system contracts in genesis
type genesisInit struct {
	state   *state.StateDB
	header  *types.Header
	genesis *Genesis
}

// callContract executes contract in EVM
func (env *genesisInit) callContract(contractName string, method string, args ...interface{}) ([]byte, error) {
	// Pack method and args for data seg
	data, err := system.ABIPack(contractName, system.ContractV0, method, args...)
	if err != nil {
		return nil, err
	}
	// Create EVM calling message
	contract := system.GetContractAddress(contractName, system.ContractV0)
	msg := types.NewMessage(env.genesis.Coinbase, &contract, 0, big.NewInt(0), math.MaxUint64, big.NewInt(0), big.NewInt(0), big.NewInt(0), data, nil, false)
	// Set up the initial access list.
	if rules := env.genesis.Config.Rules(env.header.Number); rules.IsBerlin {
		env.state.PrepareAccessList(msg.From(), msg.To(), vm.ActivePrecompiles(rules), msg.AccessList())
	}
	// Create EVM
	evm := vm.NewEVM(NewEVMBlockContext(env.header, nil, &env.header.Coinbase), NewEVMTxContext(msg), env.state, env.genesis.Config, vm.Config{})
	// Run evm call
	ret, _, err := evm.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())

	if err == vm.ErrExecutionReverted {
		reason, errUnpack := abi.UnpackRevert(common.CopyBytes(ret))
		if errUnpack != nil {
			reason = "internal error"
		}
		err = fmt.Errorf("%s: %s", err.Error(), reason)
	}

	if err != nil {
		log.Error("ExecuteMsg failed", "err", err, "ret", string(ret))
	}
	env.state.Finalise(true)
	return ret, err
}

// initStaking initializes Staking Contract
func (env *genesisInit) initStaking() error {
	contract, ok := env.genesis.Alloc[system.SystemContract]
	if !ok {
		return errors.New("Staking Contract is missing in genesis!")
	}

	if len(env.genesis.Validators) <= 0 {
		return errors.New("validators are missing in genesis!")
	}

	totalValidatorStake := big.NewInt(0)
	for _, validator := range env.genesis.Validators {
		totalValidatorStake = new(big.Int).Add(totalValidatorStake, validator.Stake)
	}

	contract.Balance = totalValidatorStake
	env.state.SetBalance(system.SystemContract, contract.Balance)

	_, err := env.callContract(system.SysContractName, "initialize",
		contract.Init.Admin,
		big.NewInt(int64(env.genesis.Config.Democracy.Epoch)),
		new(big.Int).Mul(system.MinSelfStake, big.NewInt(1000000000000000000)),
		system.CommunityPoolContract,
		system.ShareOutBonusPercent)
	return err
}

// initCommunityPool initializes CommunityPool Contract
func (env *genesisInit) initCommunityPool() error {
	contract, ok := env.genesis.Alloc[system.CommunityPoolContract]
	if !ok {
		return errors.New("CommunityPool Contract is missing in genesis!")
	}
	_, err := env.callContract(system.CommunityPoolContractName, "initialize", contract.Init.Admin)
	return err
}

// initAddressList initializes AddressList Contract
func (env *genesisInit) initAddressList() error {
	contract, ok := env.genesis.Alloc[system.AddressListContract]
	if !ok {
		return errors.New("CommunityPool Contract is missing in genesis!")
	}
	_, err := env.callContract(system.AddressListContractName, "initialize", contract.Init.Admin)
	return err
}

// initValidators add validators into Staking contracts
// and set validator addresses to header extra data
// and return new header extra data
func (env *genesisInit) initValidators() ([]byte, error) {
	if len(env.genesis.Validators) <= 0 {
		return env.header.Extra, errors.New("validators are missing in genesis!")
	}
	activeSet := make([]common.Address, 0, len(env.genesis.Validators))
	extra := make([]byte, 0, extraVanity+common.AddressLength*len(env.genesis.Validators)+extraSeal)
	extra = append(extra, env.header.Extra[:extraVanity]...)
	for _, v := range env.genesis.Validators {
		if _, err := env.callContract(system.SysContractName, "initValidator",
			v.Address, v.Manager, v.Rate, v.Stake, v.AcceptDelegation); err != nil {
			return env.header.Extra, err
		}
		extra = append(extra, v.Address[:]...)
		activeSet = append(activeSet, v.Address)
	}
	extra = append(extra, env.header.Extra[len(env.header.Extra)-extraSeal:]...)
	env.header.Extra = extra
	if _, err := env.callContract(system.SysContractName, "updateActiveValidatorSet", activeSet); err != nil {
		return extra, err
	}
	return env.header.Extra, nil
}
