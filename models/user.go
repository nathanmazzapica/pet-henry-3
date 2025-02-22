package models

type User struct {
	UserID      string `field:"user_id"`
	DisplayName string `field:"display_name"`
	SyncCode    string `field:"sync_code"`
	PetCount    int    `field:"pet_count"`
}
