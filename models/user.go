package models

import "fmt"

type User struct {
	UserID      string `field:"user_id"`
	DisplayName string `field:"display_name"`
	SyncCode    string `field:"sync_code"`
	PetCount    int    `field:"pet_count"`
}

func (u *User) PetDaisy() {
	u.PetCount++
}

func (u *User) GetHashKey() string {
	return fmt.Sprintf("users:%s", u.UserID)
}
