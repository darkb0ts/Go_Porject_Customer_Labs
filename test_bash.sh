#!/bin/bash

# cURL request 1
curl -X POST -H "Content-Type: application/json" -d '{
  "event": "top_cta_clicked",
  "event_type": "clicked",
  "app_id": "cl_app_id_001",
  "user_id": "cl_app_id_001-uid-001",
  "message_id": "cl_app_id_001-uid-001",
  "page_title": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
  "page_url": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
  "browser_language": "en-US",
  "screen_size": "1920 x 1080",
  "attributes": {
    "button_text": {
      "value": "Free trial",
      "type": "string"
    },
    "color_variation": {
      "value": "ESK0023",
      "type": "string"
    }
  },
  "traits": {
    "user_score": {
      "value": "1034",
      "type": "number"
    },
    "gender": {
      "value": "m",
      "type": "string"
    },
    "tracking_code": {
      "value": "POSERK093",
      "type": "string"
    }
  }
}' http://localhost:8080/submit

# cURL request 2
curl -X POST -H "Content-Type: application/json" -d '{
  "event": "contact_form_submitted",
  "event_type": "form_submit",
  "app_id": "cl_app_id_001",
  "user_id": "cl_app_id_001-uid-001",
  "message_id": "cl_app_id_001-uid-001",
  "page_title": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
  "page_url": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
  "browser_language": "en-US",
  "screen_size": "1920 x 1080",
  "attributes": {
    "form_varient": {
      "value": "red_top",
      "type": "string"
    },
    "ref": {
      "value": "XPOWJRICW993LKJD",
      "type": "string"
    }
  },
  "traits": {
    "name": {
      "value": "iron man",
      "type": "string"
    },
    "email": {
      "value": "ironman@avengers.com",
      "type": "string"
    },
    "age": {
      "value": "32",
      "type": "integer"
    }
  }
}' http://localhost:8080/submit

# cURL request 3
curl -X POST -H "Content-Type: application/json" -d '{
  "ev": "contact_form_submitted",
  "et": "form_submit",
  "id": "cl_app_id_001",
  "uid": "cl_app_id_001-uid-001",
  "mid": "cl_app_id_001-uid-001",
  "t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
  "p": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
  "l": "en-US",
  "sc": "1920 x 1080",
  "atrk1": "form_varient",
  "atrv1": "red_top",
  "atrt1": "string",
  "atrk2": "ref",
  "atrv2": "XPOWJRICW993LKJD",
  "atrt2": "string",
  "uatrk1": "name",
  "uatrv1": "iron man",
  "uatrt1": "string",
  "uatrk2": "email",
  "uatrv2": "ironman@avengers.com",
  "uatrt2": "string",
  "uatrk3": "age",
  "uatrv3": "32",
  "uatrt3": "integer"
}' http://localhost:8080/submit