package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var contentTypeJSON = "application/json"

type Location struct {
	ID          string
	Name        string
	Slug        string
	Description string
}

var locations = []Location{
	Location{ID: "1", Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
	Location{ID: "2", Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
	Location{ID: "3", Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and ride a T-Rex"},
	Location{ID: "4", Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
	Location{ID: "5", Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
	Location{ID: "6", Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
}

func ListLocationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentTypeJSON)
	respBody, _ := json.Marshal(locations)
	w.Write(respBody)
}

func AddLocationFeedback(w http.ResponseWriter, r *http.Request) {
	var location Location
	vars := mux.Vars(r)
	slug := vars["slug"]
	for _, l := range locations {
		if l.Slug == slug {
			location = l
		}
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	if location.Slug != "" {
		respBody, _ := json.Marshal(location)
		w.Write(respBody)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Location Not Found"))
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentTypeJSON)
	respData := map[string]string{"status": "ok"}
	respBody, _ := json.Marshal(respData)
	w.Write(respBody)
}

func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}
