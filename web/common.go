package web

import (
	"context"
	"log/slog"

	"github.com/wasilak/go-hello-world/web/chi"
	"github.com/wasilak/go-hello-world/web/common"
	"github.com/wasilak/go-hello-world/web/echo"
	"github.com/wasilak/go-hello-world/web/fiber"
	"github.com/wasilak/go-hello-world/web/gin"
	"github.com/wasilak/go-hello-world/web/gorilla"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RunWebServer(ctx context.Context, webFramework string, frameworkOptions common.FrameworkOptions) {
	caser := cases.Title(language.English)
	slog.DebugContext(ctx, "Starting server", "type", caser.String(webFramework))
	slog.DebugContext(ctx, "Features supported", "loggergo", true, "statsviz", true, "tracing", true)

	switch webFramework {
	case "gorilla":
		gorilla.Init(ctx, frameworkOptions)
	case "echo":
		echo.Init(ctx, frameworkOptions)
	case "chi":
		chi.Init(ctx, frameworkOptions)
	case "gin":
		gin.Init(ctx, frameworkOptions)
	case "fiber":
		fiber.Init(ctx, frameworkOptions)
	}
}
