package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

var resultsServiceInstance *ResultsService

type ResultsService struct {
	dao *db.ResultsDAO
}

func NewResultsService() {
	if resultsServiceInstance == nil {
		resultsServiceInstance = &ResultsService{
			dao: db.NewResultsDAO(),
		}
	}
}

func GetResultsServiceInstance() *ResultsService {
	if resultsServiceInstance == nil {
		resultsServiceInstance = &ResultsService{
			dao: db.NewResultsDAO(),
		}
	}
	return resultsServiceInstance
}

func (r *ResultsService) InsertResults(ctx context.Context, raceID, eventID, bibNumber, result string) error {
	return r.dao.InsertResult(ctx, raceID, eventID, bibNumber, convertToMilliseconds(result))
}

// select p.first_name, p.last_name, p.bib_number, r.result from participants as p join results as r on p.race_id = r.race_id and p.event_id = r.event_id and p.bib_number = r.bib_number
// select p.first_name, p.last_name, p.bib_number, p.birth_year, r.result from participants as p join results as r on p.race_id = r.race_id and p.event_id = r.event_id and p.bib_number = r.bib_number group by p.birth_year order by r.result;
