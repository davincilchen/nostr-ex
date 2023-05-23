package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"
	// btcec "github.com/btcsuite/btcd/btcec/v2"
	// schnorr "github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type Msg struct { //TODO:
	ID        string          `json:"id"`
	Pubkey    string          `json:"pubkey"`
	CreatedAt int64           `json:"created_at"`
	Kind      int             `json:"kind"`
	Content   string          `json:"content"`
	Tags      [][]interface{} `json:"tags"`
	Sig       string          `json:"sig"`
}

func NewMsg(Pubkey, Content string) *Msg {
	return &Msg{
		Pubkey:  Pubkey,
		Content: Content,
	}
}

func (t *Msg) MakeEvent(privateKeyStr string) ([]interface{}, error) {
	t.Kind = 1
	t.CreatedAt = time.Now().Unix()

	serializedData, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(serializedData)

	// convert the hash value to a lowercase hex-encoded string
	t.ID = hex.EncodeToString(hash[:])

	// .. //
	tmp, err := t.Sign(privateKeyStr, hash)
	if err != nil {
		return nil, err
	}

	t.Sig = tmp
	return []interface{}{"EVENT", t}, nil
}

func (t *Msg) Sign(privateKeyStr string, dataHash [32]byte) (string, error) {

	// 64 HEX private key
	privateKeyHex := privateKeyStr

	// decode the private key hex string
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)

	// create an ecdsa.PrivateKey from the raw bytes
	curve := elliptic.P256()
	privateKey := &ecdsa.PrivateKey{D: new(big.Int).SetBytes(privateKeyBytes), PublicKey: ecdsa.PublicKey{Curve: curve}}

	// sign the hash using the private key
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, dataHash[:])
	if err != nil {
		return "", err
	}

	signature := append(r.Bytes(), s.Bytes()...)

	//將簽名轉換為 base64 字符串
	ret := hex.EncodeToString(signature)

	//fmt.Println(ret)

	return ret, nil
}
