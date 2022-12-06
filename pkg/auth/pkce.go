package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

var Empty Params

const (
	S256 = "S256"
)

type Params struct {
	Challenge string
	Method    string
	Verifier  string
}

func (p Params) IsEmpty() bool {
	return p == Params{}
}

func NewPKCEParams(methods []string) (Params, error) {
	for _, method := range methods {
		if method == S256 {
			return newS256()
		}
	}
	return Empty, nil
}

func newS256() (Params, error) {
	b, err := randomBytes(32)
	if err != nil {
		return Empty, fmt.Errorf("error generating random number: %w", err)
	}
	return S256From(b), nil
}

func randString() (string, error) {
	b, err := randomBytes(32)
	if err != nil {
		return "", err
	}
	return base64URLEncode(b), nil
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if err := binary.Read(rand.Reader, binary.LittleEndian, b); err != nil {
		return nil, err
	}
	return b, nil
}

func S256From(b []byte) Params {
	v := base64URLEncode(b)
	s := sha256.New()
	_, _ = s.Write([]byte(v))
	return Params{
		Challenge: base64URLEncode(s.Sum(nil)),
		Method:    S256,
		Verifier:  v,
	}
}

func base64URLEncode(b []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

func getSupportedPKCEMethods(p *oidc.Provider) ([]string, error) {
	var claims struct {
		Supported []string `json:"code_challenge_methods_supported"`
	}
	if err := p.Claims(&claims); err != nil {
		return nil, fmt.Errorf("couldn't parse code challenge method claims: %w", err)
	}
	return claims.Supported, nil
}
