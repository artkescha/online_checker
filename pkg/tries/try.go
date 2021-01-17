package try

import "time"

type Try struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Solution    string    `json:"solution"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	TaskID      int       `json:"task_id"`
	LanguageID  int       `json:"language_id"`
}
