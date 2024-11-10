package router

import (
	"fmt"
	"log"
	"net/http"
	"verve/controller"
)

func Accept(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	switch request.Method {
	case http.MethodGet:
		id := request.URL.Query().Get("id")
		if id == "" {
			writer.WriteHeader(http.StatusBadRequest)
			_, err := writer.Write([]byte("failed"))
			if err != nil {
				log.Println("Failed to write:")
				return
			}
		}
		controller.AddCount(id)

		// Optional parameter
		endpoint := request.URL.Query().Get("endpoint")
		if endpoint != "" {
			var resp []byte
			method := request.URL.Query().Get("method")
			if method == "POST" {
				resp = controller.SendPOSTRequest(endpoint)
			} else {
				resp = controller.SendGETRequest(endpoint)
			}
			_, err := writer.Write(resp)
			if err != nil {
				log.Printf("Failed to write: %s\n", resp)
				return
			}
		}
	case http.MethodPost:
		// TODO: Future Scope
	}
}

func Print(writer http.ResponseWriter, request *http.Request) {
	count := request.URL.Query().Get("count")
	if count == "" {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write([]byte("The 'count' parameter must be a positive integer."))
		if err != nil {
			return
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(fmt.Sprint("Unique requests ID count:", count)))
	if err != nil {
		return
	}
}

func Init() {
	// Main Functionality Server
	server := http.NewServeMux()
	server.HandleFunc("/api/verve/accept", Accept)
	go func() {
		err := http.ListenAndServe("localhost:8080", server)
		if err != nil {
			return
		}
	}()

	// Another Server : To get count printed
	anotherServer := http.NewServeMux()
	anotherServer.HandleFunc("/print", Print)
	go func() {
		err := http.ListenAndServe("localhost:8000", anotherServer)
		if err != nil {
			return
		}
	}()
}
