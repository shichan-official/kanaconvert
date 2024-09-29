package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Request body for text conversion
type ConvertRequest struct {
	Text string `json:"text"`
}

// Response body for the conversion
type ConvertResponse struct {
	Hiragana string `json:"hiragana"`
	Katakana string `json:"katakana"`
	Romanji  string `json:"romanji"`
}

// Helper function to convert tokens to Hiragana or Katakana
func convertToKana(text string, toKatakana bool) string {
	t := tokenizer.New()
	tokens := t.Analyze(text, tokenizer.Search)
	var result strings.Builder

	for _, token := range tokens {
		features := token.Features()
		if len(features) > 7 {
			if toKatakana {
				result.WriteString(features[7]) // Katakana
			} else {
				result.WriteString(features[6]) // Hiragana
			}
		} else {
			result.WriteString(token.Surface)
		}
	}
	return result.String()
}

// Simple Romaji conversion (just a mock example)
func convertToRomanji(text string) string {
	// This is a simple, non-perfect mock-up for Romaji conversion.
	// For production use, consider using a robust library or external service.
	romanjiMap := map[string]string{
		"あ": "a", "い": "i", "う": "u", "え": "e", "お": "o",
		"か": "ka", "き": "ki", "く": "ku", "け": "ke", "こ": "ko",
		// Extend with more mappings...
	}

	var result strings.Builder
	for _, r := range text {
		if romanji, ok := romanjiMap[string(r)]; ok {
			result.WriteString(romanji)
		} else {
			result.WriteString(string(r))
		}
	}
	return result.String()
}

// API handler for the Kanji conversion
func convertHandler(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hiragana := convertToKana(req.Text, false)
	katakana := convertToKana(req.Text, true)
	romanji := convertToRomanji(hiragana) // For simplicity, converting Hiragana to Romanji

	res := ConvertResponse{
		Hiragana: hiragana,
		Katakana: katakana,
		Romanji:  romanji,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	// Existing sum handler
	http.HandleFunc("/sum", enableCors(sumHandler))

	// New Kanji to Kana conversion endpoint
	http.HandleFunc("/convert", enableCors(convertHandler))

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
