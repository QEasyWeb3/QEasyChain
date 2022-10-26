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
	// ValidatorsInteractiveABI contains all methods to interactive with validator contracts.
	ValidatorsInteractiveABI = `
    [
    {
        "inputs": [
        {
            "internalType": "address[]",
            "name": "vals",
            "type": "address[]"
        }
        ],
        "name": "initialize",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "distributeBlockReward",
        "outputs": [],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [],
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
            "internalType": "address[]",
            "name": "newSet",
            "type": "address[]"
        },
        {
            "internalType": "uint256",
            "name": "epoch",
            "type": "uint256"
        }
        ],
        "name": "updateActiveValidatorSet",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
            "internalType": "address",
            "name": "val",
            "type": "address"
            }
        ],
        "name": "getValidatorInfo",
        "outputs": [
            {
            "internalType": "address payable",
            "name": "",
            "type": "address"
            },
            {
            "internalType": "enum Validators.Status",
            "name": "",
            "type": "uint8"
            },
            {
            "internalType": "uint256",
            "name": "",
            "type": "uint256"
            },
            {
            "internalType": "uint256",
            "name": "",
            "type": "uint256"
            },
            {
            "internalType": "uint256",
            "name": "",
            "type": "uint256"
            },
            {
            "internalType": "uint256",
            "name": "",
            "type": "uint256"
            },
            {
            "internalType": "address[]",
            "name": "",
            "type": "address[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    }
]
    `

	PunishInteractiveABI = `
    [
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
                "name": "val",
                "type": "address"
            }
            ],
            "name": "punish",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [
            {
                "internalType": "uint256",
                "name": "epoch",
                "type": "uint256"
            }
            ],
            "name": "decreaseMissedBlocksCounter",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]
    `

	ProposalInteractiveABI = `
    [
        {
            "inputs": [
            {
                "internalType": "address[]",
                "name": "vals",
                "type": "address[]"
            }
            ],
            "name": "initialize",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]
    `

	SysGovInteractiveABI = `
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
        "name": "initializeV2",
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

	ValidatorsV1InteractiveABI = `[
        {
            "inputs": [
                {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                }
            ],
            "name": "activeValidators",
            "outputs": [
                {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "distributeBlockReward",
            "outputs": [],
            "stateMutability": "payable",
            "type": "function"
        },
        {
            "inputs": [],
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
                    "internalType": "address[]",
                    "name": "_candidates",
                    "type": "address[]"
                },
                {
                    "internalType": "address[]",
                    "name": "_manager",
                    "type": "address[]"
                },
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
                    "internalType": "address[]",
                    "name": "newSet",
                    "type": "address[]"
                },
                {
                    "internalType": "uint256",
                    "name": "epoch",
                    "type": "uint256"
                }
            ],
            "name": "updateActiveValidatorSet",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]`

	PunishV1InteractiveABI = `[
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
            "name": "_val",
            "type": "address"
          }
        ],
        "name": "punish",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
          {
            "internalType": "uint256",
            "name": "_epoch",
            "type": "uint256"
          }
        ],
        "name": "decreaseMissedBlocksCounter",
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
          },
          {
            "internalType": "address",
            "name": "val",
            "type": "address"
          }
        ],
        "name": "doubleSignPunish",
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
	SystemContractName        = "SystemContract"
	ValidatorsContractName    = "ValidatorsContract"
	PunishContractName        = "PunishContract"
	ProposalContractName      = "ProposalContract"
	SysGovContractName        = "SysGovContract"
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
	RedCoastVersion
	SophonVersion
	WaterdropVersion
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
	ValidatorsContract    = common.HexToAddress("0x000000000000000000000000000000000000F000")
	CommunityPoolContract = common.HexToAddress("0x000000000000000000000000000000000000F001")
	PunishContract        = common.HexToAddress("0x000000000000000000000000000000000000f001")
	ProposalContract      = common.HexToAddress("0x000000000000000000000000000000000000f002")
	SysGovContract        = common.HexToAddress("0x000000000000000000000000000000000000F003")
	AddressListContract   = common.HexToAddress("0x000000000000000000000000000000000000F004")
	ValidatorsV1Contract  = common.HexToAddress("0x000000000000000000000000000000000000F005")
	PunishV1Contract      = common.HexToAddress("0x000000000000000000000000000000000000F006")

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
		ValidatorsContractName: {
			ContractV0: {
				abi:  ValidatorsInteractiveABI,
				addr: ValidatorsContract,
			},
			ContractV1: {
				abi:  ValidatorsV1InteractiveABI,
				addr: ValidatorsV1Contract,
			},
		},
		PunishContractName: {
			ContractV0: {
				abi:  PunishInteractiveABI,
				addr: PunishContract,
			},
			ContractV1: {
				abi:  PunishInteractiveABI,
				addr: PunishV1Contract,
			},
			ContractV2: {
				abi:  PunishV1InteractiveABI,
				addr: PunishV1Contract,
			},
		},
		ProposalContractName: {
			ContractV0: {
				abi:  ProposalInteractiveABI,
				addr: ProposalContract,
			},
		},
		SysGovContractName: {
			ContractV0: {
				abi:  SysGovInteractiveABI,
				addr: SysGovContract,
			},
		},
		AddressListContractName: {
			ContractV0: {
				abi:  AddrListInteractiveABI,
				addr: AddressListContract,
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
	if config.IsWaterdrop(blockNum) {
		return WaterdropVersion
	}
	if config.IsSophon(blockNum) {
		return SophonVersion
	}
	if config.IsRedCoast(blockNum) {
		return RedCoastVersion
	}
	return GenesisVersion
}

func GetContractVersion(contractName string, blockNum *big.Int, config *params.ChainConfig) uint8 {
	sysContractVersion := GetSysContractVersion(blockNum, config)
	switch contractName {
	case ValidatorsContractName:
		{
			if sysContractVersion >= RedCoastVersion {
				return ContractV1
			}
			return ContractV0
		}
	case PunishContractName:
		{
			if sysContractVersion >= WaterdropVersion {
				return ContractV2
			} else if sysContractVersion >= RedCoastVersion {
				return ContractV1
			}
			return ContractV0
		}
	case ProposalContractName:
		{
			return ContractV0
		}
	case SysGovContractName:
		{
			if sysContractVersion >= RedCoastVersion {
				return ContractV0
			}
			break
		}
	case AddressListContractName:
		{
			if sysContractVersion >= RedCoastVersion {
				return ContractV0
			}
			break
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
