package models

import (
    "time"
)

type User struct {
    ID                           uint      `gorm:"primaryKey" json:"id"`
    Name                         string    `json:"name"`
    Email                        string    `gorm:"uniqueIndex" json:"email"`
    PasswordHash                 string    `json:"-"`
    Verified                     bool      `json:"verified"`
    VerificationToken            string    `gorm:"index" json:"-"`
    VerificationTokenExpiresAt   time.Time `json:"-"`
    CreatedAt                    time.Time `json:"createdAt"`
    UpdatedAt                    time.Time `json:"updatedAt"`
}

