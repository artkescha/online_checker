package request

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
)

func DecodePostParams(value interface{}, r *http.Request) error {
	defer r.Body.Close()
	//p:= make([]byte, 10000)
	//r.Body.Read(p)
	//log.Printf("!!!!!!!!!!!!%v", string(p))
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return fmt.Errorf("decode request params failed: %s", err)
	}
	return nil
	_, err := govalidator.ValidateStruct(value)

	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			return allErrs
		}
	}
	return nil
}
