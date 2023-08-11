package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/swe"
)

type RewardUserAddress struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Addr  string `json:"addr"`
}

func DecryptUserAddress(ctx *swe.Context, pubKey []byte, addr string) (*RewardUserAddress, error) {
	if len(addr) == 0 {
		return &RewardUserAddress{}, nil
	}

	priKey, ok := getPrivateKeyFromContext(ctx)
	if !ok || priKey == nil {
		return &RewardUserAddress{}, nil
	}

	key, err := ECDH(priKey, pubKey)
	if err != nil {
		return nil, err
	}

	data, err := Decrypt(key, addr)
	if err != nil {
		return nil, err
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	ret := RewardUserAddress{}
	err = json.Unmarshal(data, &ret)
	return &ret, err
}

func getPrivateKeyFromContext(ctx *swe.Context) ([]byte, bool) {
	if st, ok := session.GetStreamerSession(ctx); ok {
		return st.PrivateKey, true
	}
	if user, ok := session.GetDDSession(ctx); ok {
		return user.PrivateKey, true
	}
	return nil, false
}

func EncryptUserAddress(ctx *swe.Context, pubKey []byte, addr *RewardUserAddress) (string, error) {
	dd, ok := session.GetDDSession(ctx)
	if !ok || dd == nil {
		return "", nil
	}

	key, err := ECDH(dd.PrivateKey, pubKey)
	if err != nil {
		return "", err
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, _ := json.Marshal(addr)

	return Encrypt(key, data)
}
