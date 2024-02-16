package handler

import (
	"fmt"
	"io"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://zflljjcmhaqiqknshltg.supabase.co", nil)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InpmbGxqamNtaGFxaXFrbnNobHRnIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MDgxMDg4NTEsImV4cCI6MjAyMzY4NDg1MX0.eHcXNMQ1IF72xqYA-Cpe0B5L8xfpq19NSgdHPBsyLc0")
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	// get the response body as a string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, string(body))
	r.Response.StatusCode = resp.StatusCode
}
