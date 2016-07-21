package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func main() {
	log.Println("Expense Tracker")
	dbName := os.Getenv("EXPENSE_TRACKER_DB_NAME")
	db = dbSetup(dbName)

	r := mux.NewRouter()
	r.HandleFunc("/expenses", expenseIndexHandler).Methods("GET")
	r.HandleFunc("/expenses", expenseCreateHandler).Methods("POST")
	r.HandleFunc("/expenses/{id:[0-9]+}", expenseDeleteHandler).Methods("DELETE")
	r.HandleFunc("/expenses/{id:[0-9]+}", expenseUpdateHandler).Methods("PUT")
	r.HandleFunc("/expenses/{expense_id:[0-9]+}/consumptions", consumptionIndexHandler).Methods("GET")
	r.HandleFunc("/expenses/{expense_id:[0-9]+}/consumptions", consumptionCreateHandler).Methods("POST")
	r.HandleFunc("/expenses/{expense_id:[0-9]+}/consumptions/{id:[0-9]+}", consumptionDeleteHandler).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%s", "4567"), nil)
}

func dbSetup(dbName string) *gorm.DB {
	dbHost := os.Getenv("EXPENSE_TRACKER_DB_HOST")
	dbUser := os.Getenv("EXPENSE_TRACKER_DB_USER")
	dbPass := os.Getenv("EXPENSE_TRACKER_DB_PASS")

	connectStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPass)
	db, err := gorm.Open("postgres", connectStr)

	if err != nil {
		log.Panic("failed to connect to database")
	}

	if !db.HasTable(&Expense{}) {
		db.CreateTable(&Expense{})
	}
	if !db.HasTable(&Consumption{}) {
		db.CreateTable(&Consumption{})
	}

	return db
}
