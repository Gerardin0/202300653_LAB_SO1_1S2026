package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        string `json:"vm"`
	Carnet    string `json:"carnet"`
}

type ApiStatusResponse struct {
	Apiname    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	loc, err := time.LoadLocation("America/Guatemala")
	if err != nil {
		log.Fatal(err)
	}

	response := HealthResponse{
		Status:    "UP",
		Message:   "API2 is Ready",
		Timestamp: time.Now().In(loc).Format(time.RFC1123),
		VM:        "1",
		Carnet:    "202300653",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func api1Handler(w http.ResponseWriter, r *http.Request) {

	client := http.Client{
		Timeout: 3 * time.Second, //! Evitar problemas con que se trabe
	}
	host := getEnv("ENV_HOST", "localhost")
	url := "http://" + host + ":8080/health"
	resp, err := client.Get(url) //? Puerto 8080 para la API1

	if err != nil || resp.StatusCode != http.StatusOK {
		json.NewEncoder(w).Encode(ApiStatusResponse{
			Apiname:    "API1",
			Message:    "ERROR: The API1 located on the VM1 is not working",
			Connection: false,
			Carnet:     "202300653",
		})
		return
	}
	defer resp.Body.Close()

	var api1Response HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&api1Response)
	if err != nil || api1Response.Status != "UP" {
		json.NewEncoder(w).Encode(ApiStatusResponse{
			Apiname:    "API1",
			Message:    "ERROR: The API1 located on the VM1 is not UP",
			Connection: false,
			Carnet:     "202300653",
		})
		return
	}

	json.NewEncoder(w).Encode(ApiStatusResponse{
		Apiname:    "API1",
		Message:    "SUCCESS: The API1 located on the VM1 is working",
		Connection: true,
		Carnet:     "202300653",
	})

}

func api3Handler(w http.ResponseWriter, r *http.Request) {

	client := http.Client{
		Timeout: 3 * time.Second, //! Evitar problemas con que se trabe
	}

	host := getEnv("ENV_HOST", "localhost")
	url := "http://" + host + ":8083/health"
	resp, err := client.Get(url) //? Puerto 8080 para la API3 ya que no tiene ninguna interferencia entre puertos

	if err != nil || resp.StatusCode != http.StatusOK {
		json.NewEncoder(w).Encode(ApiStatusResponse{
			Apiname:    "API3",
			Message:    "ERROR: The API3 located on the VM2 is not working",
			Connection: false,
			Carnet:     "202300653",
		})
		return
	}
	defer resp.Body.Close()

	var api3Response HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&api3Response)
	if err != nil || api3Response.Status != "UP" {
		json.NewEncoder(w).Encode(ApiStatusResponse{
			Apiname:    "API3",
			Message:    "ERROR: The API3 located on the VM2 is not UP",
			Connection: false,
			Carnet:     "202300653",
		})
		return
	}

	json.NewEncoder(w).Encode(ApiStatusResponse{
		Apiname:    "API3",
		Message:    "SUCCESS: The API3 located on the VM2 is working",
		Connection: true,
		Carnet:     "202300653",
	})

}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api2/202300653/call-api1", api1Handler)
	http.HandleFunc("/api2/202300653/call-api3", api3Handler)

	port := ":8081"
	host := getEnv("ENV_HOST", "localhost")
	log.Println("API2 escuchando en http://" + host + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
