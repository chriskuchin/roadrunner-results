package services

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/rs/zerolog/log"
)

type (
	ParticipantsService struct {
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

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error().Err(err).Msg("Failed to read line from file")
		}
		fmt.Printf("%+v", line)
	}
}
