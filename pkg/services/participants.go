package services

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/rs/zerolog/log"
)

type (
	ParticipantsService struct {
		participantsDAO *db.ParticipantsDAO
	}
)

func NewParticipantsService() *ParticipantsService {
	return &ParticipantsService{}
}

func (ps *ParticipantsService) ProcessParticipantCSV(ctx context.Context, filePath string) {
	log.Info().Str("raceID", util.GetRaceIDFromContext(ctx)).Msg("Starting processing")
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to open CSV File: %s", filePath)
		return
	}

	reader := csv.NewReader(bufio.NewReader(file))

	participants := []db.Participant{}
	lines := 0
	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error().Err(err).Msg("Failed to read line from file")
		}
		lines++
		if lines == 1 {
			continue
		}
		fmt.Printf("%+v", line)
		name := strings.Split(line[1], " ")
		if len(name) < 2 {
			continue
		}
		birthYear, _ := strconv.Atoi(line[3])
		participants = append(participants, db.Participant{
			RaceID:    util.GetRaceIDFromContext(ctx),
			FirstName: name[0],
			LastName:  name[1],
			Team:      line[2],
			BirthYear: birthYear,
		})
	}
	ps.participantsDAO.BatchInsertParticipants(ctx, participants)
}
