package model

type Event struct {
	Name     string   `json:"name"`
	Results  []Result `json:"results"`
	Distance int      `json:"distance"`
}

type Result struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Place     int    `json:"place"`
	Team      string `json:"team"`
	Year      string `json:"year"`
	Time      int    `json:"time"`
}
