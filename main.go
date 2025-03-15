package main

import (
	"github.com/jiabaow/ReceiptProcessor/handlers"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	http.HandleFunc("/receipts/process", handlers.ProcessReceipt)
	http.HandleFunc("/receipts", handlers.GetPoints)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>

}
