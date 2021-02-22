package status

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Status uint32

const (
	Undefined Status = iota
	Queue
	Processing
	Success
	Failure
)

func (s Status) status() (Status, error) {
	if s >= Undefined && s <= Failure {
		return s, nil
	}
	return 0, fmt.Errorf("status undefined")
}

func (s Status) String() string {
	switch s {
	case Undefined:
		return "undefined"
	case Queue:
		return "queue"
	case Processing:
		return "processing"
	case Success:
		return "success"
	case Failure:
		return "failure"
	default:
		return fmt.Sprintf("status %d undefined", s)
	}
}

func convertToStatus(s string) (Status, error) {
	switch s {
	case "undefined":
		return Undefined, nil
	case "queue":
		return Queue, nil
	case "processing":
		return Processing, nil
	case "success":
		return Success, nil
	case "failure":
		return Failure, nil
	default:
		return 0, fmt.Errorf("wrong status")
	}
}

// MarshalJSON marshals the enum as a quoted json string
func (s Status) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	status, err := s.status()
	if err != nil {
		return []byte{}, fmt.Errorf("status %s %s", s, err)
	}
	buffer.WriteString(status.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Status) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	status, err := convertToStatus(string(j))
	if err != nil {
		return fmt.Errorf("status %s  %s", j, err)
	}
	*s = status
	return nil
}
