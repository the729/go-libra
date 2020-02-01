package types

import (
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/ed25519"
)

var (
	VerifyErrUnknownAuthor     = errors.New("unknown author")
	VerifyErrInvalidSignature  = errors.New("invalid signature")
	VerifyErrTooFewSignatures  = errors.New("too few signatures")
	VerifyErrTooManySignatures = errors.New("too many signatures")
)

type ValidatorVerifier struct {
	publicKeyMap map[string]*ValidatorPublicKeys
	epoch        uint64
	totalPower   uint64
	quorumPower  uint64
}

func (vv *ValidatorVerifier) FromValidatorSet(vs ValidatorSet, epoch uint64) error {
	vv.publicKeyMap = make(map[string]*ValidatorPublicKeys)
	vv.totalPower = 0
	for _, v := range vs {
		vv.publicKeyMap[hex.EncodeToString(v.AccountAddress[:])] = &ValidatorPublicKeys{
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

func (vv *ValidatorVerifier) ToValidatorSet() (ValidatorSet, uint64) {
	vs := make([]*ValidatorPublicKeys, 0, len(vv.publicKeyMap))
	for addr, v := range vv.publicKeyMap {
		vpk := &ValidatorPublicKeys{
			ConsensusPubkey:      cloneBytes(v.ConsensusPubkey),
			ConsensusVotingPower: v.ConsensusVotingPower,
		}
		hex.Decode(vpk.AccountAddress[:], []byte(addr))
		vs = append(vs, vpk)
	}
	return vs, vv.epoch
}

func (vv *ValidatorVerifier) verifySingle(author string, hash, sig []byte) error {
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

func (vv *ValidatorVerifier) Verify(li *LedgerInfoWithSignatures) error {
	hash := li.LedgerInfo.Hash()
	sigs := li.Sigs
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
