package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// document request this format https://docs.google.com/document/d/1kavz33uxvHF1XOugxpSPTr1W0CyhNzaZpmgN-xkN1eU
type RequestData struct { 
    Ev    string `json:"ev"`
    Et    string `json:"et"`
    ID    string `json:"id"`
    Uid   string `json:"uid"`
    Mid   string `json:"mid"`
    T     string `json:"t"`
    P     string `json:"p"`
    L     string `json:"l"`
    Sc    string `json:"sc"`
    Atrk1 string `json:"atrk1,omitempty"`
    Atrv1 string `json:"atrv1,omitempty"`
    Atrt1 string `json:"atrt1,omitempty"`
    Atrk2 string `json:"atrk2,omitempty"`
    Atrv2 string `json:"atrv2,omitempty"`
    Atrt2 string `json:"atrt2,omitempty"`
    Atrk3 string `json:"atrk3,omitempty"`
    Atrv3 string `json:"atrv3,omitempty"`
    Atrt3 string `json:"atrt3,omitempty"`
    Atrk4 string `json:"atrk4,omitempty"`
    Atrv4 string `json:"atrv4,omitempty"`
    Atrt4 string `json:"atrt4,omitempty"`
    Uatrk1 string `json:"uatrk1,omitempty"`
    Uatrv1 string `json:"uatrv1,omitempty"`
    Uatrt1 string `json:"uatrt1,omitempty"`
    Uatrk2 string `json:"uatrk2,omitempty"`
    Uatrv2 string `json:"uatrv2,omitempty"`
    Uatrt2 string `json:"uatrt2,omitempty"`
    Uatrk3 string `json:"uatrk3,omitempty"`
    Uatrv3 string `json:"uatrv3,omitempty"`
    Uatrt3 string `json:"uatrt3,omitempty"`
    Uatrk4 string `json:"uatrk4,omitempty"`
    Uatrv4 string `json:"uatrv4,omitempty"`
    Uatrt4 string `json:"uatrt4,omitempty"`
    Uatrk5 string `json:"uatrk5,omitempty"`
    Uatrv5 string `json:"uatrv5,omitempty"`
    Uatrt5 string `json:"uatrt5,omitempty"`
    Uatrk6 string `json:"uatrk6,omitempty"`
    Uatrv6 string `json:"uatrv6,omitempty"`
    Uatrt6 string `json:"uatrt6,omitempty"`
}

// document request this format https://docs.google.com/document/d/1kavz33uxvHF1XOugxpSPTr1W0CyhNzaZpmgN-xkN1eU
type TransformedData struct { 
    Event           string                 `json:"event"`
    EventType       string                 `json:"event_type"`
    AppID           string                 `json:"app_id"`
    UserID          string                 `json:"user_id"`
    MessageID       string                 `json:"message_id"`
    PageTitle       string                 `json:"page_title"`
    PageURL         string                 `json:"page_url"`
    BrowserLanguage string                 `json:"browser_language"`
    ScreenSize      string                 `json:"screen_size"`
    Attributes      map[string]Attribute   `json:"attributes"`
    Traits          map[string]Trait       `json:"traits"`
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

    dataChannel := make(chan RequestData)

    go worker(dataChannel)

    http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        var requestData RequestData
        if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
            http.Error(w, "Invalid JSON data", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Received data successfully")
        dataChannel <- requestData
    })

    log.Fatal(http.ListenAndServe(":8080", nil)) 
}


func worker(dataChannel <-chan RequestData) {
    for data := range dataChannel {
        fmt.Print(data) 
        successResponse := map[string]string{"status": "success"}
        successJSON, err := json.Marshal(successResponse)
        if err != nil {
            log.Printf("Error marshalling success response JSON: %v", err)
            continue
        }

        webhookURL := "https://webhook.site/db1ab579-43d9-457e-8f53-6ab1fdd047d6"
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
        
        transformedData := TransformedData{ //change data or transformed data 
            Event:           data.Ev,
            EventType:       data.Et,
            AppID:           data.ID,
            UserID:          data.Uid,
            MessageID:       data.Mid,
            PageTitle:       data.T,
            PageURL:         data.P,
            BrowserLanguage: data.L,
            ScreenSize:      data.Sc,
            Attributes:      make(map[string]Attribute),
            Traits:          make(map[string]Trait),
        }

        addAttributes(&transformedData, data) //check it first after use it 
        addTraits(&transformedData, data) //check tratis annd use it 

        jsonData, err := json.Marshal(transformedData)
        if err != nil {
            log.Printf("Error marshalling transformedData JSON: %v", err)
            continue
        }
        finalWebhookURL := "https://webhook.site/db1ab579-43d9-457e-8f53-6ab1fdd047d6" //final call
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

func addAttributes(td *TransformedData, rd RequestData) { //check the attr and add it 
    if rd.Atrk1 != "" {
        td.Attributes["atrk1"] = Attribute{Value: rd.Atrv1, Type: rd.Atrt1}
    }
    if rd.Atrk2 != "" {
        td.Attributes["atrk2"] = Attribute{Value: rd.Atrv2, Type: rd.Atrt2}
    }
    if rd.Atrk3 != "" {
        td.Attributes["atrk3"] = Attribute{Value: rd.Atrv3, Type: rd.Atrt3}
    }
    if rd.Atrk4 != "" {
        td.Attributes["atrk4"] = Attribute{Value: rd.Atrv4, Type: rd.Atrt4}
    }
}

func addTraits(td *TransformedData, rd RequestData) { //check the add traits
    if rd.Uatrk1 != "" {
        td.Traits["uatrk1"] = Trait{Value: rd.Uatrv1, Type: rd.Uatrt1}
    }
    if rd.Uatrk2 != "" {
        td.Traits["uatrk2"] = Trait{Value: rd.Uatrv2, Type: rd.Uatrt2}
    }
    if rd.Uatrk3 != "" {
        td.Traits["uatrk3"] = Trait{Value: rd.Uatrv3, Type: rd.Uatrt3}
    }
    if rd.Uatrk4 != "" {
        td.Traits["uatrk4"] = Trait{Value: rd.Uatrv4, Type: rd.Uatrt4}
    }
    if rd.Uatrk5 != "" {
        td.Traits["uatrk5"] = Trait{Value: rd.Uatrv5, Type: rd.Uatrt5}
    }
    if rd.Uatrk6 != "" {
        td.Traits["uatrk6"] = Trait{Value: rd.Uatrv6, Type: rd.Uatrt6}
    }
}
