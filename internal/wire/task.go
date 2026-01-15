package wire

import (
	"database/sql"

	"gin_demo/internal/app"
	"gin_demo/internal/task"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// TaskSet Task 层的 Wire 集合
var TaskSet = wire.NewSet(
	provideTaskManager,
)

// provideTaskManager 提供任务管理器
func provideTaskManager(db *sql.DB, redis redis.UniversalClient) app.TaskManager {
	return task.NewManager(redis, db)
}
