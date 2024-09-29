package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

// Initialize tokenizer with IPA dictionary
var t *tokenizer.Tokenizer

func init() {
	// Create a new tokenizer instance
	t, _ = tokenizer.New(ipa.Dict())
}

// Helper function to convert tokens to Hiragana or Katakana
func convertToKana(text string, toKatakana bool) string {
	// Analyze the text
	tokens := t.Analyze(text, tokenizer.Search) // Pass text as string
	var result strings.Builder

	for _, token := range tokens {
		features := token.Features()
		// Skip BOS and EOS tokens
		if token.Surface == "BOS" || token.Surface == "EOS" {
			continue
		}
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

// Convert Katakana to Hiragana
func convertToHiragana(text string) string {
	// Katakana to Hiragana mapping
	hiraganaMap := map[string]string{
		"ア": "あ", "イ": "い", "ウ": "う", "エ": "え", "オ": "お",
		"カ": "か", "キ": "き", "ク": "く", "ケ": "け", "コ": "こ",
		"サ": "さ", "シ": "し", "ス": "す", "セ": "せ", "ソ": "そ",
		"タ": "た", "チ": "ち", "ツ": "つ", "テ": "て", "ト": "と",
		"ナ": "な", "ニ": "に", "ヌ": "ぬ", "ネ": "ね", "ノ": "の",
		"ハ": "は", "ヒ": "ひ", "フ": "ふ", "ヘ": "へ", "ホ": "ほ",
		"マ": "ま", "ミ": "み", "ム": "む", "メ": "め", "モ": "も",
		"ヤ": "や", "ユ": "ゆ", "ヨ": "よ",
		"ラ": "ら", "リ": "り", "ル": "る", "レ": "れ", "ロ": "ろ",
		"ワ": "わ", "ヲ": "を",
		"ン": "ん",
		// Voiced consonants
		"ガ": "が", "ギ": "ぎ", "グ": "ぐ", "ゲ": "げ", "ゴ": "ご",
		"ザ": "ざ", "ジ": "じ", "ズ": "ず", "ゼ": "ぜ", "ゾ": "ぞ",
		"ダ": "だ", "ヂ": "ぢ", "ヅ": "づ", "デ": "で", "ド": "ど",
		"バ": "ば", "ビ": "び", "ブ": "ぶ", "ベ": "べ", "ボ": "ぼ",
		"パ": "ぱ", "ピ": "ぴ", "プ": "ぷ", "ペ": "ぺ", "ポ": "ぽ",
	}

	var result strings.Builder
	for _, r := range text {
		if hiragana, ok := hiraganaMap[string(r)]; ok {
			result.WriteString(hiragana)
		} else {
			result.WriteString(string(r)) // Keep the character as is if not found in the map
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

	katakana := convertToKana(req.Text, true)
	hiragana := convertToHiragana(katakana)
	romanji := convertToRomanji(hiragana)

	res := ConvertResponse{
		Hiragana: hiragana,
		Katakana: katakana,
		Romanji:  romanji,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Existing sum handler for reference
func sumHandler(w http.ResponseWriter, r *http.Request) {
	type SumRequest struct {
		Num1 float64 `json:"num1"`
		Num2 float64 `json:"num2"`
	}
	type SumResponse struct {
		Sum float64 `json:"sum"`
	}

	var req SumRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sum := req.Num1 + req.Num2
	res := SumResponse{Sum: sum}

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
