package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
)

//go:embed all:web/dist
var spaFS embed.FS

// mountSPA serves the Vue production bundle; unknown non-/api paths fall
// through to index.html so the client router handles them (login / dashboard).
func mountSPA(app *fiber.App) {
	dist, err := fs.Sub(spaFS, "web/dist")
	if err != nil {
		panic(err)
	}
	app.Get("/", serveIndex(dist))
	app.Get("/*", spaFallback(dist))
	app.Use("/", staticFromFS(dist))
}

func serveIndex(dist fs.FS) fiber.Handler {
	return func(c fiber.Ctx) error {
		b, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			return c.Status(500).SendString("SPA index.html missing in embed")
		}
		c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
		return c.Send(b)
	}
}

func spaFallback(dist fs.FS) fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()
		rel := strings.TrimPrefix(path, "/")
		if rel == "" {
			rel = "index.html"
		}
		if _, err := fs.Stat(dist, rel); err == nil {
			return c.Next()
		}
		if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".mjs") ||
			strings.HasSuffix(path, ".map") || strings.HasSuffix(path, ".woff") || strings.HasSuffix(path, ".woff2") ||
			strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".ico") ||
			strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".webmanifest") {
			return c.Status(404).SendString("not found")
		}
		b, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			return c.Status(500).SendString("SPA index.html missing in embed")
		}
		c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
		return c.Send(b)
	}
}

func staticFromFS(fsys fs.FS) fiber.Handler {
	fileServer := http.FileServer(http.FS(fsys))
	return func(c fiber.Ctx) error {
		path := c.Path()
		rel := strings.TrimPrefix(path, "/")
		if rel == "" {
			return c.Next()
		}
		req, err := http.NewRequest(http.MethodGet, "/"+rel, nil)
		if err != nil {
			return err
		}
		rec := &recWriter{header: http.Header{}}
		fileServer.ServeHTTP(rec, req)
		if rec.status == 0 {
			rec.status = 200
		}
		for k, vs := range rec.header {
			for _, v := range vs {
				c.Set(k, v)
			}
		}
		c.Set(fiber.HeaderCacheControl, "no-cache")
		c.Status(rec.status)
		return c.Send(rec.body)
	}
}

type recWriter struct {
	header http.Header
	body   []byte
	status int
}

func (r *recWriter) Header() http.Header         { return r.header }
func (r *recWriter) Write(b []byte) (int, error) { r.body = append(r.body, b...); return len(b), nil }
func (r *recWriter) WriteHeader(s int)           { r.status = s }
