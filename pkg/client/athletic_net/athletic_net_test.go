package athletic_net

import (
	"testing"
)

func Test_getEventIDAndTypeFromURL(t *testing.T) {
	type args struct {
		infoPageURL string
	}
	tests := []struct {
		name          string
		args          args
		wantMeetID    string
		wantEventType string
	}{
		{
			name: "border_wars_2023",
			args: args{
				infoPageURL: "https://www.athletic.net/CrossCountry/meet/223936/info",
			},
			wantMeetID:    "223936",
			wantEventType: "xc",
		},
		{
			name: "bob_firman_2023",
			args: args{
				infoPageURL: "https://www.athletic.net/CrossCountry/meet/221697/info",
			},
			wantMeetID:    "221697",
			wantEventType: "xc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeetID, gotEventType := getEventIDAndTypeFromURL(tt.args.infoPageURL)
			if gotMeetID != tt.wantMeetID {
				t.Errorf("getEventIDAndTypeFromURL() gotMeetID = %v, want %v", gotMeetID, tt.wantMeetID)
			}
			if gotEventType != tt.wantEventType {
				t.Errorf("getEventIDAndTypeFromURL() gotEventType = %v, want %v", gotEventType, tt.wantEventType)
			}
		})
	}
}

// func Test_getRawEventData(t *testing.T) {
// 	type args struct {
// 		meetID   string
// 		meetType string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "bob_firman_2023",
// 			args: args{
// 				meetID:   "221697",
// 				meetType: "xc",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := getRawEventData(tt.args.meetID, tt.args.meetType)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getRawEventData() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			// if !reflect.DeepEqual(gotData, tt.wantData) {
// 			// 	t.Errorf("getRawEventData() = %v, want %v", gotData, tt.wantData)
// 			// }
// 		})
// 	}
// }

// func Test_processRawEventData(t *testing.T) {
// 	type args struct {
// 		filename string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []model.Event
// 	}{
// 		{
// 			name: "border_wars_2023",
// 			args: args{
// 				filename: "testdata/bob_firman_2023.json",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			data, _ := os.ReadFile(tt.args.filename)
// 			var testdata []ANetMeetResultsData = []ANetMeetResultsData{}
// 			json.Unmarshal(data, &testdata)

// 			if got := processRawEventData(testdata); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("processRawEventData() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
