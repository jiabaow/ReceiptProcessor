package handlers

import (
	"encoding/json"
	_ "encoding/json"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
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
	"time"
	_ "time"
)

var receipts = make(map[string]models.Receipt)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	receipts[id] = receipt
	log.Printf("Stored receipt with ID: %s", id)

	response := models.ProcessReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func calculatePoints(r models.Receipt) int {
	points := 0

	// 1. one point for every alphanumeric char in the retailer name
	alnum := regexp.MustCompile("[a-zA-Z0-9]")
	alnumChars := alnum.FindAllString(r.Retailer, -1)
	points += len(alnumChars)

	// 2. 50 points if the total is a round dollar amount with no cents
	total, _ := strconv.ParseFloat(r.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

	// 3. 25 points if the total is a multiple of 0.25
	if total*100 == float64(int(total*100)) {
		if int(total*100)%25 == 0 {
			points += 25
		}
	}

	// 4. 5 points for every two items on the receipt
	points += (len(r.Items) / 2) * 5

	// 5. If the trimmed length of the item description is a multiple of 3
	// multiply the price by 0.2 and round up to the nearest integer
	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			calculatedPrice := int((price * 0.2) + 0.999)
			points += calculatedPrice
		}
	}

	// 6. If the total is greater than 10.00, add 5 points.
	if total > 10.00 {
		points += 5
	}

	// 7. 6 points if the day in the purchase date is odd
	date, _ := time.Parse("2006-01-02", r.PurchaseDate)
	if date.Day()%2 == 1 {
		points += 6
	}

	//8. 10 points if the time of purchase is after 2pm and before 4pm
	t, _ := time.Parse("15:04", r.PurchaseTime)
	if t.Hour() >= 14 || t.Hour() <= 16 {
		points += 10
	}

	return points
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	receipt, exists := receipts[id]
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
