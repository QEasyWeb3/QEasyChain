package systemcontract

import (
	"fmt"
	"github.com/QEasyWeb3/QEasyChain/contracts/system"
	"math"
	"math/big"

	"github.com/QEasyWeb3/QEasyChain/accounts/abi"

	"github.com/QEasyWeb3/QEasyChain/common"
	"github.com/QEasyWeb3/QEasyChain/core"
	"github.com/QEasyWeb3/QEasyChain/core/state"
	"github.com/QEasyWeb3/QEasyChain/core/types"
	"github.com/QEasyWeb3/QEasyChain/core/vm"
	"github.com/QEasyWeb3/QEasyChain/log"
	"github.com/QEasyWeb3/QEasyChain/params"
)

type CallContext struct {
	Statedb      *state.StateDB
	Header       *types.Header
	ChainContext core.ChainContext
	ChainConfig  *params.ChainConfig
}

func (ctx *CallContext) GetContractVersion(contractName string) uint8 {
	return system.GetContractVersion(contractName, ctx.Header.Number, ctx.ChainConfig)
}

func (ctx *CallContext) GetContractAddress(contractName string) common.Address {
	return system.GetContractAddress(contractName, ctx.GetContractVersion(contractName))
}

// CallContract executes transaction sent to system contracts.
func CallContract(ctx *CallContext, to common.Address, data []byte) (ret []byte, err error) {
	return CallContractWithValue(ctx, system.LocalAddress, to, data, big.NewInt(0))
}

// CallContract executes transaction sent to system contracts.
func CallContractWithValue(ctx *CallContext, from common.Address, to common.Address, data []byte, value *big.Int) (ret []byte, err error) {
	evm := vm.NewEVM(core.NewEVMBlockContext(ctx.Header, ctx.ChainContext, nil), vm.TxContext{
		Origin:   from,
		GasPrice: big.NewInt(0),
	}, ctx.Statedb, ctx.ChainConfig, vm.Config{})

	ret, _, err = evm.Call(vm.AccountRef(from), to, data, math.MaxUint64, value)
	// Finalise the statedb so any changes can take effect,
	// and especially if the `from` account is empty, it can be finally deleted.
	ctx.Statedb.Finalise(true)

	return ret, WrapVMError(err, ret)
}

// VMCallContract executes transaction sent to system contracts with given EVM.
func VMCallContract(evm *vm.EVM, from common.Address, to common.Address, data []byte, gas uint64) (ret []byte, err error) {
	state, ok := evm.StateDB.(*state.StateDB)
	if !ok {
		log.Crit("Unknown statedb type")
	}
	ret, _, err = evm.Call(vm.AccountRef(from), to, data, gas, big.NewInt(0))
	// Finalise the statedb so any changes can take effect,
	// and especially if the `from` account is empty, it can be finally deleted.
	state.Finalise(true)

	return ret, WrapVMError(err, ret)
}

// WrapVMError wraps vm error with readable reason
func WrapVMError(err error, ret []byte) error {
	if err == vm.ErrExecutionReverted {
		reason, errUnpack := abi.UnpackRevert(common.CopyBytes(ret))
		if errUnpack != nil {
			reason = "internal error"
		}
		return fmt.Errorf("%s: %s", err.Error(), reason)
	}
	return err
}