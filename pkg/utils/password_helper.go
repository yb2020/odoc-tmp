package utils

import (
	"strings"
)

var salt = "The#SEA%pwd&salt"

// StrengthenPassword is a static version of StrengthenPwd
// Can be used directly without creating a PasswordHelper instance
func StrengthenPassword(elements ...string) string {
	// Join all elements into a single string
	text := strings.Join(elements, "") + salt
	return MD5(text)
}
