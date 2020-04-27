package util

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tranvictor/ethashproof"
	"github.com/tranvictor/ethashproof/ethash"
	"github.com/tranvictor/ethashproof/mtree"
)

type DoubleNodeWithMerkleProof struct {
	DagNodes []string `json:"dag_nodes"`
	Proof    []string `json:"proof"`
}

// This struct is used for process interaction
type ProofOutput struct {
	HeaderRLP    string   `json:"header_rlp"`
	MerkleRoot   string   `json:"merkle_root"`
	Elements     []string `json:"elements"`
	MerkleProofs []string `json:"merkle_proofs"`
	ProofLength  uint64   `json:"proof_length"`
}

// Format ProofOutput to double node with merkle proofs
func (o *ProofOutput) Format() []DoubleNodeWithMerkleProof {
	h512s := Filter(o.Elements, func(i int, _ string) bool {
		return i%2 == 0
	})

	h512s = Map(h512s, func(i int, v string) string {
		return v + o.Elements[(i*2)+1][1:]
	})

	dnmps := []DoubleNodeWithMerkleProof{}
	sh512s := Filter(h512s, func(i int, _ string) bool {
		return i%2 == 0
	})
	Map(sh512s, func(i int, v string) string {
		dnmps = append(dnmps, DoubleNodeWithMerkleProof{
			[]string{v, h512s[i*2+1]},
			o.MerkleProofs[uint64(i)*o.ProofLength : (uint64(i)+1)*o.ProofLength],
		})
		return v
	})

	return dnmps
}

//
// Proof eth blockheader
func Proof(header *types.Header) (ProofOutput, error) {
	blockno := header.Number.Uint64()
	epoch := blockno / 30000
	output := &ProofOutput{}

	// get proof from cache
	cache, err := ethashproof.LoadCache(int(epoch))
	if err != nil {
		_, err = ethashproof.CalculateDatasetMerkleRoot(epoch, true)
		if err != nil {
			return *output, err
		}

		cache, err = ethashproof.LoadCache(int(epoch))
		if err != nil {
			return *output, err
		}
	}

	rlpheader, err := ethashproof.RLPHeader(header)
	if err != nil {
		return *output, err
	}

	// init output
	output = &ProofOutput{
		HeaderRLP:    hexutil.Encode(rlpheader),
		MerkleRoot:   cache.RootHash.Hex(),
		Elements:     []string{},
		MerkleProofs: []string{},
		ProofLength:  cache.ProofLength,
	}

	// get verification indices
	indices := ethash.Instance.GetVerificationIndices(
		blockno,
		ethash.Instance.SealHash(header),
		header.Nonce.Uint64(),
	)

	// fill output
	for _, index := range indices {
		// calculate proofs
		element, proof, err := ethashproof.CalculateProof(blockno, index, cache)
		if err != nil {
			return *output, err
		}

		es := element.ToUint256Array()
		for _, e := range es {
			output.Elements = append(output.Elements, hexutil.EncodeBig(e))
		}

		// append proofs
		allProofs := []*big.Int{}
		for _, be := range mtree.HashesToBranchesArray(proof) {
			allProofs = append(allProofs, be.Big())
		}

		for _, pr := range allProofs {
			output.MerkleProofs = append(output.MerkleProofs, hexutil.EncodeBig(pr))
		}
	}

	return *output, nil
}
