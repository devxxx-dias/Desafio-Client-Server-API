package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type AwesomeApiRequest struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type QuotationResponse struct {
	Bid string `json:"bid"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "src/database/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		cotacao TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", Quotation)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Quotation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("timeout: external API call exceeded 200ms")
			http.Error(w, "external API timeout", http.StatusGatewayTimeout)
		} else {
			log.Printf("error fetching quotation: %v", err)
			http.Error(w, "failed to fetch quotation", http.StatusInternalServerError)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	var apiResp AwesomeApiRequest
	if err := json.Unmarshal(body, &apiResp); err != nil {
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO cotacoes (cotacao) VALUES (?)")
	if err != nil {
		log.Printf("error preparing DB statement: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	_, err = stmt.ExecContext(dbCtx, apiResp.USDBRL.Bid)
	if err != nil {
		if dbCtx.Err() == context.DeadlineExceeded {
			log.Println("timeout: DB insert exceeded 10ms")
		} else {
			log.Printf("error inserting into DB: %v", err)
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	log.Printf("cotacao saved to DB: %s", apiResp.USDBRL.Bid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(QuotationResponse{Bid: apiResp.USDBRL.Bid})
}
