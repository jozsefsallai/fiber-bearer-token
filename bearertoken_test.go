package bearertoken

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/gofiber/fiber"
)

func firstTestCase(app *fiber.App) string {
	// It should not return a token if none is provided
	req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
	res, _ := app.Test(req)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return string(body)
}

func secondTestCase(app *fiber.App, key string) string {
	// It should return token from header

	if len(key) == 0 {
		key = "Bearer"
	}

	req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s nanachi-cute", key))
	res, _ := app.Test(req)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return string(body)
}

func thirdTestCase(app *fiber.App, key string) string {
	// It should return token from param

	if len(key) == 0 {
		key = "access_token"
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8000/?%s=nanachi-cute", key), nil)
	res, _ := app.Test(req)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return string(body)
}

func fourthTestCase(app *fiber.App, key string) string {
	// It should return token from body

	if len(key) == 0 {
		key = "access_token"
	}

	payload := []byte(fmt.Sprintf("%s=nanachi-cute", key))

	req, _ := http.NewRequest("POST", "http://localhost:8000", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	res, _ := app.Test(req)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return string(body)
}

func fifthTestCase(app *fiber.App, queryKey, headerKey string) int {
	// It should return HTTP status 400 if token was provided multiple times

	if len(queryKey) == 0 {
		queryKey = "access_token"
	}

	if len(headerKey) == 0 {
		headerKey = "Bearer"
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8000/?%s=nanachi-cute", queryKey), nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s nanachi-cute", headerKey))

	res, _ := app.Test(req)
	return res.StatusCode
}

func TestMiddleware(t *testing.T) {
	testCases := []Config{
		Config{},
		Config{
			BodyKey:    "custom_token",
			HeaderKey:  "custom_token",
			QueryKey:   "custom_token",
			RequestKey: "bearer",
		},
	}
	testCaseLabels := []string{
		"with default settings",
		"with custom settings",
	}

	for idx, config := range testCases {
		t.Run(testCaseLabels[idx], func(t *testing.T) {
			app := fiber.New()

			app.Use(New(&config))

			var requestKey string = "token"
			if len(config.RequestKey) > 0 {
				requestKey = config.RequestKey
			}

			app.Get("/", func(ctx *fiber.Ctx) {
				ctx.Send(ctx.Locals(requestKey))
			})

			app.Post("/", func(ctx *fiber.Ctx) {
				ctx.Send(ctx.Locals(requestKey))
			})

			first := firstTestCase(app)
			if first != "" {
				t.Errorf(`expected: <empty string>, got: "%s"`, first)
			}

			second := secondTestCase(app, config.HeaderKey)
			if second != "nanachi-cute" {
				t.Errorf(`expected: "nanachi-cute", got: "%s"`, second)
			}

			third := thirdTestCase(app, config.QueryKey)
			if third != "nanachi-cute" {
				t.Errorf(`expected: "nanachi-cute", got: "%s"`, third)
			}

			fourth := fourthTestCase(app, config.BodyKey)
			if fourth != "nanachi-cute" {
				t.Errorf(`expected: "nanachi-cute", got: "%s"`, fourth)
			}

			fifth := fifthTestCase(app, config.QueryKey, config.HeaderKey)
			if fifth != 400 {
				t.Errorf(`expected: 400, got: %d`, fifth)
			}
		})
	}
}
