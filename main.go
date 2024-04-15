package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

type mediaRequest struct {
	Magnet string `json:"magnet"`
	Media  string `json:"media"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}
	defer r.Body.Close()

	var mediaRequest mediaRequest
	err = json.Unmarshal(body, &mediaRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing JSON: %v", err)
		return
	}

	fmt.Println("Received magnet:", mediaRequest.Magnet)

	fmt.Println("Received media:", mediaRequest.Media)

	// Start to torrent from the link given which will be saved to the media type parameter
	_, err = exec.Command("qbittorrent-nox", "--save-path=/mnt/Media/"+mediaRequest.Media, mediaRequest.Magnet).Output()
	if err != nil {
		fmt.Fprintf(w, "Error running qbittorrent: %v", err)
	}
	w.WriteHeader(http.StatusContinue)
	fmt.Fprintf(w, "Torrenting started, check qbittorrent UI at 192.168.1.5")

}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server listening on port :7123")
	http.ListenAndServe(":7123", nil)
}
