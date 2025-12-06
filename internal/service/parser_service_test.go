package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

func TestNewParserService(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParserService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserService_RegisterParser(t *testing.T) {
	type fields struct {
		parsers map[model.DocumentType]DocumentParser
	}
	type args struct {
		docType model.DocumentType
		parser  DocumentParser
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &parserService{
				parsers: tt.fields.parsers,
			}
			s.RegisterParser(tt.args.docType, tt.args.parser)
		})
	}
}

func Test_parserService_ParseDocument(t *testing.T) {
	type fields struct {
		parsers map[model.DocumentType]DocumentParser
	}
	type args struct {
		ctx      context.Context
		filePath string
		docType  model.DocumentType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &parserService{
				parsers: tt.fields.parsers,
			}
			got, got1, err := s.ParseDocument(tt.args.ctx, tt.args.filePath, tt.args.docType)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parserService.ParseDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("parserService.ParseDocument() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parserService.ParseDocument() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewMarkdownParser(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMarkdownParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMarkdownParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_markdownParser_Parse(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		p       *markdownParser
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &markdownParser{}
			got, got1, err := p.Parse(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("markdownParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("markdownParser.Parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("markdownParser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_markdownParser_SupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    *markdownParser
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &markdownParser{}
			if got := p.SupportedExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markdownParser.SupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractMarkdownMetadata(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractMarkdownMetadata(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractMarkdownMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPDFParser(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPDFParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPDFParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pdfParser_Parse(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		p       *pdfParser
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pdfParser{}
			got, got1, err := p.Parse(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("pdfParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("pdfParser.Parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("pdfParser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_pdfParser_SupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    *pdfParser
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pdfParser{}
			if got := p.SupportedExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pdfParser.SupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDocxParser(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDocxParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocxParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_docxParser_Parse(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		p       *docxParser
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &docxParser{}
			got, got1, err := p.Parse(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("docxParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("docxParser.Parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("docxParser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_docxParser_SupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    *docxParser
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &docxParser{}
			if got := p.SupportedExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("docxParser.SupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSwaggerParser(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSwaggerParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSwaggerParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_swaggerParser_Parse(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		p       *swaggerParser
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &swaggerParser{}
			got, got1, err := p.Parse(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("swaggerParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("swaggerParser.Parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("swaggerParser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_swaggerParser_SupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    *swaggerParser
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &swaggerParser{}
			if got := p.SupportedExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("swaggerParser.SupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractSwaggerMetadata(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractSwaggerMetadata(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractSwaggerMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOpenAPIParser(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOpenAPIParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOpenAPIParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_openAPIParser_Parse(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		p       *openAPIParser
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &openAPIParser{}
			got, got1, err := p.Parse(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("openAPIParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("openAPIParser.Parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("openAPIParser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_openAPIParser_SupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    *openAPIParser
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &openAPIParser{}
			if got := p.SupportedExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openAPIParser.SupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJavaDocParser(t *testing.T) {
	tests := []struct {
		name string
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJavaDocParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJavaDocParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_javaDocParser_Parse(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		p       *javaDocParser
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &javaDocParser{}
			got, got1, err := p.Parse(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("javaDocParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("javaDocParser.Parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("javaDocParser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_javaDocParser_SupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    *javaDocParser
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &javaDocParser{}
			if got := p.SupportedExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("javaDocParser.SupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractJavaDocMetadata(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractJavaDocMetadata(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractJavaDocMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetParserByExtension(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want DocumentParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetParserByExtension(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParserByExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
