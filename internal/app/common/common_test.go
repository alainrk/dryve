package common

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestDecodeJSONBody(t *testing.T) {
	testCases := []struct {
		name    string
		payload string
		dst     interface{}
		wantErr bool
	}{
		{
			name:    "valid payload",
			payload: `{"name": "Alice", "age": 30}`,
			dst: &struct {
				Name string
				Age  int
			}{},
			wantErr: false,
		},
		{
			name:    "empty payload",
			payload: ``,
			dst: &struct {
				Name string
				Age  int
			}{},
			wantErr: true,
		},
		{
			name:    "payload with unknown field",
			payload: `{"name": "Alice", "age": 30, "foo": "bar"}`,
			dst: &struct {
				Name string
				Age  int
			}{},
			wantErr: true,
		},
		{
			name:    "payload with invalid type",
			payload: `{"name": "Alice", "age": "30"}`,
			dst: &struct {
				Name string
				Age  int
			}{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			w := httptest.NewRecorder()

			err = DecodeJSONBody(w, req, tc.dst)

			if tc.wantErr {
				if err == nil {
					t.Errorf("%s: expected an error but got nil", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestEncodeJSONAndSend(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		res any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid response",
			args: args{
				w: httptest.NewRecorder(),
				res: struct {
					Name string
					Age  int
				}{
					Name: "Alice",
					Age:  30,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EncodeJSONAndSend(tt.args.w, tt.args.res)

			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
			}
		})
	}
}
