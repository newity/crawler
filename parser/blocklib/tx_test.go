/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package blocklib

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var (
	tx        Tx
	invalidtx Tx
)

func TestMain(m *testing.M) {
	txsvalid, err := readTxsFromBlock("./mock/sampleblock.pb")
	if err != nil {
		log.Error(err)
	}
	tx = txsvalid[0]

	txsinvalid, err := readTxsFromBlock("./mock/mvcc_read_conflict.pb")
	if err != nil {
		log.Error(err)
	}
	invalidtx = txsinvalid[0]

	m.Run()
}

func readTxsFromBlock(pathToBlock string) ([]Tx, error) {
	file, err := ioutil.ReadFile(pathToBlock)
	if err != nil {
		return nil, err
	}

	fabBlock := &common.Block{}
	err = proto.Unmarshal(file, fabBlock)
	if err != nil {
		return nil, err
	}

	block, err := FromFabricBlock(fabBlock)
	if err != nil {
		return nil, err
	}

	txs, err := block.GetTxs()
	return txs, err
}

func TestIsValid(t *testing.T) {
	t.Run("check valid", func(t *testing.T) {
		assert.Equal(t, true, tx.IsValid())
	})
	t.Run("check invalid", func(t *testing.T) {
		assert.Equal(t, false, invalidtx.IsValid())
	})
}

func TestGetValidationCode(t *testing.T) {
	t.Run("check code 0", func(t *testing.T) {
		assert.Equal(t, int32(0), tx.ValidationCode())
	})
	t.Run("check code 11", func(t *testing.T) {
		assert.Equal(t, int32(11), invalidtx.ValidationCode())
	})
}

func TestGetValidationStatus(t *testing.T) {
	t.Run("check VALID status", func(t *testing.T) {
		assert.Equal(t, "VALID", tx.ValidationStatus())
	})
	t.Run("check MVCC_READ_CONFLICT status", func(t *testing.T) {
		assert.Equal(t, "MVCC_READ_CONFLICT", invalidtx.ValidationStatus())
	})
}

func TestGetEnvelope(t *testing.T) {
	envelope, err := tx.Envelope()
	assert.NoError(t, err)
	assert.NotNil(t, envelope.Payload)
	assert.NotNil(t, envelope.Signature)
}

func TestGetPayload(t *testing.T) {
	payload, err := tx.Payload()
	assert.NoError(t, err)
	assert.NotNil(t, payload.Header)
	assert.NotNil(t, payload.Data)
}

func TestGetChannelHeader(t *testing.T) {
	channelHeader, err := tx.ChannelHeader()
	assert.NoError(t, err)
	assert.NotNil(t, channelHeader)
}

func TestGetSignatureHeader(t *testing.T) {
	signatureHeader, err := tx.SignatureHeader()
	assert.NoError(t, err)
	assert.NotNil(t, signatureHeader)
}

func TestGetChaincodeId(t *testing.T) {
	id, err := tx.ChaincodeId()
	assert.NoError(t, err)
	assert.Equal(t, "fabcar", id.Name)
}

func TestGetEpoch(t *testing.T) {
	epoch, err := tx.Epoch()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), epoch)
}

func TestTimestamp(t *testing.T) {
	timestamp, err := tx.Timestamp()
	assert.NoError(t, err)
	assert.Equal(t, int64(1603659829097237404), timestamp.UnixNano())
}

func TestGetTxId(t *testing.T) {
	txid, err := tx.TxId()
	assert.NoError(t, err)
	assert.Equal(t, "23e7c409b6849a71e6b5d7767a4e6c7efcd4bafba02b932ca5e6559e4d050dea", txid)
}

func TestGetPeerTransaction(t *testing.T) {
	peerTransaction, err := tx.PeerTransaction()
	assert.NoError(t, err)
	assert.NotNil(t, peerTransaction.Actions)
}

func TestGetActions(t *testing.T) {
	actions, err := tx.Actions()
	assert.NoError(t, err)
	for _, action := range actions {
		payloadHash := sha256.New()
		payloadHash.Write(action.Payload.ChaincodeProposalPayload)
		creatorHash := sha256.New()
		creatorHash.Write(action.SignatureHeader.Creator)
		assert.Equal(t, "51353a437c811a1ec7d4ffe061d1f38907ad443c5c1847b4877c4a65c5efa24e", hex.EncodeToString(payloadHash.Sum(nil)))
		assert.Equal(t, "3b2106648e7b0773db03d160dbfef48a514f0871f8e18524a10a2de19fb21dd9", hex.EncodeToString(creatorHash.Sum(nil)))
		assert.Equal(t, "2693c3593e74c1984b12aca0dcd625619bdbb59fe05abf76", hex.EncodeToString(action.SignatureHeader.Nonce))
	}
}