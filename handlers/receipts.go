package handlers

import (
	"encoding/json"
	_ "encoding/json"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jiabaow/ReceiptProcessor/models"
	"log"
	"net/http"
	_ "net/http"
	"regexp"
	_ "regexp"
	"strconv"
	_ "strconv"
	"strings"
	_ "strings"
	"sync"
	"time"
	_ "time"
)

var (
	receipts = make(map[string]models.Receipt)
	mu       sync.Mutex
)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding receipt: %v", err)
		return
	}

	id := uuid.New().String()

	mu.Lock()
	receipts[id] = receipt
	mu.Unlock()
	log.Printf("Stored receipt with ID: %s", id)

	response := models.ProcessReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
}

func calculatePoints(r models.Receipt) int {
	points := 0

	// 1. one point for every alphanumeric char in the retailer name
	alnum := regexp.MustCompile("[a-zA-Z0-9]")
	alnumChars := alnum.FindAllString(r.Retailer, -1)
	points += len(alnumChars)
	log.Printf("1. Add %d points for retailer name: %s", len(alnumChars), r.Retailer)

	// 2. 50 points if the total is a round dollar amount with no cents
	total, _ := strconv.ParseFloat(r.Total, 64)
	if total == float64(int(total)) {
		points += 50
		log.Printf("2. Add 50 points for total: %s", r.Total)
	}

	// 3. 25 points if the total is a multiple of 0.25
	if total*100 == float64(int(total*100)) {
		if int(total*100)%25 == 0 {
			points += 25
			log.Printf("3. Add 25 points for total: %s", r.Total)
		}
	}

	// 4. 5 points for every two items on the receipt
	addition := (len(r.Items) / 2) * 5
	points += addition
	log.Printf("4. Add %d points for two items on the receipt", addition)

	// 5. If the trimmed length of the item description is a multiple of 3
	// multiply the price by 0.2 and round up to the nearest integer
	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			calculatedPrice := int((price * 0.2) + 0.999)
			points += calculatedPrice
			log.Printf("5. Add %d points for item description", calculatedPrice)
		}
	}

	// 6. If the total is greater than 10.00, add 5 points.
	// no need to add coz not generated using a llm
	//if total > 10.00 {
	//	points += 5
	//	log.Printf("6. Add 5 points for total > 10")
	//}

	// 7. 6 points if the day in the purchase date is odd
	date, _ := time.Parse("2006-01-02", r.PurchaseDate)
	if date.Day()%2 == 1 {
		points += 6
		log.Printf("7. Add 6 points for odd day in purchase date")
	}

	//8. 10 points if the time of purchase is after 2pm and before 4pm
	t, _ := time.Parse("15:04", r.PurchaseTime)
	if t.Hour() >= 14 && t.Hour() <= 16 {
		log.Printf("Time of purchase: %d", t.Hour())
		points += 10
		log.Printf("8. Add 10 points for time of purchase")
	}

	return points
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	receipt, exists := receipts[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		log.Printf("Receipt not found for ID: %s", id)
		return
	}

	points := calculatePoints(receipt)
	log.Printf("Calculated %d points for receipt ID: %s", points, id)

	response := models.PointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
