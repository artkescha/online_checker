package try

import (
	"github.com/artkescha/checker/online_checker/pkg/tries/status"
	"time"
)

type Try struct {
	ID         int           `json:"id,string"`
	UserID     int64         `json:"user_id,string"`
	Solution   string        `json:"solution"`
	Status     status.Status `json:"status"`
	Created    time.Time     `json:"created"`
	TaskID     int           `json:"task_id,string"`
	LanguageID int           `json:"language_id,string"`
}
