// Copyright 2018 Drone.IO Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package header

import (
	"net/http"
	"time"
	"github.com/labstack/echo"
	"log"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.

func NoCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
			c.Request().Header.Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
			c.Request().Header.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			return next(c)
		}
	}
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != "OPTIONS" {
				return next(c)
			} else {
				log.Println("get OPTIONS request")

				c.Request().Header.Set("Access-Control-Allow-Origin", "*")
				c.Request().Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
				c.Request().Header.Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
				c.Request().Header.Set("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
				c.Request().Header.Set("Content-Type", "application/json")
				return c.NoContent(http.StatusNoContent)
			}
		}
	}
}


// Secure is a middleware function that appends security
// and resource access headers.
func Secure() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set("Access-Control-Allow-Origin", "*")
			c.Request().Header.Set("X-Frame-Options", "DENY")
			c.Request().Header.Set("X-Content-Type-Options", "nosniff")
			c.Request().Header.Set("X-XSS-Protection", "1; mode=block")
			return next(c)
		}
	}
}