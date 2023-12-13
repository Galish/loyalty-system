package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/loyalty-system/internal/loyalty"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/api/orders/{id}", getOrder)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	order := loyalty.Order{
		ID: chi.URLParam(r, "id"),
	}

	switch order.ID {
	case "277431151":
		order.Status = loyalty.StatusProcessing

	case "2774311589":
		order.Status = loyalty.StatusInvalid

	default:
		order.Status = loyalty.StatusProcessed
		order.Accrual = 100
	}

	payload, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Err:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
