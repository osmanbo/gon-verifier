package verifier

import (
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
)

type A5Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassID           string // Wasm Contract Addr
	NFTID             string
	ChainID           string // Dest Chain ID
}

type A5Verifier struct {
	r *chain.Registry
}

func (v A5Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A5Params)
	if !ok {
		result.Reason = ReasonIllegalParams
		res <- result
		return
	}

	if len(params.TxHash) == 0 {
		result.Reason = ReasonTxHashEmpty
		res <- result
		return
	}

	srcChain := v.r.GetChain(params.ChainAbbreviation)
	tx, err := srcChain.GetTx(params.TxHash)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	if req.User.Address[params.ChainAbbreviation] != tx.Sender {
		result.Reason = ReasonSenderNotMatch
		res <- result
		return
	}

	// FIXME: receiver must be the iris addr

	if !srcChain.HasClass(params.ClassID) {
		result.Reason = ReasonClassNotFound
		res <- result
		return
	}

	// FIXME: check original class on iris
	originalClassID := ""
	iris := v.r.GetChain(chain.ChainIdAbbreviationIris)
	nft, err := iris.GetNFT(originalClassID, params.NFTID)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	if nft.Owner != req.User.Address[chain.ChainIdAbbreviationIris] {
		result.Reason = ReasonNFTOwnerNotMatch
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A5Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) != 1 {
		return nil, errors.New("rows length not match")
	}

	param := rows[0]
	chainAbbr := chain.ChainIdAbbreviationStars
	if param[3] == chain.ChainIdValueJuno {
		chainAbbr = chain.ChainIdAbbreviationJuno
	}

	return A5Params{
		ChainAbbreviation: chainAbbr,
		TxHash:            param[0],
		ClassID:           param[1], // Wasm Contract Addr
		NFTID:             param[2],
		ChainID:           param[3],
	}, nil
}
