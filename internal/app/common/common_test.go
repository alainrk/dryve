package common

import (
	"reflect"
	"testing"
	"time"
)

func TestParseAndValidateDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "Valid date",
			args: args{
				date: "2021-01-01",
			},
			want:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Valid date",
			args: args{
				date: "2021-12-31",
			},
			want:    time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Valid date",
			args: args{
				date: "2056-01-01",
			},
			want:    time.Date(2056, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Invalid date",
			args: args{
				date: "2021-13-01",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid date",
			args: args{
				date: "2021-00-01",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid date",
			args: args{
				date: "0000-13-43",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid date",
			args: args{

				date: "2021-01-01T00:00:00Z",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid date",
			args: args{
				date: "2021-01-01T00:00:00+00:00",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndValidateDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAndValidateDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAndValidateDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
