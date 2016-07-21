package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ExpenseListItem - Output formatting for Expense
type ExpenseListItem struct {
	Name      string  `json:"name"`
	ID        uint    `json:"id"`
	Interval  uint    `json:"interval"`
	Amount    float64 `json:"amount"`
	Rollover  bool    `json:"rollover"`
	Remaining float64 `json:"remaining"`
	Consumed  float64 `json:"consumed"`
	Rate      float64 `json:"rate"`
}

func expenseAPIOutput(expense Expense) ExpenseListItem {
	return ExpenseListItem{
		Name:      expense.Name,
		ID:        expense.ID,
		Interval:  expense.Interval,
		Amount:    expense.Amount,
		Rollover:  expense.Rollover,
		Remaining: expense.Remaining(),
		Consumed:  expense.Consumed(),
		Rate:      expense.Rate()}
}

func expenseIndexHandler(writer http.ResponseWriter, request *http.Request) {
	var expenses []Expense
	db.Find(&expenses)

	expenseListItems := make([]ExpenseListItem, len(expenses), len(expenses))
	for idx, expense := range expenses {
		expenseListItems[idx] = expenseAPIOutput(expense)
	}

	body, _ := json.Marshal(expenseListItems)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(body)
}

func expenseCreateHandler(writer http.ResponseWriter, request *http.Request) {
	if request.ParseForm() != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	startDate, _ := time.Parse(time.RFC3339, request.Form["start_date"][0])
	interval, _ := strconv.ParseUint(request.Form["interval"][0], 10, 64)
	amount, _ := strconv.ParseFloat(request.Form["amount"][0], 64)
	name := request.Form["name"][0]
	rollover := request.Form["rollover"][0] == "true"

	expense := Expense{StartDate: startDate, Interval: uint(interval), Name: name, Rollover: rollover, Amount: amount}
	db.Create(&expense)

	body, _ := json.Marshal(expenseAPIOutput(expense))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(body)
}

func expenseUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	if request.ParseForm() != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(request)
	expenseID, _ := strconv.ParseUint(vars["id"], 10, 64)

	var expense Expense
	db.Where("id = ?", expenseID).Find(&expense)
	if expense.ID == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	expense.StartDate, _ = time.Parse(time.RFC3339, request.Form["start_date"][0])
	interval, _ := strconv.ParseUint(request.Form["interval"][0], 10, 64)
	expense.Interval = uint(interval)
	expense.Amount, _ = strconv.ParseFloat(request.Form["amount"][0], 64)
	expense.Name = request.Form["name"][0]
	expense.Rollover = request.Form["rollover"][0] == "true"
	db.Save(&expense)

	body, _ := json.Marshal(expenseAPIOutput(expense))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(body)
}

func expenseDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	expenseID, _ := strconv.ParseUint(vars["id"], 10, 64)

	var expense Expense
	db.Where("id = ?", expenseID).Find(&expense)
	if expense.ID == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	db.Where("expense_id = ?", expense.ID).Delete(Consumption{})
	db.Delete(&expense)
}

func consumptionIndexHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	expenseID, _ := strconv.ParseUint(vars["expense_id"], 10, 64)

	var consumptions []Consumption
	db.Where("expense_id = ?", expenseID).Find(&consumptions)

	type ConsumptionListItem struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Amount    float64   `json:"amount"`
	}
	consumptionListItems := make([]ConsumptionListItem, len(consumptions), len(consumptions))
	for idx, consumption := range consumptions {
		consumptionListItems[idx] = ConsumptionListItem{
			ID:        consumption.ID,
			CreatedAt: consumption.CreatedAt,
			Amount:    consumption.Amount}
	}
	body, _ := json.Marshal(consumptionListItems)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(body)
}

func consumptionCreateHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	if request.ParseForm() != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	expenseID, _ := strconv.ParseUint(vars["expense_id"], 10, 64)
	amount, _ := strconv.ParseFloat(request.Form["amount"][0], 64)

	consumption := Consumption{ExpenseID: uint(expenseID), Amount: amount}
	db.Save(&consumption)

	body, _ := json.Marshal(consumption)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(body)
}

func consumptionDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	consumptionID, _ := strconv.ParseUint(vars["id"], 10, 64)

	var consumption Consumption
	db.Where("id = ?", consumptionID).Find(&consumption)
	if consumption.ID == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	db.Delete(&consumption)
}
