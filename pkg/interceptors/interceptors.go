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

package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	jwtauth "github.com/onosproject/onos-lib-go/pkg/auth/jwt"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

const (
	ContextMetadataTokenKey = "bearer"
)

var log = logging.GetLogger("interceptors")

// AuthenticationInterceptor
func AuthenticationInterceptor(ctx context.Context) (context.Context, error) {
	log.Info("authentication interceptor")
	// Extract token from metadata in the context
	tokenString, err := grpc_auth.AuthFromMD(ctx, ContextMetadataTokenKey)
	if err != nil {
		return nil, err
	}

	// Parse token to extract a jwt token
	token, err := jwtauth.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid %d", codes.Unauthenticated)
	}
	return ctx, nil

}
