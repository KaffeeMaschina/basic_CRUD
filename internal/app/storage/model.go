package storage

import "time"

type Task struct {
	ID          int       `json:"id"`                    // Уникальный идентификатор
	Title       string    `json:"title"`                 // Заголовок задачи
	Description *string   `json:"description,omitempty"` // Описание (может быть null)
	Status      string    `json:"status"`                // Статус задачи (new, in_progress, done)
	CreatedAt   time.Time `json:"created_at"`            // Дата создания
	UpdatedAt   time.Time `json:"updated_at"`            // Дата обновления
}
