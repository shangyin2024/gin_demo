package wire

import (
	"gin_demo/internal/app"

	"github.com/google/wire"
)

// AppSet App 层的 Wire 集合
var AppSet = wire.NewSet(
	app.NewHandlers,
	app.New,
)
