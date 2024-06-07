package remote

import (
	"io"

	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	io.Closer
	listenAddr string
	app        *fiber.App
	callbacks  []func(app *fiber.App)
}

func NewWebServer(listenAddr string) *WebServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	return &WebServer{
		listenAddr: listenAddr,
		app:        app,
	}
}

func (w *WebServer) AddCallback(callback func(app *fiber.App)) {
	w.callbacks = append(w.callbacks, callback)
}

func (w *WebServer) Listen() error {
	for _, callback := range w.callbacks {
		callback(w.app)
	}

	return w.app.Listen(w.listenAddr)
}

func (w *WebServer) Close() error {
	return w.app.Shutdown()
}
