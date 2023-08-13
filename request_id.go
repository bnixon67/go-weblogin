/*
Copyright 2023 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package weblogin

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

// Key to use when setting the request ID.
type ctxReqIDKey int

// ReqIDKey is the key for the unique request ID in a request context.
const ReqIDKey ctxReqIDKey = 0

var (
	reqIDPrefix string // random prefix for the request id
	reqIDVal    uint64 // current request id value
)

// init generates a unique required id prefix each time the server starts.
func init() {
	reqIDPrefix, _ = GenerateRandomString(6)
}

// generateRequestID generates a unique request ID.
func generateRequestID() string {
	id := atomic.AddUint64(&reqIDVal, 1)
	return fmt.Sprintf("%s%06d", reqIDPrefix, id)
}

// GetReqID returns the request ID from ctx if present, otherwise "".
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	reqID, ok := ctx.Value(ReqIDKey).(string)
	if !ok {
		return ""
	}
	return reqID
}

// RequestIDHandler is middleware that adds a unique request ID to the request context.
func RequestIDHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()
		ctx := context.WithValue(r.Context(), ReqIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
