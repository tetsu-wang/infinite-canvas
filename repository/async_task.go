package repository

import (
	"time"

	"github.com/basketikun/infinite-canvas/model"
	"github.com/google/uuid"
)

// CreateAsyncTask 创建异步任务。
func CreateAsyncTask(userID string, taskType model.TaskType, modelName string, requestBody string) (*model.AsyncTask, error) {
	db, err := DB()
	if err != nil {
		return nil, err
	}

	now := time.Now().Format(time.RFC3339)
	task := &model.AsyncTask{
		ID:          uuid.NewString(),
		UserID:      userID,
		TaskType:    taskType,
		ModelName:   modelName,
		Status:      model.TaskStatusPending,
		RequestBody: requestBody,
		Progress:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = db.Create(task).Error
	return task, err
}

// GetAsyncTask 根据 ID 查询异步任务。
func GetAsyncTask(id string) (*model.AsyncTask, error) {
	db, err := DB()
	if err != nil {
		return nil, err
	}

	var task model.AsyncTask
	err = db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateAsyncTaskStatus 更新任务状态和进度。
func UpdateAsyncTaskStatus(id string, status model.TaskStatus, progress int) error {
	db, err := DB()
	if err != nil {
		return err
	}

	return db.Model(&model.AsyncTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"progress":   progress,
		"updated_at": time.Now().Format(time.RFC3339),
	}).Error
}

// UpdateAsyncTaskResult 更新任务结果（成功）。
func UpdateAsyncTaskResult(id string, result string) error {
	db, err := DB()
	if err != nil {
		return err
	}

	return db.Model(&model.AsyncTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     model.TaskStatusCompleted,
		"result":     result,
		"progress":   100,
		"updated_at": time.Now().Format(time.RFC3339),
	}).Error
}

// UpdateAsyncTaskError 更新任务错误（失败）。
func UpdateAsyncTaskError(id string, errorMsg string) error {
	db, err := DB()
	if err != nil {
		return err
	}

	return db.Model(&model.AsyncTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     model.TaskStatusFailed,
		"error":      errorMsg,
		"updated_at": time.Now().Format(time.RFC3339),
	}).Error
}

// DeleteOldAsyncTasks 删除 N 天前的已完成或失败任务。
func DeleteOldAsyncTasks(days int) error {
	db, err := DB()
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)
	return db.Where("status IN ? AND created_at < ?", []model.TaskStatus{
		model.TaskStatusCompleted,
		model.TaskStatusFailed,
	}, cutoff).Delete(&model.AsyncTask{}).Error
}

// ListUserAsyncTasks 查询用户的异步任务列表。
func ListUserAsyncTasks(userID string, q model.Query) ([]model.AsyncTask, int64, error) {
	db, err := DB()
	if err != nil {
		return nil, 0, err
	}

	q.Normalize()
	tx := db.Model(&model.AsyncTask{}).Where("user_id = ?", userID)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tasks []model.AsyncTask
	err = tx.Order("created_at desc").Offset(q.Offset()).Limit(q.PageSize).Find(&tasks).Error
	return tasks, total, err
}
