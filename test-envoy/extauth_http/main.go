package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type errResponse struct {
	Type    *string `json:"type,omitempty"`
	Message *string `json:"message,omitempty"`
}

const requiredPermission = "permission6"

func handler(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, headerValues := range headers {
			if strings.ToLower(name) == "my-credential-header" {
				perms := strings.Split(headerValues, ",")
				for _, perm := range perms {
					if perm == requiredPermission {
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}
	}

	errType := "FORBIDDEN"
	msg := "you are forbidden :("
	err := &errResponse{
		Type:    &errType,
		Message: &msg,
	}

	responseBytes, _ := json.MarshalIndent(err, "", "  ")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(responseBytes))
}

func main() {
	fmt.Println("LISTENING on 10003...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":10003", nil)
}
