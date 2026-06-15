package service

import (
	"log"
	"time"

	"github.com/basketikun/infinite-canvas/repository"
)

// StartAsyncTaskCleanupScheduler 启动异步任务清理定时器，每天清理一次 7 天前的已完成任务。
func StartAsyncTaskCleanupScheduler() {
	const cleanupInterval = 24 * time.Hour
	const retentionDays = 7

	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		// 启动时立即执行一次
		cleanupOldTasks(retentionDays)

		for range ticker.C {
			cleanupOldTasks(retentionDays)
		}
	}()
}

func cleanupOldTasks(days int) {
	if err := repository.DeleteOldAsyncTasks(days); err != nil {
		log.Printf("Failed to cleanup old async tasks: %v", err)
	} else {
		log.Printf("Cleaned up async tasks older than %d days", days)
	}
}
