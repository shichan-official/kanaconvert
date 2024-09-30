package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

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
	Debug    string `json:"debug"`
}

// Initialize tokenizer with IPA dictionary
var t *tokenizer.Tokenizer

func init() {
	// Create a new tokenizer instance
	t, _ = tokenizer.New(ipa.Dict())
}

// Convert Katakana to Hiragana
func katakanaToHiragana(text string) string {
	var result string
	for _, r := range text {
		if unicode.In(r, unicode.Katakana) {
			result += string(r - 0x60)
		} else {
			result += string(r)
		}
	}
	return result
}

// Convert Hiragana to Katakana
func hiraganaToKatakana(text string) string {
	var result string
	for _, r := range text {
		if unicode.In(r, unicode.Hiragana) {
			result += string(r + 0x60)
		} else {
			result += string(r)
		}
	}
	return result
}

// Convert Kanji and Hiragana/Katakana to Hiragana
func convertToHiragana(text string) string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		fmt.Println("Error initializing tokenizer:", err)
		return ""
	}

	tokens := t.Analyze(text, tokenizer.Search)
	var hiraganaResult string
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}
		// Get the reading (pronunciation) in Hiragana
		features := token.Features()
		if len(features) > 6 && features[6] != "*" {
			hiraganaResult += features[6] // Feature index 6 contains the Hiragana reading
		} else {
			hiraganaResult += token.Surface // If no reading, add the surface as is
		}
	}
	return hiraganaResult
}

// Convert Kanji and Hiragana/Katakana to Katakana
func convertToKatakana(text string) string {
	hiragana := convertToHiragana(text) // Get the Hiragana conversion first
	return hiraganaToKatakana(hiragana) // Convert Hiragana to Katakana
}

// Simple Romaji conversion (just a mock example)
func convertToRomanji(text string) string {
	// This is a simple, non-perfect mock-up for Romaji conversion.
	// For production use, consider using a robust library or external service.
	romanjiMap := map[string]string{
		"あ": "a", "い": "i", "う": "u", "え": "e", "お": "o",
		"か": "ka", "き": "ki", "く": "ku", "け": "ke", "こ": "ko",
		"さ": "sa", "し": "shi", "す": "su", "せ": "se", "そ": "so",
		"た": "ta", "ち": "chi", "つ": "tsu", "て": "te", "と": "to",
		"な": "na", "に": "ni", "ぬ": "nu", "ね": "ne", "の": "no",
		"は": "ha", "ひ": "hi", "ふ": "fu", "へ": "he", "ほ": "ho",
		"ま": "ma", "み": "mi", "む": "mu", "め": "me", "も": "mo",
		"や": "ya", "ゆ": "yu", "よ": "yo",
		"ら": "ra", "り": "ri", "る": "ru", "れ": "re", "ろ": "ro",
		"わ": "wa", "を": "wo",
		"ん": "n",
		// Combine with diacritics for voiced consonants
		"が": "ga", "ぎ": "gi", "ぐ": "gu", "げ": "ge", "ご": "go",
		"ざ": "za", "じ": "ji", "ず": "zu", "ぜ": "ze", "ぞ": "zo",
		"だ": "da", "ぢ": "ji", "づ": "zu", "で": "de", "ど": "do",
		"ば": "ba", "び": "bi", "ぶ": "bu", "べ": "be", "ぼ": "bo",
		"ぱ": "pa", "ぴ": "pi", "ぷ": "pu", "ぺ": "pe", "ぽ": "po",
		// Long vowels
		"ああ": "aa", "いい": "ii", "うう": "uu", "ええ": "ee", "おお": "oo",
		"かあ": "kaa", "きい": "kii", "くう": "kuu", "けえ": "kee", "こお": "koo",
		"さあ": "saa", "しい": "shii", "すう": "suu", "せえ": "see", "そお": "soo",
		"たあ": "taa", "ちい": "chii", "つう": "tsuu", "てえ": "tee", "とお": "too",
		"なあ": "naa", "にい": "nii", "ぬう": "nuu", "ねえ": "nee", "のお": "noo",
		"はあ": "haa", "ひい": "hii", "ふう": "fuu", "へえ": "hee", "ほお": "hoo",
		"まあ": "maa", "みい": "mii", "むう": "muu", "めえ": "mee", "もお": "moo",
		"やあ": "yaa", "ゆう": "yuu", "よお": "yoo",
		"らあ": "raa", "りい": "rii", "るう": "ruu", "れえ": "ree", "ろお": "roo",
		"わあ": "waa", "をお": "woo",
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
	romanji := convertToRomanji(hiragana)

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
