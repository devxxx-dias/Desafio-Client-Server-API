package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type QuotationResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatal("timeout: server request exceeded 300ms")
		}
		log.Fatalf("error fetching quotation: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}

	var quotation QuotationResponse
	if err := json.Unmarshal(body, &quotation); err != nil {
		log.Fatalf("error parsing response: %v", err)
	}

	file, err := os.OpenFile("cotacao.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening cotacao.txt: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dólar: %s\n", quotation.Bid)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}

	log.Printf("cotacao.txt updated — Dólar: %s", quotation.Bid)
}
