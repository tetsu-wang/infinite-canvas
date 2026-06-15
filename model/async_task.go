package model

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type TaskType string

const (
	TaskTypeImageGeneration TaskType = "image_generation"
	TaskTypeImageEdit       TaskType = "image_edit"
	TaskTypeChatCompletion  TaskType = "chat_completion"
	TaskTypeAudioSpeech     TaskType = "audio_speech"
	TaskTypeVideo           TaskType = "video"
)

// AsyncTask 异步任务记录。
type AsyncTask struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	UserID      string     `json:"userId" gorm:"index"`
	TaskType    TaskType   `json:"taskType"`
	ModelName   string     `json:"modelName"`
	Status      TaskStatus `json:"status" gorm:"index"`
	RequestBody string     `json:"requestBody,omitempty" gorm:"type:text"`
	Result      string     `json:"result,omitempty" gorm:"type:text"`
	Error       string     `json:"error,omitempty"`
	Progress    int        `json:"progress"` // 0-100
	CreatedAt   string     `json:"createdAt" gorm:"index"`
	UpdatedAt   string     `json:"updatedAt"`
}

// AsyncTaskList 异步任务分页结果。
type AsyncTaskList struct {
	Items []AsyncTask `json:"items"`
	Total int         `json:"total"`
}
