package handler

import (
	"log"
	"net/http"

	"github.com/basketikun/infinite-canvas/repository"
	"github.com/basketikun/infinite-canvas/service"
)

// GetAsyncTask 查询异步任务状态。
func GetAsyncTask(w http.ResponseWriter, r *http.Request, taskID string) {
	user, ok := service.UserFromContext(r.Context())
	if !ok {
		Fail(w, "未登录或权限不足")
		return
	}

	task, err := repository.GetAsyncTask(taskID)
	if err != nil {
		log.Printf("Failed to get async task: taskId=%s err=%v", taskID, err)
		Fail(w, "任务不存在")
		return
	}

	// 权限检查：只能查询自己的任务
	if task.UserID != user.ID {
		Fail(w, "无权访问此任务")
		return
	}

	// 优化：轮询时不返回完整的 result（可能很大），只返回状态和进度
	// 前端通过 URL 参数 includeResult=true 请求完整结果
	includeResult := r.URL.Query().Get("includeResult") == "true"
	if !includeResult && task.Status != "completed" {
		// 未完成的任务不需要返回 result
		task.Result = ""
	}

	OK(w, task)
}
