package misc

// This file is for swapping coins from Phore to Synapse

import (
	"github.com/phoreproject/synapse/chainhash"
)

type proofEntry struct {
	hash *chainhash.Hash
	left bool
}

func computeCombinedHash(left *chainhash.Hash, right *chainhash.Hash) chainhash.Hash {
	return chainhash.HashH(append(left[:], right[:]...))
}

func computeProofRootHash(hash *chainhash.Hash, entries []proofEntry) chainhash.Hash {
	h := *hash
	for i := 0; i < len(entries); i++ {
		if entries[i].left {
			h = computeCombinedHash(entries[i].hash, &h)
		} else {
			h = computeCombinedHash(&h, entries[i].hash)
		}
	}
	return h
}

// The text is generated by SynapseSwap::proofListToText in C++
func textToProofList(text string) []proofEntry {
	proofList := []proofEntry{}

	for i := 0; i <= len(text)-65; i += 65 {
		direction := text[i : i+1]
		hex := text[i+1 : i+65]
		entry := proofEntry{}
		entry.left = (direction == "L")
		entry.hash, _ = chainhash.NewHashFromStr(hex)
		proofList = append(proofList, entry)
	}

	return proofList
}
