package web

import (
	"context"
	"log/slog"
	"os"

	"github.com/wasilak/go-hello-world/web/chi"
	"github.com/wasilak/go-hello-world/web/common"
	"github.com/wasilak/go-hello-world/web/echo"
	"github.com/wasilak/go-hello-world/web/fiber"
	"github.com/wasilak/go-hello-world/web/gin"
	"github.com/wasilak/go-hello-world/web/gorilla"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RunWebServer(ctx context.Context, frameworkOptions common.FrameworkOptions) {
	caser := cases.Title(language.English)
	var server common.WebServerInterface
	var isRunning bool

	for {
		select {
		case <-ctx.Done():
			// Context cancellation received, stop the server and exit
			if isRunning && server != nil {
				slog.DebugContext(ctx, "Shutting down server before exiting")
				server.Stop(ctx)
			}
			return

		case webFramework := <-common.FrameworkChannel:
			// Stop the currently running server if one exists
			if isRunning && server != nil {
				slog.DebugContext(ctx, "Stopping server", "type", caser.String(webFramework))
				server.Stop(ctx)
				isRunning = false
			}

			// Initialize the selected framework
			slog.DebugContext(ctx, "Starting server", "type", caser.String(webFramework))
			slog.DebugContext(ctx, "Features supported", "loggergo", true, "statsviz", true, "tracing", true)

			ws := common.WebServer{
				Framework:        webFramework,
				FrameworkOptions: frameworkOptions,
			}

			switch webFramework {
			case "gorilla":
				server = &gorilla.Server{WebServer: &ws}
			case "echo":
				server = &echo.Server{WebServer: &ws}
			case "chi":
				server = &chi.Server{WebServer: &ws}
			case "gin":
				server = &gin.Server{WebServer: &ws}
			case "fiber":
				server = &fiber.Server{WebServer: &ws}
			default:
				slog.ErrorContext(ctx, "No valid web framework selected", "type", caser.String(webFramework))
				os.Exit(1)
			}

			// Start the new server
			if server != nil {
				server.Start(ctx)
				isRunning = true
			}
		}
	}
}
