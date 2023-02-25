package services

import (
	_ "github.com/mattn/go-sqlite3"
)

// var (
// 	testDBUrl       = "../../results.db"
// 	raceDAO         = db.NewRaceDAO(sqlx.MustOpen("sqlite3", testDBUrl))
// 	participantsDAO = db.NewParticipantsDAO(*sqlx.MustOpen("sqlite3", testDBUrl))
// )

// func setup() context.Context {
// 	NewRaceService(raceDAO)
// 	raceID, _ := GetRaceServiceInstance().CreateRace(context.TODO(), "test race 2022")
// 	ctx := context.WithValue(context.Background(), util.RaceID, raceID)

// 	return ctx
// }

// func TestParticipantsService_ProcessParticipantCSV(t *testing.T) {
// 	type fields struct {
// 		participantsDAO *db.ParticipantsDAO
// 	}
// 	type args struct {
// 		ctx      context.Context
// 		filePath string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		{
// 			name: "basic",
// 			fields: fields{
// 				participantsDAO: participantsDAO,
// 			},
// 			args: args{
// 				filePath: "testdata/participants-test.csv",
// 				ctx:      setup(),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ps := &ParticipantsService{
// 				participantsDAO: tt.fields.participantsDAO,
// 			}
// 			ps.ProcessParticipantCSV(tt.args.ctx, tt.args.filePath)
// 		})
// 	}
// }
