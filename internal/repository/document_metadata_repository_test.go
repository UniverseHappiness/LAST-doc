package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"gorm.io/gorm"
)

func TestNewDocumentMetadataRepository(t *testing.T) {
	db := SetupTestDB(t)

	// 测试创建仓库实例
	repo := NewDocumentMetadataRepository(db)

	if repo == nil {
		t.Error("NewDocumentMetadataRepository() 返回 nil")
	}

	// 验证仓库实例类型
	if _, ok := repo.(*documentMetadataRepository); !ok {
		t.Error("NewDocumentMetadataRepository() 返回的类型不正确")
	}
}

// TestDocumentMetadataRepository_Create 测试创建文档元数据功能
func TestDocumentMetadataRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	metadata := CreateTestDocumentMetadata("test-doc-id")

	err := repo.Create(ctx, metadata)
	if err != nil {
		t.Fatalf("创建文档元数据失败: %v", err)
	}

	// 验证元数据是否创建成功
	var savedMeta model.DocumentMetadata
	err = db.Where("id = ?", metadata.ID).First(&savedMeta).Error
	if err != nil {
		t.Fatalf("无法查询创建的文档元数据: %v", err)
	}

	if savedMeta.DocumentID != metadata.DocumentID {
		t.Errorf("文档ID不匹配，预期 %s, 实际 %s", metadata.DocumentID, savedMeta.DocumentID)
	}
}

func Test_documentMetadataRepository_GetByID(t *testing.T) {
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
		want    *model.DocumentMetadata
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &documentMetadataRepository{
				db: tt.fields.db,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentMetadataRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentMetadataRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDocumentMetadataRepository_GetByDocumentID 测试根据文档ID获取元数据列表功能
func TestDocumentMetadataRepository_GetByDocumentID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个元数据记录
	metadata1 := CreateTestDocumentMetadata(documentID)
	metadata1.Metadata["version"] = "1.0.0"

	metadata2 := CreateTestDocumentMetadata(documentID)
	metadata2.Metadata["version"] = "2.0.0"

	err := repo.Create(ctx, metadata1)
	if err != nil {
		t.Fatalf("创建第一个文档元数据失败: %v", err)
	}

	err = repo.Create(ctx, metadata2)
	if err != nil {
		t.Fatalf("创建第二个文档元数据失败: %v", err)
	}

	// 测试获取元数据列表
	metadataList, err := repo.GetByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("获取文档元数据列表失败: %v", err)
	}

	if len(metadataList) != 2 {
		t.Errorf("预期返回2个元数据记录，实际返回 %d 个", len(metadataList))
	}

	// 测试获取不存在文档的元数据
	emptyList, err := repo.GetByDocumentID(ctx, "non-existent-doc")
	if err != nil {
		t.Fatalf("获取不存在文档的元数据不应该返回错误: %v", err)
	}

	if len(emptyList) != 0 {
		t.Errorf("不存在文档的元数据列表应该为空，实际返回 %d 个", len(emptyList))
	}
}

// TestDocumentMetadataRepository_GetLatestMetadata 测试获取最新文档元数据功能
func TestDocumentMetadataRepository_GetLatestMetadata(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个元数据记录（按时间顺序）
	metadata1 := CreateTestDocumentMetadata(documentID)
	metadata1.Metadata["version"] = "1.0.0"

	// 等待一段时间确保时间差
	// metadata2 := CreateTestDocumentMetadata(documentID)
	// metadata2.Metadata["version"] = "2.0.0"

	err := repo.Create(ctx, metadata1)
	if err != nil {
		t.Fatalf("创建第一个文档元数据失败: %v", err)
	}

	// 测试获取最新元数据
	latestMeta, err := repo.GetLatestMetadata(ctx, documentID)
	if err != nil {
		t.Fatalf("获取最新文档元数据失败: %v", err)
	}

	if latestMeta.DocumentID != documentID {
		t.Errorf("文档ID不匹配，预期 %s, 实际 %s", documentID, latestMeta.DocumentID)
	}

	// 测试获取不存在文档的最新元数据
	_, err = repo.GetLatestMetadata(ctx, "non-existent-doc")
	if err == nil {
		t.Error("获取不存在文档的最新元数据应该返回错误")
	}
}

// TestDocumentMetadataRepository_Update 测试更新文档元数据功能
func TestDocumentMetadataRepository_Update(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个元数据记录
	metadata := CreateTestDocumentMetadata(documentID)
	metadata.Metadata["version"] = "1.0.0"

	err := repo.Create(ctx, metadata)
	if err != nil {
		t.Fatalf("创建文档元数据失败: %v", err)
	}

	// 准备更新数据
	updates := map[string]interface{}{
		"version": "2.0.0",
		"updated": true,
		"author":  "updated-author",
	}

	// 更新元数据
	err = repo.Update(ctx, metadata.ID, updates)
	if err != nil {
		t.Fatalf("更新文档元数据失败: %v", err)
	}

	// 验证更新是否成功
	updatedMeta, err := repo.GetByID(ctx, metadata.ID)
	if err != nil {
		t.Fatalf("获取更新后的文档元数据失败: %v", err)
	}

	if updatedMeta.Metadata["version"] != "2.0.0" {
		t.Errorf("元数据版本更新失败，预期 2.0.0, 实际 %v", updatedMeta.Metadata["version"])
	}

	if !updatedMeta.Metadata["updated"].(bool) {
		t.Error("元数据updated字段更新失败")
	}

	if updatedMeta.Metadata["author"] != "updated-author" {
		t.Errorf("元数据author字段更新失败，预期 updated-author, 实际 %v", updatedMeta.Metadata["author"])
	}

	// 测试更新不存在的元数据
	err = repo.Update(ctx, "non-existent-id", updates)
	if err == nil {
		t.Error("更新不存在的元数据应该返回错误")
	}
}

// TestDocumentMetadataRepository_Delete 测试删除文档元数据功能
func TestDocumentMetadataRepository_Delete(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个元数据记录
	metadata := CreateTestDocumentMetadata(documentID)
	metadata.Metadata["version"] = "1.0.0"

	err := repo.Create(ctx, metadata)
	if err != nil {
		t.Fatalf("创建文档元数据失败: %v", err)
	}

	// 验证元数据已创建
	retrievedMeta, err := repo.GetByID(ctx, metadata.ID)
	if err != nil {
		t.Fatalf("获取文档元数据失败: %v", err)
	}
	if retrievedMeta == nil {
		t.Fatal("元数据创建失败")
	}

	// 删除元数据
	err = repo.Delete(ctx, metadata.ID)
	if err != nil {
		t.Fatalf("删除文档元数据失败: %v", err)
	}

	// 验证元数据已被删除
	deletedMeta, err := repo.GetByID(ctx, metadata.ID)
	if err == nil {
		t.Error("删除后的元数据仍然可以获取，删除失败")
	}
	if deletedMeta != nil {
		t.Error("删除后的元数据不应该存在")
	}

	// 测试删除不存在的元数据
	err = repo.Delete(ctx, "non-existent-id")
	if err == nil {
		t.Error("删除不存在的元数据应该返回错误")
	}
}

// TestDocumentMetadataRepository_DeleteByDocumentID 测试按文档ID删除所有元数据功能
func TestDocumentMetadataRepository_DeleteByDocumentID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个元数据记录
	metadata1 := CreateTestDocumentMetadata(documentID)
	metadata1.Metadata["version"] = "1.0.0"

	metadata2 := CreateTestDocumentMetadata(documentID)
	metadata2.Metadata["version"] = "2.0.0"

	err := repo.Create(ctx, metadata1)
	if err != nil {
		t.Fatalf("创建第一个文档元数据失败: %v", err)
	}

	err = repo.Create(ctx, metadata2)
	if err != nil {
		t.Fatalf("创建第二个文档元数据失败: %v", err)
	}

	// 验证元数据已创建
	metaList, err := repo.GetByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("获取文档元数据列表失败: %v", err)
	}
	if len(metaList) != 2 {
		t.Fatalf("预期创建2个元数据记录，实际创建 %d 个", len(metaList))
	}

	// 按文档ID删除所有元数据
	err = repo.DeleteByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("按文档ID删除元数据失败: %v", err)
	}

	// 验证所有元数据已被删除
	deletedMetaList, err := repo.GetByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("获取删除后的文档元数据列表失败: %v", err)
	}
	if len(deletedMetaList) != 0 {
		t.Errorf("删除后仍有 %d 个元数据记录存在，删除失败", len(deletedMetaList))
	}

	// 测试删除不存在文档的元数据
	err = repo.DeleteByDocumentID(ctx, "non-existent-doc")
	if err == nil {
		t.Error("删除不存在文档的元数据应该返回错误")
	}
}

// TestDocumentMetadataRepository_Count 测试统计文档元数据数量功能
func TestDocumentMetadataRepository_Count(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentMetadataRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 初始时应该没有元数据记录
	count, err := repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档元数据数量失败: %v", err)
	}
	if count != 0 {
		t.Errorf("初始状态下元数据数量应该为0，实际为 %d", count)
	}

	// 创建第一个元数据记录
	metadata1 := CreateTestDocumentMetadata(documentID)
	metadata1.Metadata["version"] = "1.0.0"

	err = repo.Create(ctx, metadata1)
	if err != nil {
		t.Fatalf("创建第一个文档元数据失败: %v", err)
	}

	// 验证数量为1
	count, err = repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档元数据数量失败: %v", err)
	}
	if count != 1 {
		t.Errorf("创建第一个元数据后数量应该为1，实际为 %d", count)
	}

	// 创建第二个元数据记录
	metadata2 := CreateTestDocumentMetadata(documentID)
	metadata2.Metadata["version"] = "2.0.0"

	err = repo.Create(ctx, metadata2)
	if err != nil {
		t.Fatalf("创建第二个文档元数据失败: %v", err)
	}

	// 验证数量为2
	count, err = repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档元数据数量失败: %v", err)
	}
	if count != 2 {
		t.Errorf("创建第二个元数据后数量应该为2，实际为 %d", count)
	}

	// 测试统计不存在文档的元数据数量
	count, err = repo.Count(ctx, "non-existent-doc")
	if err != nil {
		t.Fatalf("统计不存在文档的元数据数量失败: %v", err)
	}
	if count != 0 {
		t.Errorf("不存在文档的元数据数量应该为0，实际为 %d", count)
	}
}
