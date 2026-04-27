package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Claims struct {
	Hash string `json:"hash"`
	Exp  int64  `json:"exp"`
}

func generateToken(password string) (string, error) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	claims := Claims{
		Hash: hash,
		Exp:  time.Now().Add(8 * time.Hour).Unix(),
	}

	header := `{"alg":"HS256","typ":"JWT"}`
	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(header))

	payloadBytes, _ := json.Marshal(claims)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signatureInput := headerB64 + "." + payloadB64
	h := hmac.New(sha256.New, []byte(password))
	h.Write([]byte(signatureInput))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return headerB64 + "." + payloadB64 + "." + signature, nil
}

func validateToken(token, password string) bool {
	if token == "" || password == "" {
		return false
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	signatureInput := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(password))
	h.Write([]byte(signatureInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if parts[2] != expectedSig {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return false
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if claims.Hash != hash {
		return false
	}

	return claims.Exp > time.Now().Unix()
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			writeJSON(w, map[string]string{"error": "Authentication required"}, http.StatusUnauthorized)
			return
		}

		if !validateToken(cookie.Value, pass) {
			writeJSON(w, map[string]string{"error": "Authentication required"}, http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
