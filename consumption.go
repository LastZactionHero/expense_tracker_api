package main

import "github.com/jinzhu/gorm"

// Consumption - Model for a Consumption on an Expense
type Consumption struct {
	gorm.Model
	ExpenseID uint    `json:"expense_id"`
	Amount    float64 `json:"amount"`
}
