package v1

import (
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/go-chi/chi"
)

type EventResultsResources struct{}

func (rs EventResultsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handleEventResults)

	return r
}

func handleEventResults(w http.ResponseWriter, r *http.Request) {
	services.GetEventResultsInstance().GetEventResults(r.Context())
}

// select p.gender, p.bib_number, p.birth_year, r.result from participants as p join results as r on r.bib_number = p.bib_number GROUP by p.birth_year, p.gender, p.bib_number;
// select p.gender, p.bib_number, p.birth_year, r.result from participants as p join results as r on r.bib_number = p.bib_number GROUP by p.birth_year, p.gender, p.bib_number ORDER by p.birth_year, p.gender, r.result;
// select p.gender, p.bib_number, p.birth_year, r.result from participants as p join results as r on r.bib_number = p.bib_number ORDER by p.birth_year, p.gender, r.result;

// select * from participants where birth_year = 2012 and gender = "F";

// select p.gender, p.bib_number, p.birth_year, r.result from participants as p outer join results as r on r.bib_number = p.bib_number;
// select p.gender, p.birth_year, p.bib_number, r.result from results as r join participants as p on r.bib_number = p.bib_number WHERE birth_year = 2012 GROUP BY p.birth_year, p.gender;
