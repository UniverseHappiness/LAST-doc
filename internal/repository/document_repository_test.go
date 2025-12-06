package repository

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用PostgreSQL的内存数据库用于测试
	dsn := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果连接失败，尝试使用内存中的PostgreSQL
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Skipf("无法连接测试数据库，跳过测试: %v", err)
		}
	}

	// 创建测试数据库
	err = db.Exec("CREATE DATABASE IF NOT EXISTS test_db").Error
	if err != nil {
		// 如果创建数据库失败，尝试使用默认数据库
		t.Logf("创建测试数据库失败，使用默认数据库: %v", err)
	} else {
		// 重新连接到新创建的数据库
		dsn = "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Skipf("无法连接到新创建的测试数据库，跳过测试: %v", err)
		}
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Document{}, &model.DocumentVersion{}, &model.DocumentMetadata{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 清理测试数据
	db.Exec("DELETE FROM document_versions")
	db.Exec("DELETE FROM document_metadata")
	db.Exec("DELETE FROM documents")

	return db
}

// TestDocumentRepository_Create 测试创建文档功能
func TestDocumentRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()
	document := &model.Document{
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Tags:        []string{"测试", "文档"},
		FilePath:    "/tmp/test.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "这是一个测试文档",
		Library:     "测试库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, document)
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
	}

	// 验证文档是否创建成功
	var savedDoc model.Document
	err = db.Where("id = ?", document.ID).First(&savedDoc).Error
	if err != nil {
		t.Fatalf("无法查询创建的文档: %v", err)
	}

	if savedDoc.Name != document.Name {
		t.Errorf("文档名称不匹配，预期 %s, 实际 %s", document.Name, savedDoc.Name)
	}

	if savedDoc.Type != document.Type {
		t.Errorf("文档类型不匹配，预期 %v, 实际 %v", document.Type, savedDoc.Type)
	}
}

// TestDocumentRepository_GetByID 测试根据ID获取文档功能
func TestDocumentRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()
	document := &model.Document{
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Tags:        []string{"测试", "文档"},
		FilePath:    "/tmp/test.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "这是一个测试文档",
		Library:     "测试库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 先创建文档
	err := repo.Create(ctx, document)
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
	}

	// 测试获取文档
	retrievedDoc, err := repo.GetByID(ctx, document.ID)
	if err != nil {
		t.Fatalf("获取文档失败: %v", err)
	}

	if retrievedDoc.Name != document.Name {
		t.Errorf("文档名称不匹配，预期 %s, 实际 %s", document.Name, retrievedDoc.Name)
	}

	if retrievedDoc.Type != document.Type {
		t.Errorf("文档类型不匹配，预期 %v, 实际 %v", document.Type, retrievedDoc.Type)
	}
}

// TestDocumentRepository_List 测试获取文档列表功能
func TestDocumentRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()

	// 创建测试文档
	for i := 0; i < 5; i++ {
		document := &model.Document{
			Name:        fmt.Sprintf("测试文档%d", i),
			Type:        model.DocumentTypeMarkdown,
			Version:     "1.0.0",
			Tags:        []string{"测试", "文档"},
			FilePath:    fmt.Sprintf("/tmp/test%d.md", i),
			FileSize:    1024,
			Status:      model.DocumentStatusCompleted,
			Description: fmt.Sprintf("这是第%d个测试文档", i),
			Library:     "测试库",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(ctx, document)
		if err != nil {
			t.Fatalf("创建文档失败: %v", err)
		}
	}

	// 测试获取文档列表
	documents, total, err := repo.List(ctx, 1, 10, nil)
	if err != nil {
		t.Fatalf("获取文档列表失败: %v", err)
	}

	if total != 5 {
		t.Errorf("文档总数不匹配，预期 5, 实际 %d", total)
	}

	if len(documents) != 5 {
		t.Errorf("返回的文档数量不匹配，预期 5, 实际 %d", len(documents))
	}

	// 测试分页功能
	documents, total, err = repo.List(ctx, 1, 2, nil)
	if err != nil {
		t.Fatalf("获取文档列表失败: %v", err)
	}

	if len(documents) != 2 {
		t.Errorf("分页返回的文档数量不匹配，预期 2, 实际 %d", len(documents))
	}

	// 测试过滤功能
	documents, total, err = repo.List(ctx, 1, 10, map[string]interface{}{"name": "测试文档1"})
	if err != nil {
		t.Fatalf("获取文档列表失败: %v", err)
	}

	if total != 1 {
		t.Errorf("过滤后的文档总数不匹配，预期 1, 实际 %d", total)
	}

	if len(documents) != 1 {
		t.Errorf("过滤后返回的文档数量不匹配，预期 1, 实际 %d", len(documents))
	}

	if documents[0].Name != "测试文档1" {
		t.Errorf("过滤后的文档名称不匹配，预期 测试文档1, 实际 %s", documents[0].Name)
	}
}

// TestDocumentRepository_Update 测试更新文档功能
func TestDocumentRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()
	document := &model.Document{
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Tags:        []string{"测试", "文档"},
		FilePath:    "/tmp/test.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "这是一个测试文档",
		Library:     "测试库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 先创建文档
	err := repo.Create(ctx, document)
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
	}

	// 更新文档
	updates := map[string]interface{}{
		"name":        "更新后的文档",
		"description": "这是更新后的描述",
		"tags":        []string{"更新", "文档"},
	}

	err = repo.Update(ctx, document.ID, updates)
	if err != nil {
		t.Fatalf("更新文档失败: %v", err)
	}

	// 验证更新结果
	updatedDoc, err := repo.GetByID(ctx, document.ID)
	if err != nil {
		t.Fatalf("获取更新后的文档失败: %v", err)
	}

	if updatedDoc.Name != "更新后的文档" {
		t.Errorf("更新后的文档名称不匹配，预期 更新后的文档, 实际 %s", updatedDoc.Name)
	}

	if updatedDoc.Description != "这是更新后的描述" {
		t.Errorf("更新后的文档描述不匹配，预期 这是更新后的描述, 实际 %s", updatedDoc.Description)
	}
}

// TestDocumentRepository_Delete 测试删除文档功能
func TestDocumentRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()
	document := &model.Document{
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Tags:        []string{"测试", "文档"},
		FilePath:    "/tmp/test.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "这是一个测试文档",
		Library:     "测试库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 先创建文档
	err := repo.Create(ctx, document)
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
	}

	// 删除文档
	err = repo.Delete(ctx, document.ID)
	if err != nil {
		t.Fatalf("删除文档失败: %v", err)
	}

	// 验证文档是否已删除
	_, err = repo.GetByID(ctx, document.ID)
	if err == nil {
		t.Error("文档未被删除")
	}
}

// TestDocumentRepository_GetByLibrary 测试根据库获取文档功能
func TestDocumentRepository_GetByLibrary(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()

	// 创建测试文档
	for i := 0; i < 3; i++ {
		document := &model.Document{
			Name:        fmt.Sprintf("测试文档%d", i),
			Type:        model.DocumentTypeMarkdown,
			Version:     "1.0.0",
			Tags:        []string{"测试", "文档"},
			FilePath:    fmt.Sprintf("/tmp/test%d.md", i),
			FileSize:    1024,
			Status:      model.DocumentStatusCompleted,
			Description: fmt.Sprintf("这是第%d个测试文档", i),
			Library:     "测试库",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(ctx, document)
		if err != nil {
			t.Fatalf("创建文档失败: %v", err)
		}
	}

	// 创建不同库的文档
	otherDoc := &model.Document{
		Name:        "其他库文档",
		Type:        model.DocumentTypePDF,
		Version:     "1.0.0",
		Tags:        []string{"其他", "文档"},
		FilePath:    "/tmp/other.pdf",
		FileSize:    2048,
		Status:      model.DocumentStatusCompleted,
		Description: "这是其他库的文档",
		Library:     "其他库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, otherDoc)
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
	}

	// 测试根据库获取文档
	documents, total, err := repo.GetByLibrary(ctx, "测试库", 1, 10)
	if err != nil {
		t.Fatalf("根据库获取文档失败: %v", err)
	}

	if total != 3 {
		t.Errorf("根据库获取的文档总数不匹配，预期 3, 实际 %d", total)
	}

	if len(documents) != 3 {
		t.Errorf("根据库返回的文档数量不匹配，预期 3, 实际 %d", len(documents))
	}

	// 验证所有文档都属于指定库
	for _, doc := range documents {
		if doc.Library != "测试库" {
			t.Errorf("文档库不匹配，预期 测试库, 实际 %s", doc.Library)
		}
	}
}

// TestDocumentRepository_Count 测试获取文档总数功能
func TestDocumentRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDocumentRepository(db)

	ctx := context.Background()

	// 创建测试文档
	for i := 0; i < 5; i++ {
		document := &model.Document{
			Name:        fmt.Sprintf("测试文档%d", i),
			Type:        model.DocumentTypeMarkdown,
			Version:     "1.0.0",
			Tags:        []string{"测试", "文档"},
			FilePath:    fmt.Sprintf("/tmp/test%d.md", i),
			FileSize:    1024,
			Status:      model.DocumentStatusCompleted,
			Description: fmt.Sprintf("这是第%d个测试文档", i),
			Library:     "测试库",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(ctx, document)
		if err != nil {
			t.Fatalf("创建文档失败: %v", err)
		}
	}

	// 测试获取文档总数
	total, err := repo.Count(ctx, nil)
	if err != nil {
		t.Fatalf("获取文档总数失败: %v", err)
	}

	if total != 5 {
		t.Errorf("文档总数不匹配，预期 5, 实际 %d", total)
	}

	// 测试带过滤条件获取文档总数
	total, err = repo.Count(ctx, map[string]interface{}{"library": "测试库"})
	if err != nil {
		t.Fatalf("获取过滤后的文档总数失败: %v", err)
	}

	if total != 5 {
		t.Errorf("过滤后的文档总数不匹配，预期 5, 实际 %d", total)
	}

	total, err = repo.Count(ctx, map[string]interface{}{"library": "不存在的库"})
	if err != nil {
		t.Fatalf("获取过滤后的文档总数失败: %v", err)
	}

	if total != 0 {
		t.Errorf("不存在的库的文档总数不匹配，预期 0, 实际 %d", total)
	}
}

func TestNewDocumentRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want DocumentRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDocumentRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocumentRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentRepository_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx      context.Context
		document *model.Document
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
			r := &documentRepository{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("documentRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentRepository_GetByID(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentRepository_List(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		page    int
		size    int
		filters map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.List(tt.args.ctx, tt.args.page, tt.args.size, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_Update(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		id      string
		updates map[string]interface{}
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
			r := &documentRepository{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.id, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("documentRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentRepository_Delete(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  string
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
			r := &documentRepository{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("documentRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentRepository_GetByLibrary(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		library string
		page    int
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.GetByLibrary(tt.args.ctx, tt.args.library, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByLibrary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByLibrary() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.GetByLibrary() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_GetByType(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		docType model.DocumentType
		page    int
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.GetByType(tt.args.ctx, tt.args.docType, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.GetByType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_GetByTag(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx  context.Context
		tag  string
		page int
		size int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.GetByTag(tt.args.ctx, tt.args.tag, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByTag() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.GetByTag() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_GetByVersion(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		version string
		page    int
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.GetByVersion(tt.args.ctx, tt.args.version, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.GetByVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_GetByLibraryAndVersion(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		library string
		version string
		page    int
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.GetByLibraryAndVersion(tt.args.ctx, tt.args.library, tt.args.version, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByLibraryAndVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByLibraryAndVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.GetByLibraryAndVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_GetByName(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx  context.Context
		name string
		page int
		size int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.GetByName(tt.args.ctx, tt.args.name, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.GetByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentRepository.GetByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentRepository.GetByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentRepository_Count(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		filters map[string]interface{}
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
			r := &documentRepository{
				db: tt.fields.db,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentRepository.Count() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("documentRepository.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
