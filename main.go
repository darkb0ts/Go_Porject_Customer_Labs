package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestData struct {
	Ev     string `json:"ev"`
	Et     string `json:"et"`
	Id     string `json:"id"`
	Uid    string `json:"uid"`
	Mid    string `json:"mid"`
	T      string `json:"t"`
	P      string `json:"p"`
	L      string `json:"l"`
	Sc     string `json:"sc"`
	Atrk1  string `json:"atrk1"`
	Atrv1  string `json:"atrv1"`
	Atrt1  string `json:"atrt1"`
	Atrk2  string `json:"atrk2"`
	Atrv2  string `json:"atrv2"`
	Atrt2  string `json:"atrt2"`
	Uatrk1 string `json:"uatrk1"`
	Uatrv1 string `json:"uatrv1"`
	Uatrt1 string `json:"uatrt1"`
	Uatrk2 string `json:"uatrk2"`
	Uatrv2 string `json:"uatrv2"`
	Uatrt2 string `json:"uatrt2"`
	Uatrk3 string `json:"uatrk3"`
	Uatrv3 string `json:"uatrv3"`
	Uatrt3 string `json:"uatrt3"`
}

type AltRequestData struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	Traits          map[string]Trait     `json:"traits"`
}

type TransformedData struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	Traits          map[string]Trait     `json:"traits"`
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Trait struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func main() {
	dataChannel := make(chan interface{})

	go worker(dataChannel)

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Println("Received request body:", string(body))

		var altRequestData AltRequestData
		var requestData RequestData

		err = json.Unmarshal(body, &altRequestData)
		if err != nil {
			log.Println("Failed to unmarshal AltRequestData:", err)
			err = json.Unmarshal(body, &requestData)
			if err != nil {
				log.Println("Failed to unmarshal RequestData:", err)
				http.Error(w, "Invalid JSON data", http.StatusBadRequest)
				return
			}
			dataChannel <- requestData
		} else {
			dataChannel <- altRequestData
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Received data successfully")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func worker(dataChannel <-chan interface{}) {
	for data := range dataChannel {
		successResponse := map[string]string{"status": "success"}
		successJSON, err := json.Marshal(successResponse)
		if err != nil {
			log.Printf("Error marshalling success response JSON: %v", err)
			continue
		}

		webhookURL := "https://webhook.site/6d64707f-35de-41af-bcf2-021a35d9f99a"
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(successJSON))
		if err != nil {
			log.Printf("Error sending success response to webhook: %v", err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("First webhook returned non-200 status: %v", resp.Status)
			continue
		}
		log.Println("Success response sent to first webhook")

		var transformedData TransformedData

		switch v := data.(type) {
		case RequestData:
			transformedData = TransformedData{
				Event:           v.Ev,
				EventType:       v.Et,
				AppID:           v.Id,
				UserID:          v.Uid,
				MessageID:       v.Mid,
				PageTitle:       v.T,
				PageURL:         v.P,
				BrowserLanguage: v.L,
				ScreenSize:      v.Sc,
				Attributes: map[string]Attribute{
					v.Atrk1: {Value: v.Atrv1, Type: v.Atrt1},
					v.Atrk2: {Value: v.Atrv2, Type: v.Atrt2},
				},
				Traits: map[string]Trait{
					v.Uatrk1: {Value: v.Uatrv1, Type: v.Uatrt1},
					v.Uatrk2: {Value: v.Uatrv2, Type: v.Uatrt2},
					v.Uatrk3: {Value: v.Uatrv3, Type: v.Uatrt3},
				},
			}
		case AltRequestData:
			transformedData = TransformedData{
				Event:           v.Event,
				EventType:       v.EventType,
				AppID:           v.AppID,
				UserID:          v.UserID,
				MessageID:       v.MessageID,
				PageTitle:       v.PageTitle,
				PageURL:         v.PageURL,
				BrowserLanguage: v.BrowserLanguage,
				ScreenSize:      v.ScreenSize,
				Attributes:      v.Attributes,
				Traits:          v.Traits,
			}
		default:
			log.Printf("Unknown data type received")
			continue
		}

		jsonData, err := json.Marshal(transformedData)
		if err != nil {
			log.Printf("Error marshalling transformedData JSON: %v", err)
			continue
		}
		finalWebhookURL := "https://webhook.site/6d64707f-35de-41af-bcf2-021a35d9f99a"
		resp, err = http.Post(finalWebhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error sending data to final webhook: %v", err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Final webhook returned non-200 status: %v", resp.Status)
		} else {
			log.Println("Data sent successfully to final webhook")
		}
	}
}
