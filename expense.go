package main

import (
	"math"
	"time"

	"github.com/jinzhu/gorm"
)

// Expense - Model for Expense type
type Expense struct {
	gorm.Model
	Name         string    `json:"name"`
	Interval     uint      `json:"interval"` // Days
	Amount       float64   `json:"amount"`
	StartDate    time.Time `json:"start_date"`
	Rollover     bool      `json:"rollover"`
	Consumptions []Consumption
}

// HoursSinceStartDate - Number of hours since the start date
func (e Expense) HoursSinceStartDate() int {
	d := time.Now().Sub(e.StartDate)
	return int(d.Hours())
}

// Rate - Amount / Hour
func (e Expense) Rate() float64 {
	return e.Amount / (float64(e.Interval) * 24.0)
}

// Accumulation - Amount available since StartDate
func (e Expense) Accumulation() float64 {
	hoursSinceStart := time.Now().Sub(e.StartDate).Hours()
	return math.Floor(hoursSinceStart * e.Rate())
}

// Consumed - Total amount consumed for this expense
func (e Expense) Consumed() float64 {
	var consumptions []Consumption
	db.Model(&e).Related(&consumptions)

	sum := 0.0
	for _, consumption := range consumptions {
		sum += consumption.Amount
	}
	return sum
}

// Remaining - Amount available for consumption
func (e Expense) Remaining() float64 {
	return e.Accumulation() - e.Consumed()
}
