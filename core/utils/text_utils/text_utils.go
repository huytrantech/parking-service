package text_utils

import (
	"encoding/json"
	"regexp"
)

func ParseStringToStruct(body string , data interface{}){
	json.Unmarshal([]byte(body), &data)
}

func CheckRegexPhoneNumber(phone string) bool {

	if m, _ := regexp.MatchString("(84|0[3|5|7|8|9])+([0-9]{8})\\b", phone); !m {
		return false
	}

	return true
}