package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

func TestExpenseHoursSinceStartDate(t *testing.T) {
	var tests = []struct {
		durationAgo float64
		expected    int
	}{
		{400.0, 400},
		{400.9, 400},
	}

	for _, testcase := range tests {
		startDate := time.Now().Add(time.Duration(testcase.durationAgo*-1) * time.Hour)
		e := Expense{StartDate: startDate}
		hoursAgo := e.HoursSinceStartDate()
		if testcase.expected != hoursAgo {
			t.Error("Expected", testcase.expected, "got", hoursAgo)
		}
	}
}

func TestExpenseRate(t *testing.T) {
	var tests = []struct {
		amount   float64
		interval uint
		expected float64
	}{
		{100.0, 30, 100.0 / (30 * 24)},
		{100.0, 1, 100.0 / (1 * 24)},
	}

	for _, testcase := range tests {
		e := Expense{Amount: testcase.amount, Interval: testcase.interval}
		rate := e.Rate()
		if testcase.expected != rate {
			t.Error("Expect", testcase.expected, "got", rate)
		}
	}
}

func TestExpenseAccumulation(t *testing.T) {
	var tests = []struct {
		hoursAgo float64
		amount   float64
		interval uint
		expected float64
	}{
		{24, 100, 1, 100},
		{1, 240, 1, 10},
	}

	for _, testcase := range tests {
		startDate := time.Now().Add(time.Duration(testcase.hoursAgo*-1) * time.Hour)
		e := Expense{Amount: testcase.amount, Interval: testcase.interval, StartDate: startDate}
		accumulation := e.Accumulation()
		if testcase.expected != accumulation {
			t.Error("Expect", testcase.expected, "got", accumulation)
		}
	}
}

func TestExpenseConsumed(t *testing.T) {
	db = setupTestDB()

	e := Expense{}
	db.Create(&e)

	consumptionAmounts := []float64{10, 20, 30, 40}
	for _, amount := range consumptionAmounts {
		c := Consumption{ExpenseID: e.ID, Amount: amount}
		db.Create(&c)
	}

	consumed := e.Consumed()
	expected := float64(10 + 20 + 30 + 40)
	if consumed != expected {
		t.Error("Expected", expected, "got", consumed)
	}
	destroyTestDB(db)
}

func TestExpenseRemaining(t *testing.T) {
	db = setupTestDB()

	var tests = []struct {
		hoursAgo float64
		interval uint
		amount   float64
		consumed float64
		expected float64
	}{
		{240.0, 1, 10.0, 80.0, 20.0},
	}

	for _, testcase := range tests {
		startDate := time.Now().Add(time.Duration(testcase.hoursAgo*-1) * time.Hour)
		e := Expense{Amount: testcase.amount, Interval: testcase.interval, StartDate: startDate}
		db.Create(&e)

		c := Consumption{ExpenseID: e.ID, Amount: testcase.consumed}
		db.Create(&c)

		remaining := e.Remaining()
		if testcase.expected != remaining {
			t.Error("Expect", testcase.expected, "got", remaining)
		}
	}

	destroyTestDB(db)
}

func setupTestDB() *gorm.DB {
	devDBName := os.Getenv("EXPENSE_TRACKER_DB_NAME")
	testDBName := fmt.Sprintf("%s_test", devDBName)
	return dbSetup(testDBName)
}

func destroyTestDB(db *gorm.DB) {
	db.DropTable(&Consumption{})
	db.DropTable(&Expense{})
}
