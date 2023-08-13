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
	"log/slog"
	"net/http"
)

// LogRequestHandler is middleware that logs all HTTP requests.
func LogRequestHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("LogRequest",
			slog.Group("request",
				slog.String("id", GetReqID(r.Context())),
				slog.String("remoteAddr", GetRealRemoteAddr(r)),
				slog.String("method", r.Method),
				slog.String("url", r.RequestURI),
			),
		)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
