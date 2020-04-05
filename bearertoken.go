// Package bearertoken is a middleware for the Fiber web framework that allows you to
// extract bearer authorization tokens from the HTTP requests sent to your application.
// The middleware is RFC6750-compliant, however, it does allow you to specify
// non-compliant settings.
//
// Quick Start
//
//   package main
//
//   import (
//     "github.com/gofiber/fiber"
//     bearertoken "github.com/jozsefsallai/fiber-bearer-token"
//   )
//
//   func main() {
//     app := fiber.New()
//     app.Use(bearertoken.New(nil))
//     app.Listen(3000)
//   }
//
// You can access the bearer token from the request's context using the designated local
// variable. By default, the variable is called "token", but you can change it to anything.
//
//   app.Get("/", func (ctx *fiber.Ctx)) {
//     bearer := ctx.Locals("token")
//     if len(bearer) == 0 {
//       ctx.Send("Unauthorized!")
//     } else {
//       ctx.Send("You're in!")
//     }
//   })
//
// The middleware searches for the bearer token inside the query paramters, then the body
// parameters, and then the authorization header (if its value starts with the specified
// key), in this order. Support for cookies is a planned feature.
//
// As per RFC6750, if the middleware finds more than one bearer tokens, it will return an
// HTTP 400 status code, aborting the request.
//
// You can customize the names of the keys, as well as the name of the local variable during
// the initialization of the middleware:
//
//   app.Use(bearertoken.New(&bearertoken.Config{
//     BodyKey: "auth_token",
//     HeaderKey: "Auth",
//     QueryKey: "auth_token",
//     RequestKey: "bearer"
//   }))
package bearertoken

import (
	"strings"

	"github.com/gofiber/fiber"
)

// Config holds the configuration of the middleware. It is completely optional
// and should only be provided if your application uses token keys that are not
// RFC6750-compliant.
type Config struct {
	// BodyKey defines the key to use when searching for the bearer token inside the
	// request's body.
	// Optional. Default: "access_token".
	BodyKey string

	// HeaderKey defines the prefix of the Authorization header's value, used when
	// searching for the bearer token inside the request's headers.
	// Optional. Default: "Bearer".
	HeaderKey string

	// QueryKey defines the key to use when searching for the bearer token inside the
	// request's query parameters.
	// Optional. Default: "access_token".
	QueryKey string

	// RequestKey defines the name of the local variable that will be created in the
	// request's context, which will contain the bearer token extracted from the
	// request.
	// Optional. Default: "token".
	RequestKey string
}

// New creates a middleware for use in Fiber.
func New(opts *Config) func(*fiber.Ctx) {
	config := &Config{
		BodyKey:    "access_token",
		HeaderKey:  "Bearer",
		QueryKey:   "access_token",
		RequestKey: "token",
	}

	if opts != nil {
		if len(opts.BodyKey) > 0 {
			config.BodyKey = opts.BodyKey
		}

		if len(opts.HeaderKey) > 0 {
			config.HeaderKey = opts.HeaderKey
		}

		if len(opts.QueryKey) > 0 {
			config.QueryKey = opts.QueryKey
		}

		if len(opts.RequestKey) > 0 {
			config.RequestKey = opts.RequestKey
		}
	}

	return func(ctx *fiber.Ctx) {
		var token string
		var errored bool = false

		// query parameter
		queryValue := ctx.Query(config.QueryKey)

		if len(queryValue) > 0 {
			token = queryValue
		}

		// body parameter
		bodyValue := ctx.Body(config.BodyKey)

		if len(bodyValue) > 0 {
			if len(token) > 0 {
				errored = true
			}

			token = bodyValue
		}

		// request authorization header
		headerValue := ctx.Get("authorization")

		if len(headerValue) > 0 {
			components := strings.SplitN(headerValue, " ", 2)

			if len(components) == 2 && components[0] == config.HeaderKey {
				if len(token) > 0 {
					errored = true
				}

				token = components[1]
			}
		}

		if errored {
			ctx.Status(400).Send()
		} else {
			ctx.Locals(config.RequestKey, token)
			ctx.Next()
		}
	}
}
