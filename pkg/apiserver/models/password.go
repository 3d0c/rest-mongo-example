package models

import (
	"fmt"
	"net/http"
)

// Password is a meta object, for binding requests only
type Password struct {
	OldPassword string `bson:"old_password" json:"old_password,omitempty" `
	NewPassword string `bson:"new_password" json:"new_password,omitempty" `
}

// Bind interface
func (u *Password) Bind(r *http.Request) error {
	if u.OldPassword == "" {
		return fmt.Errorf("old password is required")
	}

	if u.OldPassword == "" {
		return fmt.Errorf("new password is required")
	}

	return nil
}
