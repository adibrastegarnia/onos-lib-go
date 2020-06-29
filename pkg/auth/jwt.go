// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/codes"

	"github.com/dgrijalva/jwt-go"
)

const (
	// SharedSecretKey shared secret key for signing a token
	SharedSecretKey = "SHARED_SECRET_KEY"
)

// JwtAuthenticator jwt authenticator
type JwtAuthenticator struct {
}

// NewJwtAuthenticator create an instance of jwt authenticator
func NewJwtAuthenticator() *JwtAuthenticator {
	return &JwtAuthenticator{}
}

// ParseToken parse token and Ensure that the JWT conforms to the structure of a JWT.
func (j *JwtAuthenticator) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// HS256, HS384, or HS512
		if strings.HasPrefix(token.Method.Alg(), "HS") {
			key := os.Getenv(SharedSecretKey)
			return []byte(key), nil
		}
		return nil, fmt.Errorf("unknown signining algorithm: %s", token.Method.Alg())
	})

	return token, claims, err

}

// Authenticate parse a jwt string token and authenticate it
func (j *JwtAuthenticator) Authenticate(tokenString string) (jwt.MapClaims, error) {
	token, claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid %d", codes.Unauthenticated)
	}

	return claims, err
}
