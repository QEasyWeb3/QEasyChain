package system

import (
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// SystemInteractiveABI contains all methods to interactive with system contracts.
	SystemInteractiveABI = `
    [
    {
      "inputs": [],
      "name": "decreaseMissedBlocksCounter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "distributeBlockFee",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_punishHash",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "doubleSignPunish",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getActiveValidators",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint8",
          "name": "_count",
          "type": "uint8"
        }
      ],
      "name": "getTopValidators",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "val",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "manager",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "rate",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "stake",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "acceptDelegation",
          "type": "bool"
        }
      ],
      "name": "initValidator",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "adminAddress",
          "type": "address"
        },
        {
          "internalType": "uint8",
          "name": "maxValidators",
          "type": "uint8"
        },
        {
          "internalType": "uint256",
          "name": "epoch",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "minSelfStake",
          "type": "uint256"
        },
        {
          "internalType": "address payable",
          "name": "communityAddress",
          "type": "address"
        },
        {
          "internalType": "uint8",
          "name": "shareOutBonusPercent",
          "type": "uint8"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "punishHash",
          "type": "bytes32"
        }
      ],
      "name": "isDoubleSignPunished",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "lazyPunish",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address[]",
          "name": "newSet",
          "type": "address[]"
        }
      ],
      "name": "updateActiveValidatorSet",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

	OnChainDaoInteractiveABI = `
    [
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "finishProposalById",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint32",
          "name": "index",
          "type": "uint32"
        }
      ],
      "name": "getPassedProposalByIndex",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "id",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "action",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "data",
          "type": "bytes"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getPassedProposalCount",
      "outputs": [
        {
          "internalType": "uint32",
          "name": "",
          "type": "uint32"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getProposalById",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "_id",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "action",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "data",
          "type": "bytes"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getProposalsTotalCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_admin",
          "type": "address"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

	CommunityPoolInteractiveABI = `
    [
        {
		  "inputs": [
			{
			  "internalType": "address",
			  "name": "_admin",
			  "type": "address"
			}
		  ],
		  "name": "initialize",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
		}
    ]
    `

	AddrListInteractiveABI = `
    [
        {
        "inputs": [],
        "name": "blackLastUpdatedNumber",
        "outputs": [
            {
            "internalType": "uint256",
            "name": "",
            "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [],
        "name": "devVerifyEnabled",
        "outputs": [
            {
            "internalType": "bool",
            "name": "",
            "type": "bool"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [],
        "name": "getBlacksFrom",
        "outputs": [
            {
            "internalType": "address[]",
            "name": "",
            "type": "address[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [],
        "name": "getBlacksTo",
        "outputs": [
            {
            "internalType": "address[]",
            "name": "",
            "type": "address[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [
            {
            "internalType": "uint32",
            "name": "i",
            "type": "uint32"
            }
        ],
        "name": "getRuleByIndex",
        "outputs": [
            {
            "internalType": "bytes32",
            "name": "",
            "type": "bytes32"
            },
            {
            "internalType": "uint128",
            "name": "",
            "type": "uint128"
            },
            {
            "internalType": "enum AddressList.CheckType",
            "name": "",
            "type": "uint8"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [],
        "name": "initialize",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
        },
        {
        "inputs": [
            {
            "internalType": "address",
            "name": "_admin",
            "type": "address"
            }
        ],
        "name": "initialize",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
        },
        {
        "inputs": [
            {
            "internalType": "address",
            "name": "addr",
            "type": "address"
            }
        ],
        "name": "isDeveloper",
        "outputs": [
            {
            "internalType": "bool",
            "name": "",
            "type": "bool"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [],
        "name": "rulesLastUpdatedNumber",
        "outputs": [
            {
            "internalType": "uint256",
            "name": "",
            "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        },
        {
        "inputs": [],
        "name": "rulesLen",
        "outputs": [
            {
            "internalType": "uint32",
            "name": "",
            "type": "uint32"
            }
        ],
        "stateMutability": "view",
        "type": "function"
        }
    ]`
)

// DevMappingPosition is the position of the state variable `devs`.
// Since the state variables are as follows:
//    bool public initialized;
//    bool public devVerifyEnabled;
//	  bool public checkInnerCreation;
//    address public admin;
//    address public pendingAdmin;
//
//    mapping(address => bool) private devs;
//
//    //NOTE: make sure this list is not too large!
//    address[] blacksFrom;
//    address[] blacksTo;
//    mapping(address => uint256) blacksFromMap;      // address => index+1
//    mapping(address => uint256) blacksToMap;        // address => index+1
//
//    uint256 public blackLastUpdatedNumber; // last block number when the black list is updated
//    uint256 public rulesLastUpdatedNumber;  // last block number when the rules are updated
//    // event check rules
//    EventCheckRule[] rules;
//    mapping(bytes32 => mapping(uint128 => uint256)) rulesMap;   // eventSig => checkIdx => indexInArray+1
//
// according to [Layout of State Variables in Storage](https://docs.soliditylang.org/en/v0.8.4/internals/layout_in_storage.html),
// and after optimizer enabled, the `initialized`, `devVerifyEnabled`, `checkInnerCreation` and `admin` will be packed, and stores at slot 0,
// `pendingAdmin` stores at slot 1, so the position for `devs` is 2.
const DevMappingPosition = 2

const (
	SysContractName           = "SystemContract"
	OnChainDaoContractName    = "OnChainDaoContract"
	AddressListContractName   = "AddressListContract"
	CommunityPoolContractName = "CommunityPoolContract"
)

const (
	ContractV0 = iota // 0
	ContractV1        // 1
	ContractV2        // 2
	ContractV3        // 3
	ContractV4        // 4
	ContractV5        // 5
)

const (
	GenesisVersion SysContractVersion = iota
	EarthVersion
)

type SysContractVersion int

var (
	BlackLastUpdatedNumberPosition = common.BytesToHash([]byte{0x07})
	RulesLastUpdatedNumberPosition = common.BytesToHash([]byte{0x08})
)

var (
	AdminForDevelopChain common.Address
)

var (
	MaxValidators        = 21
	MinSelfStake         = big.NewInt(100)
	ShareOutBonusPercent = 10
)

var (
	SystemContract        = common.HexToAddress("0x000000000000000000000000000000000000F000")
	OnChainDaoContract    = common.HexToAddress("0x000000000000000000000000000000000000F001")
	AddressListContract   = common.HexToAddress("0x000000000000000000000000000000000000F002")
	CommunityPoolContract = common.HexToAddress("0x000000000000000000000000000000000000F003")

	addrMap map[string]map[uint8]common.Address // ContractName->version->address
	abiMap  map[string]map[uint8]abi.ABI        // ContractName->version->abi
)

type ContractInfo struct {
	abi  string
	addr common.Address
}

// init the addrMap abiMap
func init() {
	addrMap = make(map[string]map[uint8]common.Address, 0)
	abiMap = make(map[string]map[uint8]abi.ABI, 0)
	for contractName, contractInfo := range map[string]map[uint8]ContractInfo{
		SysContractName: {
			ContractV0: {
				abi:  SystemInteractiveABI,
				addr: SystemContract,
			},
		},
		OnChainDaoContractName: {
			ContractV0: {
				abi:  OnChainDaoInteractiveABI,
				addr: OnChainDaoContract,
			},
		},
		AddressListContractName: {
			ContractV0: {
				abi:  AddrListInteractiveABI,
				addr: AddressListContract,
			},
		},
		CommunityPoolContractName: {
			ContractV0: {
				abi:  CommunityPoolInteractiveABI,
				addr: CommunityPoolContract,
			},
		},
	} {

		addrSubMap := make(map[uint8]common.Address, 0)
		abiSubMap := make(map[uint8]abi.ABI, 0)
		for version, info := range contractInfo {
			if abiJson, err := abi.JSON(strings.NewReader(info.abi)); err != nil {
				panic("abi json error: " + err.Error())
			} else {
				addrSubMap[version] = info.addr
				abiSubMap[version] = abiJson
			}
		}
		addrMap[contractName] = addrSubMap
		abiMap[contractName] = abiSubMap
	}
}

// ABI return abi for given contract calling
func ABI(contractName string, version uint8) abi.ABI {
	contractABI, ok := abiMap[contractName][version]
	if !ok {
		log.Crit("Unknown system abi: ", "ContractName", contractName, "Version", version)
	}
	return contractABI
}

// ABIPack generates the data field for given contract calling
func ABIPack(contractName string, version uint8, method string, args ...interface{}) ([]byte, error) {
	return ABI(contractName, version).Pack(method, args...)
}

func GetSysContractVersion(blockNum *big.Int, config *params.ChainConfig) SysContractVersion {
	if config.IsEarth(blockNum) {
		return EarthVersion
	}
	return GenesisVersion
}

func GetContractVersion(contractName string, blockNum *big.Int, config *params.ChainConfig) uint8 {
	sysContractVersion := GetSysContractVersion(blockNum, config)
	switch contractName {
	case SysContractName:
		{
			return ContractV0
		}
	case OnChainDaoContractName:
		{
			return ContractV0
		}
	case AddressListContractName:
		{
			return ContractV0
		}
	case CommunityPoolContractName:
		{
			return ContractV0
		}
	}
	log.Crit("Unknown system contract name: "+contractName, "SysContractVersion", sysContractVersion)
	return 0
}

func GetContractAddress(contractName string, version uint8) common.Address {
	addr, ok := addrMap[contractName][version]
	if !ok {
		log.Crit("Unknown system address: ", "ContractName", contractName, "Version", version)
	}
	return addr
}

func GetContractAddressByConfig(contractName string, blockNum *big.Int, config *params.ChainConfig) common.Address {
	return GetContractAddress(contractName, GetContractVersion(contractName, blockNum, config))
}
