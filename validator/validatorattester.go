package validator

import (
	"github.com/phoreproject/prysm/shared/ssz"
	"github.com/phoreproject/synapse/bls"
	"github.com/phoreproject/synapse/chainhash"
	"github.com/phoreproject/synapse/primitives"
)

var zeroHash = [32]byte{}

func getAttestation(information attestationAssignment) (*primitives.AttestationData, [32]byte, error) {
	a := primitives.AttestationData{
		Slot:                information.slot,
		Shard:               information.shard,
		BeaconBlockHash:     information.beaconBlockHash,
		EpochBoundaryHash:   information.epochBoundaryRoot,
		ShardBlockHash:      chainhash.Hash{}, // only attest to 0 hashes in phase 0
		LatestCrosslinkHash: information.latestCrosslinks[information.shard].ShardBlockHash,
		JustifiedSlot:       information.justifiedSlot,
		JustifiedBlockHash:  information.justifiedRoot,
	}

	attestationAndPoCBit := primitives.AttestationDataAndCustodyBit{Data: a, PoCBit: false}
	hashAttestation, err := ssz.TreeHash(attestationAndPoCBit)
	if err != nil {
		return nil, [32]byte{}, err
	}

	return &a, hashAttestation, nil
}

func (v *Validator) signAttestation(hashAttestation [32]byte, data primitives.AttestationData, committeeSize uint64, committeeIndex uint64) (*primitives.Attestation, error) {
	signature, err := bls.Sign(v.keystore.GetKeyForValidator(v.id), hashAttestation[:], primitives.GetDomain(*v.forkData, data.Slot, bls.DomainAttestation))
	if err != nil {
		return nil, err
	}

	participationBitfield := make([]uint8, (committeeSize+7)/8)
	custodyBitfield := make([]uint8, (committeeSize+7)/8)
	participationBitfield[committeeIndex/8] = 1 << (committeeIndex % 8)

	att := &primitives.Attestation{
		Data:                  data,
		ParticipationBitfield: participationBitfield,
		CustodyBitfield:       custodyBitfield,
		AggregateSig:          signature.Serialize(),
	}

	return att, nil
}

func (v *Validator) attestBlock(information attestationAssignment) (*primitives.Attestation, error) {
	// create attestation
	attData, hash, err := getAttestation(information)
	if err != nil {
		return nil, err
	}

	// sign attestation
	att, err := v.signAttestation(hash, *attData, information.committeeSize, information.committeeIndex)
	if err != nil {
		return nil, err
	}

	return att, nil
}