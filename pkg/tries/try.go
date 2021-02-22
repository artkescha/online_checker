package try

import (
	"github.com/artkescha/grader/online_checker/pkg/tries/status"
	"time"
)

type Try struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	Solution    string        `json:"solution"`
	Status      status.Status `json:"status"`
	Created     time.Time     `json:"created"`
	TaskID      int           `json:"task_id"`
	LanguageID  int           `json:"language_id"`
}
