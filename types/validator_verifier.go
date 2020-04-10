package types

import (
	"errors"

	"golang.org/x/crypto/ed25519"
)

var (
	VerifyErrUnknownAuthor     = errors.New("unknown author")
	VerifyErrInvalidSignature  = errors.New("invalid signature")
	VerifyErrTooFewSignatures  = errors.New("too few signatures")
	VerifyErrTooManySignatures = errors.New("too many signatures")
)

// ValidatorVerifier is a validator set that can verify ledger infos.
// It implements LedgerInfoVerifier.
type ValidatorVerifier struct {
	publicKeyMap map[AccountAddress]*ValidatorInfo
	epoch        uint64
	totalPower   uint64
	quorumPower  uint64
}

// FromValidatorSet builds a ValidatorVerifier from a validator set and a certain epoch number.
func (vv *ValidatorVerifier) FromValidatorSet(vs *ValidatorSet, epoch uint64) error {
	vv.publicKeyMap = make(map[AccountAddress]*ValidatorInfo)
	vv.totalPower = 0
	for _, v := range vs.Payload {
		vv.publicKeyMap[v.AccountAddress] = &ValidatorInfo{
			ConsensusPubkey:      cloneBytes(v.ConsensusPubkey),
			ConsensusVotingPower: v.ConsensusVotingPower,
		}
		vv.totalPower += v.ConsensusVotingPower
	}
	vv.quorumPower = vv.totalPower*2/3 + 1
	if vv.totalPower == 0 {
		vv.quorumPower = 0
	}
	vv.epoch = epoch
	return nil
}

// ToValidatorSet exports a list of validators and the epoch number.
func (vv *ValidatorVerifier) ToValidatorSet() (*ValidatorSet, uint64) {
	vs := make([]*ValidatorInfo, 0, len(vv.publicKeyMap))
	for addr, v := range vv.publicKeyMap {
		vpk := &ValidatorInfo{
			AccountAddress:       addr,
			ConsensusPubkey:      cloneBytes(v.ConsensusPubkey),
			ConsensusVotingPower: v.ConsensusVotingPower,
		}
		vs = append(vs, vpk)
	}
	vs1 := &ValidatorSet{
		Scheme:  SchemeED25519{},
		Payload: vs,
	}
	return vs1, vv.epoch
}

func (vv *ValidatorVerifier) verifySingle(author AccountAddress, hash, sig []byte) error {
	pubk, ok := vv.publicKeyMap[author]
	if !ok {
		return VerifyErrUnknownAuthor
	}
	ok = ed25519.Verify(ed25519.PublicKey(pubk.ConsensusPubkey), hash, sig)
	if !ok {
		return VerifyErrInvalidSignature
	}
	return nil
}

// Verify a LedgerInfoWithSignatures
func (vv *ValidatorVerifier) Verify(li *LedgerInfoWithSignatures) error {
	li0 := li.Value.(*LedgerInfoWithSignaturesV0)
	hash := li0.LedgerInfo.Hash()
	sigs := li0.Sigs
	if len(sigs) > len(vv.publicKeyMap) {
		return VerifyErrTooManySignatures
	}
	power := uint64(0)
	for author, sig := range sigs {
		err := vv.verifySingle(author, hash, sig)
		if err != nil {
			return err
		}
		power += vv.publicKeyMap[author].ConsensusVotingPower
	}
	if power < vv.quorumPower {
		return VerifyErrTooFewSignatures
	}
	return nil
}

// EpochChangeVerificationRequired returns true in case the given epoch is larger
// than the existing verifier can support.
// In this case the ValidatorChangeProof should be verified and the verifier updated.
func (vv *ValidatorVerifier) EpochChangeVerificationRequired(epoch uint64) bool {
	return vv.epoch < epoch
}
