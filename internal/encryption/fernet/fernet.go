package fernet

import "github.com/fernet/fernet-go"

type FernetCipherator struct {
}

func NewFernetCipherator() *FernetCipherator {
	return &FernetCipherator{}
}

func (fc *FernetCipherator) Encrypt(msg []byte, key string) ([]byte, error) {
	tokKey, err := fernet.DecodeKey(key)
	if err != nil {
		return nil, err
	}

	return fernet.EncryptAndSign(msg, tokKey)
}
func (fc *FernetCipherator) Decrypt(msg []byte, key string) ([]byte, error) {
	tokKey, err := fernet.DecodeKey(key)
	if err != nil {
		return nil, err
	}

	tokKeys := []*fernet.Key{tokKey}

	return fernet.VerifyAndDecrypt(msg, 0, tokKeys), nil
}
