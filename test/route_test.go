package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	go_web "github.com/cyuzhe1994-commits/go-web"
	"github.com/cyuzhe1994-commits/go-web/middleware"
)

func routeTest(method string, path string, router *go_web.Router) {
	router.Add(method, path, func(ctx *go_web.Context) {
		ctx.Echo(http.StatusOK, method+" "+ctx.Request.URL.Path)
	})
}

func routeCall(method string, path string, t *testing.T) {
	res, err := http.Get("http://localhost:8080" + path)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	matchString := method + " " + path
	if string(body) != matchString {
		t.Fatalf("Expected response body '%s', got '%s'", matchString, string(body))
	}
}

// TestRoute 测试路由功能，包括路由注册和路由调用，路由分组功能
func TestRoute(t *testing.T) {
	engin := go_web.NewEngine(nil)
	router := engin.Router
	routeTest(http.MethodGet, "/test", router)
	routeTest(http.MethodGet, "/test/:id", router)
	v1 := router.Group("/api/v1")
	routeTest(http.MethodGet, "/test", v1)
	healthy := v1.Group("/healthy")
	routeTest(http.MethodGet, "/check", healthy)

	go engin.Run(":8080")

	routeCall(http.MethodGet, "/test", t)
	routeCall(http.MethodGet, "/test/123", t)
	routeCall(http.MethodGet, "/api/v1/test", t)
	routeCall(http.MethodGet, "/api/v1/healthy/check", t)
}

func TestRouteBindJson(t *testing.T) {
	engin := go_web.NewEngine(nil)
	router := engin.Router
	router.Add(http.MethodPost, "/test/bindjson", func(ctx *go_web.Context) {
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var user User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.Echo(http.StatusBadRequest, "Invalid JSON")
			return
		}
		ctx.Echo(http.StatusOK, fmt.Sprintf("Received user: %s, age: %d", user.Name, user.Age))
	})

	go engin.Run(":8080")

	// 发送 POST 请求测试 BindJSON
	jsonData := `{"name": "Alice", "age": 30}`
	res, err := http.Post("http://localhost:8080/test/bindjson",
		"application/json",
		io.NopCloser(io.Reader(strings.NewReader(jsonData))),
	)
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	expectedResponse := "Received user: Alice, age: 30"
	if string(body) != expectedResponse {
		t.Fatalf("Expected response body '%s', got '%s'", expectedResponse, string(body))
	}
}

func TestRouteMiddlewarePanic(t *testing.T) {
	engin := go_web.NewEngine(nil)
	engin.Use(middleware.Recovery)
	engin.Use(middleware.Logger)
	router := engin.Router
	router.Use(middleware.Cors)
	router.Get("/healthy", func(ctx *go_web.Context) {
		panic("healthy")
	})

	go engin.Run(":8080")
	res, err := http.Get("http://localhost:8080/healthy")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	} else if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code 500, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	expectedResponse := "healthy"
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if errorResponse.Error != expectedResponse {
		t.Fatalf("Expected response body %s", string(body))
	}
}

func TestRouteMiddlewareSingleRoute(t *testing.T) {
	engin := go_web.NewEngine(nil)
	router := engin.Router
	router.Use(middleware.Cors)
	router.Get("/healthy", func(ctx *go_web.Context) {
		ctx.Echo(http.StatusOK, "healthy")
	}, func(next go_web.HandlerFunc) go_web.HandlerFunc {
		return func(c *go_web.Context) {
			c.Writer.Header().Set("Single-Middleware", "success")
			next(c)
		}
	})
	router.Get("/healthy001", func(ctx *go_web.Context) {
		ctx.Echo(http.StatusOK, "healthy001")
	})
	go engin.Run(":8080")
	res, err := http.Get("http://localhost:8080/healthy")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	expectedResponse := "healthy"

	if string(body) != expectedResponse {
		t.Fatalf("Expected response body %s", string(body))
	} else {
		singleMiddlewareHeader := res.Header.Get("Single-Middleware")
		if singleMiddlewareHeader != "success" {
			t.Fatalf("Expected Single-Middleware header to be 'success', got '%s'", singleMiddlewareHeader)
		}
	}

	res, err = http.Get("http://localhost:8080/healthy001")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	body, _ = io.ReadAll(res.Body)
	res.Body.Close()
	expectedResponse = "healthy001"

	if string(body) != expectedResponse {
		t.Fatalf("Expected response body %s", string(body))
	} else {
		singleMiddlewareHeader := res.Header.Get("Single-Middleware")
		if singleMiddlewareHeader == "success" {
			t.Fatalf("Expected Single-Middleware header to be 'success', got '%s'", singleMiddlewareHeader)
		}
	}
}
