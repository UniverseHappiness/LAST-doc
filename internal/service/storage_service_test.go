package service

import (
	"context"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestNewLocalStorageService(t *testing.T) {
	type args struct {
		baseDir string
	}
	tests := []struct {
		name string
		args args
		want *LocalStorageService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocalStorageService(tt.args.baseDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocalStorageService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorageService_SaveFile(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx  context.Context
		file *multipart.FileHeader
		path string
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
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if err := s.SaveFile(tt.args.ctx, tt.args.file, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("LocalStorageService.SaveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorageService_DeleteFile(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx  context.Context
		path string
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
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if err := s.DeleteFile(tt.args.ctx, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("LocalStorageService.DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorageService_GetFile(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			got, err := s.GetFile(tt.args.ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LocalStorageService.GetFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalStorageService.GetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorageService_GenerateFilePath(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		documentID string
		fileName   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if got := s.GenerateFilePath(tt.args.documentID, tt.args.fileName); got != tt.want {
				t.Errorf("LocalStorageService.GenerateFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorageService_GetBaseDir(t *testing.T) {
	type fields struct {
		baseDir string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if got := s.GetBaseDir(); got != tt.want {
				t.Errorf("LocalStorageService.GetBaseDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorageService_FileExists(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if got := s.FileExists(tt.args.ctx, tt.args.path); got != tt.want {
				t.Errorf("LocalStorageService.FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorageService_CopyFile(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx     context.Context
		srcPath string
		dstPath string
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
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if err := s.CopyFile(tt.args.ctx, tt.args.srcPath, tt.args.dstPath); (err != nil) != tt.wantErr {
				t.Errorf("LocalStorageService.CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorageService_MoveFile(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx     context.Context
		srcPath string
		dstPath string
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
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			if err := s.MoveFile(tt.args.ctx, tt.args.srcPath, tt.args.dstPath); (err != nil) != tt.wantErr {
				t.Errorf("LocalStorageService.MoveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorageService_GetFileSize(t *testing.T) {
	type fields struct {
		baseDir string
	}
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorageService{
				baseDir: tt.fields.baseDir,
			}
			got, err := s.GetFileSize(tt.args.ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LocalStorageService.GetFileSize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("LocalStorageService.GetFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
