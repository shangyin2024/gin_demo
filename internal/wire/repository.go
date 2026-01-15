package wire

import (
	"gin_demo/internal/repository"

	"github.com/google/wire"
)

// RepositorySet Repository 层 Provider 集合
var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
	wire.Bind(new(repository.UserRepositoryInterface), new(*repository.UserRepository)),
	// 未来可以在这里添加其他 Repository
	// repository.NewArticleRepository,
	// repository.NewCommentRepository,
)
