# fiber-bearer-token

`fiber-bearer-token` is an [RFC6750](https://tools.ietf.org/html/rfc6750)-compliant middleware for the [Fiber web framework](https://github.com/gofiber/fiber) that makes it easy to work with bearer tokens.

## Getting Started

**Install the module:**

```
go get github.com/jozsefsallai/fiber-bearer-token
```

**Include the middleware in your Fiber setup:**

```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/jozsefsallai/fiber-bearer-token"
)

func main() {
  app := fiber.New()
  app.Use(bearertoken.New(nil))
  app.Listen(3000)
}
```

**Read more in the [documentation](https://pkg.go.dev/github.com/jozsefsallai/fiber-bearer-token).**

## TODO

- [ ] Add support for cookie-based bearer tokens

## License

MIT.
