package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ChatRequest struct {
	Prompt string `json:"prompt"`
}

func handleChat(w http.ResponseWriter, r *http.Request) {

	// ✅ Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// ✅ Ambil prompt dari query
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	prompt := req.Prompt

	// Kirim ke TinyLLaMA
	payload := fmt.Sprintf(`{"model":"tinyllama", "prompt":%q, "stream":true}`, prompt)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", strings.NewReader(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	notify := r.Context().Done()

	// Streaming respons
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-notify:
			return // stop streaming kalau client disconnect
		default:
		}

		line := scanner.Text()
		fmt.Println(line)
		if strings.HasPrefix(line, "{") {
			var data struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}
			if err := json.Unmarshal([]byte(line), &data); err == nil {
				if data.Response != "" {
					fmt.Fprintf(w, "data: %s\n\n", data.Response)
					flusher.Flush()
				}
				if data.Done {
					fmt.Fprint(w, "data: [DONE]\n\n")
					flusher.Flush()
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading stream:", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", handleChat)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	handler := enableCORS(mux)

	fmt.Println("✅ Streaming on http://localhost:8080/chat")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
