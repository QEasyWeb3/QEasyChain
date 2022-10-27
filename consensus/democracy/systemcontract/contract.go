package systemcontract

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/contracts/system"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

// AddrAscend implements the sort interface to allow sorting a list of addresses
type AddrAscend []common.Address

func (s AddrAscend) Len() int           { return len(s) }
func (s AddrAscend) Less(i, j int) bool { return bytes.Compare(s[i][:], s[j][:]) < 0 }
func (s AddrAscend) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Proposal struct {
	Id     *big.Int
	Action *big.Int
	From   common.Address
	To     common.Address
	Value  *big.Int
	Data   []byte
}

// GetTopValidators return the result of calling method `getTopValidators` in Staking contract
func GetTopValidators(ctx *CallContext) ([]common.Address, error) {
	const method = "getTopValidators"
	result, err := contractRead(ctx, system.SysContractName, method, system.MaxValidators)
	if err != nil {
		log.Error("GetTopValidators contractRead failed", "err", err)
		return []common.Address{}, err
	}
	validators, ok := result.([]common.Address)
	if !ok {
		return []common.Address{}, errors.New("GetTopValidators: invalid validator format")
	}
	sort.Sort(AddrAscend(validators))
	return validators, nil
}

// UpdateActiveValidatorSet return the result of calling method `updateActiveValidatorSet` in Staking contract
func UpdateActiveValidatorSet(ctx *CallContext, newValidators []common.Address) error {
	const method = "updateActiveValidatorSet"
	err := contractWrite(ctx, system.SysContractName, method, newValidators)
	if err != nil {
		log.Error("UpdateActiveValidatorSet failed", "newValidators", newValidators, "err", err)
	}
	return err
}

// DecreaseMissedBlocksCounter return the result of calling method `decreaseMissedBlocksCounter` in Staking contract
func DecreaseMissedBlocksCounter(ctx *CallContext) error {
	const method = "decreaseMissedBlocksCounter"
	err := contractWrite(ctx, system.SysContractName, method)
	if err != nil {
		log.Error("DecreaseMissedBlocksCounter failed", "err", err)
	}
	return err
}

// DistributeBlockFee return the result of calling method `distributeBlockFee` in Staking contract
func DistributeBlockFee(ctx *CallContext, fee *big.Int) error {
	contractName := system.SysContractName
	const method = "distributeBlockFee"
	data, err := system.ABIPack(contractName, ctx.GetContractVersion(contractName), method)
	if err != nil {
		log.Error("Can't pack data for distributeBlockFee", "error", err)
		return err
	}
	if _, err := CallContractWithValue(ctx, ctx.Header.Coinbase, ctx.GetContractAddress(contractName), data, fee); err != nil {
		log.Error("distributeBlockFee failed", "fee", fee, "err", err)
		return err
	}
	return nil
}

// LazyPunish return the result of calling method `lazyPunish` in Staking contract
func LazyPunish(ctx *CallContext, validator common.Address) error {
	const method = "lazyPunish"
	err := contractWrite(ctx, system.SysContractName, method, validator)
	if err != nil {
		log.Error("LazyPunish failed", "validator", validator, "err", err)
	}
	return err
}

// DoubleSignPunish return the result of calling method `doubleSignPunish` in Staking contract
func DoubleSignPunish(ctx *CallContext, punishHash common.Hash, validator common.Address) error {
	const method = "doubleSignPunish"
	err := contractWrite(ctx, system.SysContractName, method, punishHash, validator)
	if err != nil {
		log.Error("DoubleSignPunish failed", "punishHash", punishHash, "validator", validator, "err", err)
	}
	return err
}

// DoubleSignPunishWithGivenEVM return the result of calling method `doubleSignPunish` in Staking contract with given EVM
func DoubleSignPunishWithGivenEVM(evm *vm.EVM, from common.Address, punishHash common.Hash, validator common.Address) error {
	contractName := system.SysContractName
	version := system.GetContractVersion(contractName, evm.Context.BlockNumber, evm.ChainConfig())
	contract := system.GetContractAddress(contractName, version)

	const method = "doubleSignPunish"
	data, err := system.ABIPack(contractName, version, method, punishHash, validator)
	if err != nil {
		log.Error("Can't pack data for doubleSignPunish", "error", err)
		return err
	}
	if _, err := VMCallContract(evm, from, contract, data, math.MaxUint64); err != nil {
		log.Error("DoubleSignPunishWithGivenEVM failed", "punishHash", punishHash, "validator", validator, "err", err)
		return err
	}
	return nil
}

// IsDoubleSignPunished return the result of calling method `isDoubleSignPunished` in Staking contract
func IsDoubleSignPunished(ctx *CallContext, punishHash common.Hash) (bool, error) {
	const method = "isDoubleSignPunished"
	result, err := contractRead(ctx, system.SysContractName, method, punishHash)
	if err != nil {
		log.Error("IsDoubleSignPunished contractRead failed", "punishHash", punishHash, "err", err)
		return true, err
	}
	punished, ok := result.(bool)
	if !ok {
		return true, errors.New("IsDoubleSignPunished: invalid result format, punishHash" + punishHash.Hex())
	}
	return punished, nil
}

// GetBlacksFrom return the access tx-from list
func GetBlacksFrom(ctx *CallContext) ([]common.Address, error) {
	const method = "getBlacksFrom"
	result, err := contractRead(ctx, system.AddressListContractName, method)
	if err != nil {
		log.Error("GetBlacksFrom contractRead failed", "err", err)
		return []common.Address{}, err
	}
	from, ok := result.([]common.Address)
	if !ok {
		return []common.Address{}, errors.New("GetBlacksFrom: invalid result format")
	}
	return from, nil
}

// GetBlacksTo return access tx-to list
func GetBlacksTo(ctx *CallContext) ([]common.Address, error) {
	const method = "getBlacksTo"
	result, err := contractRead(ctx, system.AddressListContractName, method)
	if err != nil {
		log.Error("GetBlacksTo contractRead failed", "err", err)
		return []common.Address{}, err
	}
	to, ok := result.([]common.Address)
	if !ok {
		return []common.Address{}, errors.New("GetBlacksTo: invalid result format")
	}
	return to, nil
}

// GetRuleByIndex return event log rules
func GetRuleByIndex(ctx *CallContext, idx uint32) (common.Hash, int, common.AddressCheckType, error) {
	const method = "getRuleByIndex"
	results, err := contractReadAll(ctx, system.AddressListContractName, method, idx)
	if err != nil {
		log.Error("GetRuleByIndex contractRead failed", "err", err)
		return common.Hash{}, 0, common.CheckNone, err
	}
	if len(results) != 3 {
		return common.Hash{}, 0, common.CheckNone, errors.New("GetRuleByIndex: invalid results' length")
	}
	var (
		ok    bool
		sig   common.Hash
		index *big.Int
		ctype uint8
	)
	if sig, ok = results[0].([32]byte); !ok {
		return common.Hash{}, 0, common.CheckNone, errors.New("GetRuleByIndex: invalid result sig format")
	}
	if index, ok = results[1].(*big.Int); !ok {
		return common.Hash{}, 0, common.CheckNone, errors.New("GetRuleByIndex: invalid result index format")
	}
	if ctype, ok = results[2].(uint8); !ok {
		return common.Hash{}, 0, common.CheckNone, errors.New("GetRuleByIndex: invalid result checktype format")
	}
	return sig, int(index.Int64()), common.AddressCheckType(ctype), nil
}

// GetRulesLen return event log rules length
func GetRulesLen(ctx *CallContext) (uint32, error) {
	const method = "rulesLen"
	result, err := contractRead(ctx, system.AddressListContractName, method)
	if err != nil {
		log.Error("GetRulesLen contractRead failed", "err", err)
		return 0, err
	}
	n, ok := result.(uint32)
	if !ok {
		return 0, errors.New("GetRulesLen: invalid result format")
	}
	return n, nil
}

// IsDeveloperVerificationEnabled Since the state variables are as follow:
//    bool public initialized;
//    bool public enabled;
//    address public admin;
//    address public pendingAdmin;
//    mapping(address => bool) private devs;
//
// according to [Layout of State Variables in Storage](https://docs.soliditylang.org/en/v0.8.4/internals/layout_in_storage.html),
// and after optimizer enabled, the `initialized`, `enabled` and `admin` will be packed, and stores at slot 0,
// `pendingAdmin` stores at slot 1, and the position for `devs` is 2.
func IsDeveloperVerificationEnabled(state consensus.StateReader, height *big.Int, config *params.ChainConfig) bool {
	contractName := system.AddressListContractName
	contract := system.GetContractAddressByConfig(contractName, height, config)
	compactValue := state.GetState(contract, common.Hash{})
	// Layout of slot 0:
	// [0   -    9][10-29][  30   ][    31     ]
	// [zero bytes][admin][enabled][initialized]
	enabledByte := compactValue.Bytes()[common.HashLength-2]
	return enabledByte == 0x01
}

// LastBlackUpdatedNumber returns LastBlackUpdatedNumber of address list
func LastBlackUpdatedNumber(state consensus.StateReader, height *big.Int, config *params.ChainConfig) uint64 {
	contractName := system.AddressListContractName
	contract := system.GetContractAddressByConfig(contractName, height, config)
	value := state.GetState(contract, system.BlackLastUpdatedNumberPosition)
	return value.Big().Uint64()
}

// LastRulesUpdatedNumber returns LastRulesUpdatedNumber of address list
func LastRulesUpdatedNumber(state consensus.StateReader, height *big.Int, config *params.ChainConfig) uint64 {
	contractName := system.AddressListContractName
	contract := system.GetContractAddressByConfig(contractName, height, config)
	value := state.GetState(contract, system.RulesLastUpdatedNumberPosition)
	return value.Big().Uint64()
}

// GetPassedProposalCount returns passed proposal count
func GetPassedProposalCount(ctx *CallContext) (uint32, error) {
	const method = "getPassedProposalCount"
	result, err := contractRead(ctx, system.OnChainDaoContractName, method)
	if err != nil {
		log.Error("GetPassedProposalCount contractRead failed", "err", err)
		return 0, err
	}
	count, ok := result.(uint32)
	if !ok {
		return 0, errors.New("GetPassedProposalCount: invalid result format")
	}
	return count, nil
}

// GetPassedProposalByIndex returns passed proposal by index
func GetPassedProposalByIndex(ctx *CallContext, idx uint32) (*Proposal, error) {
	contractName := system.OnChainDaoContractName
	const method = "getPassedProposalByIndex"
	abi := system.ABI(contractName, ctx.GetContractVersion(contractName))
	result, err := contractReadBytes(ctx, ctx.GetContractAddress(contractName), &abi, method, idx)
	if err != nil {
		log.Error("GetPassedProposalByIndex contractReadBytes failed", "idx", idx, "err", err)
		return nil, err
	}
	// unpack data
	prop := &Proposal{}
	if err = abi.UnpackIntoInterface(prop, method, result); err != nil {
		log.Error("GetPassedProposalByIndex UnpackIntoInterface failed", "idx", idx, "err", err)
		return nil, err
	}
	return prop, nil
}

// FinishProposalById finish passed proposal by id
func FinishProposalById(ctx *CallContext, id *big.Int) error {
	const method = "finishProposalById"
	err := contractWrite(ctx, system.OnChainDaoContractName, method, id)
	if err != nil {
		log.Error("FinishProposalById failed", "id", id, "err", err)
	}
	return err
}

// ExecuteProposal executes proposal
func ExecuteProposal(ctx *CallContext, prop *Proposal) error {
	_, err := CallContractWithValue(ctx, prop.From, prop.To, prop.Data, prop.Value)
	if err != nil {
		log.Error("ExecuteProposal failed", "proposal", prop, "err", err)
	}
	return err
}

// ExecuteProposalWithGivenEVM executes proposal by given evm
func ExecuteProposalWithGivenEVM(evm *vm.EVM, prop *Proposal, gas uint64) (ret []byte, err error) {
	if ret, err = VMCallContract(evm, prop.From, prop.To, prop.Data, gas); err != nil {
		log.Error("ExecuteProposalWithGivenEVM failed", "proposal", prop, "err", err)
	}
	return
}

// contractRead perform contract read
func contractRead(ctx *CallContext, contractName string, method string, args ...interface{}) (interface{}, error) {
	ret, err := contractReadAll(ctx, contractName, method, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) != 1 {
		return nil, errors.New(method + ": invalid result length")
	}
	return ret[0], nil
}

// contractReadAll perform contract Read and return all results
func contractReadAll(ctx *CallContext, contractName string, method string, args ...interface{}) ([]interface{}, error) {
	abi := system.ABI(contractName, ctx.GetContractVersion(contractName))
	result, err := contractReadBytes(ctx, ctx.GetContractAddress(contractName), &abi, method, args...)
	if err != nil {
		return nil, err
	}
	// unpack data
	ret, err := abi.Unpack(method, result)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// contractReadBytes perform read contract and returns bytes
func contractReadBytes(ctx *CallContext, contract common.Address, abi *abi.ABI, method string, args ...interface{}) ([]byte, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		log.Error("Can't pack data", "method", method, "error", err)
		return nil, err
	}
	result, err := CallContract(ctx, contract, data)
	if err != nil {
		log.Error("Failed to execute", "method", method, "err", err)
		return nil, err
	}
	return result, nil
}

// contractWrite perform write contract
func contractWrite(ctx *CallContext, contractName string, method string, args ...interface{}) error {
	data, err := system.ABIPack(contractName, ctx.GetContractVersion(contractName), method, args...)
	if err != nil {
		log.Error("Can't pack data", "method", method, "error", err)
		return err
	}
	if _, err := CallContract(ctx, ctx.GetContractAddress(contractName), data); err != nil {
		log.Error("Failed to execute", "method", method, "err", err)
		return err
	}
	return nil
}
