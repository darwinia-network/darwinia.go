package core

import (
	// "encoding/hex"
	"fmt"

	"github.com/darwinia-network/darwinia.go/util"
	"github.com/ethereum/go-ethereum/core/types"
)

// Dimmy shadow service
type Shadow int

/**
 * GetEthHeaderByNumber
 */
type GetEthHeaderByNumberParams struct {
	Number uint64 `json:"number"`
}

type GetEthHeaderByNumberResp struct {
	Header types.Header `json:"header"`
}

func (s *Shadow) GetEthHeaderByNumber(
	params GetEthHeaderByNumberParams,
	resp *GetEthHeaderByNumberResp,
) error {
	var err error
	resp.Header, err = util.Header(params.Number)
	return err
}

/**
 * GetEthHeaderWithProofByNumber
 */
type GetEthHeaderWithProofByNumberOptions struct {
	Format string `json:"format"`
}

type GetEthHeaderWithProofByNumberParams struct {
	Number  uint64                               `json:"block_num"`
	Options GetEthHeaderWithProofByNumberOptions `json:"options"`
}

type GetEthHeaderWithProofByNumberRawResp struct {
	Header util.DarwiniaEthHeader           `json:"eth_header"`
	Proof  []util.DoubleNodeWithMerkleProof `json:"proof"`
}

type GetEthHeaderWithProofByNumberJSONResp struct {
	Header util.DarwiniaEthHeaderHexFormat  `json:"eth_header"`
	Proof  []util.DoubleNodeWithMerkleProof `json:"proof"`
}

type GetEthHeaderWithProofByNumberCodecResp struct {
	Header string `json:"header"`
	Proof  string `json:"proof"`
}

func (s *Shadow) GetEthHeaderWithProofByNumber(
	params GetEthHeaderWithProofByNumberParams,
	resp *interface{},
) error {
	// Fetch header from cache
	cache := EthHeaderWithProofCache{Number: params.Number}
	rawResp, err := cache.Fetch()

	fmt.Println("requsting")
	// Fetch header from infura
	if err != nil {
		ethHeader, err := util.Header(params.Number)
		if err != nil {
			return err
		}

		rawResp.Header, err = util.IntoDarwiniaEthHeader(ethHeader)
		if err != nil {
			return err
		}

		// Proof header
		proof, err := util.Proof(&ethHeader)
		rawResp.Proof = proof.Format()
		if err != nil {
			return err
		}

		// Create cache
		err = cache.FromResp(rawResp)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s\n", rawResp)
	// Set response
	*resp = rawResp

	// Check if need codec
	if params.Options.Format == "scale" {
		*resp = GetEthHeaderWithProofByNumberCodecResp{
			encodeDarwiniaEthHeader(rawResp.Header),
			encodeProofArray(rawResp.Proof),
		}
	} else if params.Options.Format == "json" {
		*resp = GetEthHeaderWithProofByNumberJSONResp{
			rawResp.Header.HexFormat(),
			rawResp.Proof,
		}
	}

	return nil
}
