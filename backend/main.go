package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/mukesh1352/splitwise-backend/ledger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}
	defer pool.Close()
	//
	//converting the pgxpool to *sql
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("database connection established successfully")

	l := ledger.New(sqlDB)
	mux := http.NewServeMux()

	// Get the balances of the users
	mux.HandleFunc("/balances/user", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		balances, err := l.GetUserBalances(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(balances)
	})

	// Get the balances for the groups
	mux.HandleFunc("/balances/groups", func(w http.ResponseWriter, r *http.Request) {
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id is required..", http.StatusBadRequest)
			return
		}
		balances, err := l.GetGroupBalances(groupID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(balances)
	})

	// create the expenses
	mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method is not allowed..", http.StatusMethodNotAllowed)
			return
		}
		var input ledger.ExpenseInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request body..", http.StatusBadRequest)
			return
		}
		if err := l.CreateExpense(r.Context(), input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Expense created successfully..",
		})
	})

	// settling the balance
	mux.HandleFunc("/settle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var request struct {
			FromUserID string  `json:"from_user_id"`
			ToUserID   string  `json:"to_user_id"`
			Amount     float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := l.SettleBalance(
			r.Context(),
			request.FromUserID,
			request.ToUserID,
			request.Amount,
		); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Settlement is recorded succesfully..",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on.. : " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
