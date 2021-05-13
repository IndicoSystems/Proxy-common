package metadata

import (
	"testing"
	"time"

	"github.com/go-test/deep"
)

var (
	timeA = time.Now()
	timeB = timeA.Add(time.Second)
	timeC = timeA.Add(2 * time.Second)
)

func pint(n int64) *int64 {
	return &n
}
func pstring(s string) *string {
	return &s
}
func pevent(s EventKind) *EventKind {
	return &s
}

func TestRecordingEvents_WithEndings(t *testing.T) {
	tests := []struct {
		name    string
		re      RecordingEvents
		want    RecordingEvents
		wantErr bool
	}{
		{
			name: "should combine events logically",
			re: RecordingEvents{
				{
					Id:         "a",
					Kind:       RecordingEventStarted,
					Start:      timeA,
					StartFrame: 0,
				},
				{
					Id:         "b",
					Kind:       RecordingEventPaused,
					Start:      timeB,
					StartFrame: 60,
				},
				{
					Id:         "c",
					Kind:       RecordingEventUserEnded,
					Start:      timeC,
					StartFrame: 120,
				},
			},
			want: RecordingEvents{
				{
					Id:         "a",
					Kind:       RecordingEventStarted,
					Start:      timeA,
					StartFrame: 0,
					End:        &timeB,
					EndFrame:   pint(60),
					EndKind:    pevent(RecordingEventPaused),
					EndId:      pstring("b"),
				},
				{
					Id:         "c",
					Kind:       RecordingEventUserEnded,
					Start:      timeC,
					StartFrame: 120,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.re.WithEndings()
			if (err != nil) != tt.wantErr {
				t.Errorf("WithEndings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Err did not match wanted result: \n diff %+v, \nGot:    '%+v', \nWanted: '%+v'", diff, got, tt.want)
				for i, event := range got {
					t.Logf("%d Got ID %s", i, event.Id)
					t.Logf("%d Got Kind %s", i, event.Kind)
				}
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	if len(got) != len(tt.want) {
			//		t.Errorf("Length did not match: expectes %d got %d", len(tt.want), len(got))
			//	}
			//  for i, event := range tt.want {
			//
			//  }
			//	t.Errorf("WithEndings() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
