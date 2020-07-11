package core

import (
	"fmt"
	"strings"

	"github.com/darwinia-network/shadow/internal"
	"github.com/darwinia-network/shadow/internal/eth"
	"github.com/darwinia-network/shadow/internal/util"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
)

// Shadow genesis block error message
const GENESIS_ERROR = "The requested block number is too low, only support blocks heigher than %v"
const PROOF_LOCK = "proof.lock"

// Dimmy shadow service
type Shadow struct {
	Config internal.Config
	DB     *gorm.DB
}

func NewShadow() (Shadow, error) {
	conf := new(internal.Config)
	err := conf.Load()
	if err != nil {
		return Shadow{}, err
	}

	db, err := ConnectDb()
	if err != nil {
		return Shadow{}, err
	}

	return Shadow{
		*conf,
		db,
	}, err
}

/**
 * Genesis block checker
 */
func (s *Shadow) checkGenesis(genesis uint64, block interface{}) (uint64, error) {
	block, err := util.NumberOrString(block)
	if err != nil {
		return genesis, err
	}

	switch b := block.(type) {
	case uint64:
		if b < genesis {
			return genesis, fmt.Errorf(GENESIS_ERROR, genesis)
		}

		return b, nil
	case string:
		eH, err := eth.Header(b, s.Config.Api)
		if err != nil {
			return genesis, err
		}

		// convert ethHeader to darwinia header
		dH, err := eth.IntoDarwiniaEthHeader(eH)
		if err != nil {
			return dH.Number, err
		}

		// Check hash empty response
		if util.IsEmpty(dH) {
			return genesis, fmt.Errorf("Empty block: %s", b)
		}

		// Check genesis by number
		if dH.Number <= genesis {
			return genesis, fmt.Errorf(GENESIS_ERROR, genesis)
		}

		return dH.Number, nil
	default:
		return genesis, fmt.Errorf("Invaild block param: %v", block)
	}
}

/**
 * GetEthHeaderByNumber
 */
func (s *Shadow) GetHeader(
	chain Chain,
	block interface{},
) (types.Header, error) {
	switch chain {
	default:
		num, err := s.checkGenesis(s.Config.Genesis, block)
		if err != nil {
			return types.Header{}, err
		}

		return eth.Header(num, s.Config.Api)
	}

}

func (s *Shadow) GetHeaderWithProof(
	chain Chain,
	block interface{},
) (GetEthHeaderWithProofRawResp, error) {
	switch chain {
	default:
		num, err := s.checkGenesis(s.Config.Genesis, block)
		if err != nil {
			return GetEthHeaderWithProofRawResp{}, err
		}

		// Fetch header from cache
		cache := EthHeaderWithProofCache{Number: num}
		err = cache.Fetch(s.Config, s.DB)
		if err != nil {
			return GetEthHeaderWithProofRawResp{}, err
		}

		err = cache.ApplyProof(s.Config, s.DB)
		if err != nil {
			return GetEthHeaderWithProofRawResp{}, err
		}

		rawResp, err := cache.IntoResp()
		if err != nil {
			return GetEthHeaderWithProofRawResp{}, err
		}

		return rawResp, nil
	}
}

/**
 * BatchEthHeaderWithProofByNumber
 */
func (s *Shadow) BatchHeaderWithProof(
	block uint64,
	batch int,
) ([]GetEthHeaderWithProofRawResp, error) {
	var nps []GetEthHeaderWithProofRawResp
	for i := 0; i < batch; i++ {
		np, err := s.GetHeaderWithProof(
			Ethereum,
			block+uint64(i),
		)

		if err != nil {
			return nps, err
		}

		nps = append(nps, np)
	}

	return nps, nil
}

/**
 * Get proposal headers
 */
func (s *Shadow) GetProposalHeaders(
	numbers []uint64,
	format ProofFormat,
) ([]ProposalHeader, error) {
	var (
		phs []ProposalHeader
	)

	for _, i := range numbers {
		rawp, err := s.GetHeaderWithProof(
			Ethereum,
			uint64(i),
		)
		if err != nil {
			return phs, err
		}

		phs = append(phs, ProposalHeader{
			Header: rawp.Header,
			Proof:  rawp.Proof,
			Root:   rawp.Root,
		})
	}

	return phs, nil
}

/**
 * Get proposal headers
 */
func (s *Shadow) GetReceipt(
	tx string,
) (resp GetReceiptResp, err error) {
	proof, hash, err := eth.GetReceipt(tx)
	if err != nil {
		return
	}

	resp.ReceiptProof = proof
	cache := EthHeaderWithProofCache{Hash: hash}
	err = cache.Fetch(s.Config, s.DB)
	if err != nil {
		return
	}

	resp.MMRProof = strings.Split(cache.MMRProof, ",")
	return
}
