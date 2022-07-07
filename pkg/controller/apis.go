package controller

import "github.com/go-chi/chi"

type APIsResource struct{}

func (rs APIsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Mount("/races", RacesResources{}.Routes())

	return r
}

// /races
// POST /races
// GET /races/{id}

// /races/{id}/participants
// POST /races/{id}/participants/csv
// GET /races/{id}/participants
// GET /races/{id}/participants/{id}/results

// /races/{id}/events
// POST /races/{id}/events
// GET /races/{id}/events
// GET /races/{id}/events/{id}

// /races/{id}/results
// POST /races/{id}/results

// /races/{id}/events/{id}/results
// POST /races/{id}/events/{id}/results

// race table - id, name, date,
// event table - id, name,
// participant table - id, name, age group, birth year,
// results table - id, event id, race id, participant id, place, heat?, result
//

// views
// Race List
// Race View
// Events List
// Per Event View
// Per Participant View
// Participant Linking??
