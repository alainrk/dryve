package service

import (
	"dryve/internal/datastruct"
	"dryve/internal/mocks"
	"dryve/internal/repository"
	"io"
	"mime/multipart"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

var file1 = datastruct.File{
	UUID:     "11111111-1111-1111-1111-111111111111",
	Name:     "test",
	Size:     100,
	Filename: "test.txt",
}

func Test_fileService_Get(t *testing.T) {
	dao := &mocks.DAO{}
	fq := &mocks.FileQuery{}
	dao.On("NewFileQuery").Return(fq)
	fq.On("Get", "1").Return(file1, nil)
	fq.On("Get", "2").Return(datastruct.File{}, gorm.ErrRecordNotFound)

	type fields struct {
		dao             repository.DAO
		fileStoragePath string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    datastruct.File
		wantErr bool
	}{
		{
			name: "existing file",
			fields: fields{
				dao: dao,
			},
			args: args{
				id: "1",
			},
			want:    file1,
			wantErr: false,
		},
		{
			name: "not existing file",
			fields: fields{
				dao: dao,
			},
			args: args{
				id: "2",
			},
			want:    datastruct.File{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileService{
				dao:             tt.fields.dao,
				fileStoragePath: tt.fields.fileStoragePath,
			}
			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileService_Upload(t *testing.T) {
	type fields struct {
		dao             repository.DAO
		fileStoragePath string
	}
	type args struct {
		file       multipart.File
		fileHeader *multipart.FileHeader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    datastruct.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileService{
				dao:             tt.fields.dao,
				fileStoragePath: tt.fields.fileStoragePath,
			}
			got, err := s.Upload(tt.args.file, tt.args.fileHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileService.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileService.Upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileService_Delete(t *testing.T) {
	type fields struct {
		dao             repository.DAO
		fileStoragePath string
	}
	type args struct {
		metaFile datastruct.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileService{
				dao:             tt.fields.dao,
				fileStoragePath: tt.fields.fileStoragePath,
			}
			if err := s.Delete(tt.args.metaFile); (err != nil) != tt.wantErr {
				t.Errorf("fileService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileService_LoadFile(t *testing.T) {
	type fields struct {
		dao             repository.DAO
		fileStoragePath string
	}
	type args struct {
		metaFile datastruct.File
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantFile io.ReadCloser
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileService{
				dao:             tt.fields.dao,
				fileStoragePath: tt.fields.fileStoragePath,
			}
			gotFile, err := s.LoadFile(tt.args.metaFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileService.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile, tt.wantFile) {
				t.Errorf("fileService.LoadFile() = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}

func Test_fileService_SearchByDateRange(t *testing.T) {
	type fields struct {
		dao             repository.DAO
		fileStoragePath string
	}
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []datastruct.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileService{
				dao:             tt.fields.dao,
				fileStoragePath: tt.fields.fileStoragePath,
			}
			got, err := s.SearchByDateRange(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileService.SearchByDateRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileService.SearchByDateRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
