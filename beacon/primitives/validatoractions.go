package primitives

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/phoreproject/synapse/bls"
)

// DepositParameters are the parameters the depositer needs
// to provide.
type DepositParameters struct {
	PubKey                bls.PublicKey
	ProofOfPossession     bls.Signature
	WithdrawalCredentials chainhash.Hash
	RandaoCommitment      chainhash.Hash
}

// Copy returns a copy of the deposit parameters
func (dp *DepositParameters) Copy() DepositParameters {
	newDP := *dp
	newDP.PubKey = dp.PubKey.Copy()
	newSig := dp.ProofOfPossession.Copy()
	newDP.ProofOfPossession = *newSig
	return newDP
}

// Deposit is a new deposit from a shard.
type Deposit struct {
	Parameters DepositParameters
}

// Copy returns a copy of the deposit.
func (d Deposit) Copy() Deposit {
	return Deposit{d.Parameters.Copy()}
}

// Exit exits the validator.
type Exit struct {
	Slot           uint64
	ValidatorIndex uint64
	Signature      bls.Signature
}

// Copy returns a copy of the exit.
func (e *Exit) Copy() Exit {
	newSig := e.Signature.Copy()
	return Exit{
		e.Slot,
		e.ValidatorIndex,
		*newSig,
	}
}
