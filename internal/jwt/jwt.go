package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
}

type Claims struct {
	Sub  int64  `json:"sub"`
	Name string `json:"name"`
	Iat  int64  `json:"iat"`
}

type JWT struct {
	secret string
}

func New(secret string) *JWT {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) Encode(claims Claims) (string, error) {
	jsonHeader, err := json.Marshal(Header{
		Alg: "HS256",
	})
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}

	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	header := base64Encode(jsonHeader)
	payload := base64Encode(jsonClaims)
	unsignedBase64 := header + "." + payload

	signature := j.sign(unsignedBase64, j.secret)
	signatureBase64 := base64Encode(signature)

	return unsignedBase64 + "." + signatureBase64, nil
}

func (j *JWT) sign(content, secret string) []byte {
	hasher := hmac.New(sha256.New, []byte(secret))

	if _, err := hasher.Write([]byte(content)); err != nil {
		panic(err)
	}

	return hasher.Sum(nil)
}

func (j *JWT) Verify(token string) error {
	if token == "" {
		return fmt.Errorf("empty token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("token does not have 3 parts")
	}

	jsonHeader, err := base64Decode(parts[0])
	if err != nil {
		return fmt.Errorf("decode json header: %w", err)
	}

	var header Header
	if err := json.Unmarshal(jsonHeader, &header); err != nil {
		return fmt.Errorf("unmarshal json header: %w", err)
	}

	if header.Alg != "HS256" {
		return errors.New("expected alg HS256")
	}

	signature := j.sign(parts[0]+"."+parts[1], j.secret)
	signatureBase64 := base64Encode(signature)
	if signatureBase64 != parts[2] {
		return fmt.Errorf("invalid signature")
	}

	jsonPayload, err := base64Decode(parts[1])
	if err != nil {
		return fmt.Errorf("decode json payload: %w", err)
	}

	var payload Claims
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		return fmt.Errorf("unmarshal json payload: %w", err)
	}

	now := time.Now().UTC().Unix()
	if payload.Iat > 0 && payload.Iat > now {
		return fmt.Errorf("used before issued: %d", now)
	}

	return nil
}

func (j *JWT) Decode(token string) (*Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("token does not have 3 parts")
	}

	jsonPayload, err := base64Decode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode json payload: %w", err)
	}

	var payload Claims
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal json payload: %w", err)
	}

	return &payload, nil
}

func (j *JWT) VerifyAndDecode(token string) (*Claims, error) {
	if err := j.Verify(token); err != nil {
		return nil, err
	}

	return j.Decode(token)
}

func base64Encode(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

func base64Decode(src string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(src)
}
