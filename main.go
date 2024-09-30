package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode"

	"github.com/gojp/kana"
	"github.com/ikawaha/kagome-dict/ipa"
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

// Convert Katakana to Hiragana for mixed text
func katakanaToHiragana(text string) string {
	var result string
	for _, r := range text {
		// If the character is Katakana, convert it to Hiragana
		if unicode.In(r, unicode.Katakana) {
			result += string(r - 0x60) // Katakana to Hiragana shift
		} else {
			result += string(r) // Leave other characters unchanged
		}
	}
	return result
}

// Convert Hiragana to Katakana for mixed text
func hiraganaToKatakana(text string) string {
	var result string
	for _, r := range text {
		// If the character is Hiragana, convert it to Katakana
		if unicode.In(r, unicode.Hiragana) {
			result += string(r + 0x60) // Hiragana to Katakana shift
		} else {
			result += string(r) // Leave other characters unchanged
		}
	}
	return result
}

// Convert all text (including Kanji) to Hiragana
func convertToHiragana(text string) string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		fmt.Println("Error initializing tokenizer:", err)
		return ""
	}

	tokens := t.Analyze(text, tokenizer.Normal)
	var hiraganaResult string
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}
		features := token.Features()
		if len(features) > 7 && features[7] != "*" { // Kagome: Features[7] for phonetic reading (Hiragana)
			hiraganaResult += features[7] // Use the Hiragana reading
		} else if len(features) > 6 && features[6] != "*" { // Fallback to feature 6 if feature 7 is unavailable
			hiraganaResult += features[6]
		} else {
			hiraganaResult += token.Surface // If no reading, add the surface as is
		}
	}

	// Convert any Katakana in the result to Hiragana
	return katakanaToHiragana(hiraganaResult)
}

// Convert all text (including Kanji) to Katakana
func convertToKatakana(text string) string {
	hiragana := convertToHiragana(text) // First, convert everything to Hiragana
	return hiraganaToKatakana(hiragana) // Then, convert Hiragana to Katakana
}

// CORS Middleware to enable cross-origin requests
func enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins or specify frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	}
}

// API handler for the Kanji conversion
func convertHandler(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hiragana := convertToHiragana(req.Text)
	katakana := convertToKatakana(req.Text)
	romanji := kana.KanaToRomaji(hiragana)

	res := ConvertResponse{
		Hiragana: hiragana,
		Katakana: katakana,
		Romanji:  romanji,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	// Kanji to Kana conversion endpoint
	http.HandleFunc("/convert", enableCors(convertHandler))

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
