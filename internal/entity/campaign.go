package entity

import "github.com/jinzhu/gorm"

// configured sizes are for mysql, since version 5 mysql counts characters, not bytes

type Campaign struct {
	gorm.Model
	Subject string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL;unique_index:subject_uidx"`
	Body string `gorm:"type:text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	Recipients []Recipient
}

type Recipient struct {
	gorm.Model
	ToAddress string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL;index:to_adr_idx"`
}
