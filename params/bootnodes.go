// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import "github.com/QEasyWeb3/QEasyChain/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// Ethereum Foundation Go Bootnodes
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
var TestnetBootnodes = []string{
	"enode://32821e6ff0e176b78d95f321c0739db895d3a978ca3826c9215e7b38b3181e82b64e1083a0691eefb42947e472791abd979675674aed23c310a10cfb321edc74@1.117.105.199:32664",
	"enode://03c7f01f375dcf94632078f64a4d12f64fb1fe2558f07e7faff15fa296ae781eb42765d09147e554e7471aea39e284a65903c2a0ce83e9695528bb145391add8@1.117.105.199:32665",
}

var V5Bootnodes = []string{}

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	return ""
}
