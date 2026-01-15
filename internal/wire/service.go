package wire

import (
	"gin_demo/internal/domain/service"

	"github.com/google/wire"
)

// ServiceSet Service 层 Provider 集合
var ServiceSet = wire.NewSet(
	service.NewUserService,
	// 未来可以在这里添加其他 Service
	// service.NewArticleService,
	// service.NewCommentService,
)
