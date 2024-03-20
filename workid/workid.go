package workid

import "context"

type Worker interface {
	Get(ctx context.Context) Conn
	SetAppName(appName string) // 设置应用名
	SetModName(modName string) // 设置模块名，如果一个应用不同的模块需要单独的workID
}

type Conn interface {
	GetWorkID(ctx context.Context) (int, error) // 获取workID
	CleanWorkID(ctx context.Context) error      // 清理workID
}
