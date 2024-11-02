package ytdlp

import (
	"testing"
)

func TestExtract(t *testing.T) {
	type args struct {
		url        string
		formatCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "youtube", args: args{url: "https://www.youtube.com/watch?v=6-WNdFhbK5k", formatCode: "95"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Extract(tt.args.url, tt.args.formatCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Extract() = %v", got)
		})
	}
}
