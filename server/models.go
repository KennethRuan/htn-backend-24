package main

import (
	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID   `json:"id"`
	Name    string      `json:"name"`
	Email   string      `json:"email"`
	Company string      `json:"company"`
	Phone   string      `json:"phone"`
	Skills  []UserSkill `json:"skills"`
}

type UserSkill struct {
	Name   string `json:"name"`
	Rating int    `json:"rating"`
}

type SkillFrequency struct {
	Name      string `json:"name"`
	Frequency int    `json:"frequency"`
}

type UserUpdate struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Company string `json:"company,omitempty"`
	Phone   string `json:"phone,omitempty"`
}
