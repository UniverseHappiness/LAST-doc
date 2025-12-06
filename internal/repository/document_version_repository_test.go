package repository

import (
	"context"
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

// TestNewDocumentVersionRepository 测试创建文档版本仓库实例
func TestNewDocumentVersionRepository(t *testing.T) {
	db := SetupTestDB(t)

	repo := NewDocumentVersionRepository(db)

	if repo == nil {
		t.Fatal("NewDocumentVersionRepository 返回了 nil")
	}

	// 检查返回的类型是否正确
	if _, ok := repo.(*documentVersionRepository); !ok {
		t.Errorf("NewDocumentVersionRepository 返回的类型错误，预期 *documentVersionRepository，实际 %T", repo)
	}

	// 检查数据库连接是否正确设置
	versionRepo, ok := repo.(*documentVersionRepository)
	if !ok {
		t.Fatal("类型断言失败")
	}

	if versionRepo.db == nil {
		t.Error("文档版本仓库的数据库连接未正确设置")
	}

	if versionRepo.db != db {
		t.Error("文档版本仓库的数据库连接与传入的数据库连接不一致")
	}
}

// TestDocumentVersionRepository_Create 测试创建文档版本功能
func TestDocumentVersionRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个文档版本记录
	version := CreateTestDocumentVersion(documentID)
	version.Version = "1.0.0"
	version.Content = "这是第一个版本的文档内容"

	err := repo.Create(ctx, version)
	if err != nil {
		t.Fatalf("创建文档版本失败: %v", err)
	}

	// 验证版本是否创建成功
	retrievedVersion, err := repo.GetByID(ctx, version.ID)
	if err != nil {
		t.Fatalf("获取创建的文档版本失败: %v", err)
	}

	if retrievedVersion.DocumentID != documentID {
		t.Errorf("文档ID不匹配，预期 %s, 实际 %s", documentID, retrievedVersion.DocumentID)
	}

	if retrievedVersion.Version != "1.0.0" {
		t.Errorf("版本号不匹配，预期 1.0.0, 实际 %s", retrievedVersion.Version)
	}

	if retrievedVersion.Content != "这是第一个版本的文档内容" {
		t.Errorf("文档内容不匹配，预期 '这是第一个版本的文档内容', 实际 %s", retrievedVersion.Content)
	}

	// 测试创建重复ID的版本
	duplicateVersion := CreateTestDocumentVersion(documentID)
	duplicateVersion.ID = version.ID // 使用相同的ID

	err = repo.Create(ctx, duplicateVersion)
	if err == nil {
		t.Error("创建重复ID的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_GetByID 测试根据ID获取文档版本功能
func TestDocumentVersionRepository_GetByID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个文档版本记录
	version := CreateTestDocumentVersion(documentID)
	version.Version = "1.0.0"
	version.Content = "测试文档内容"

	err := repo.Create(ctx, version)
	if err != nil {
		t.Fatalf("创建文档版本失败: %v", err)
	}

	// 测试获取已存在的版本
	retrievedVersion, err := repo.GetByID(ctx, version.ID)
	if err != nil {
		t.Fatalf("获取文档版本失败: %v", err)
	}

	if retrievedVersion.DocumentID != documentID {
		t.Errorf("文档ID不匹配，预期 %s, 实际 %s", documentID, retrievedVersion.DocumentID)
	}

	if retrievedVersion.Version != "1.0.0" {
		t.Errorf("版本号不匹配，预期 1.0.0, 实际 %s", retrievedVersion.Version)
	}

	if retrievedVersion.Content != "测试文档内容" {
		t.Errorf("文档内容不匹配，预期 '测试文档内容', 实际 %s", retrievedVersion.Content)
	}

	// 测试获取不存在的版本
	_, err = repo.GetByID(ctx, "non-existent-id")
	if err == nil {
		t.Error("获取不存在的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_GetByDocumentID 测试根据文档ID获取所有版本功能
func TestDocumentVersionRepository_GetByDocumentID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个文档版本记录
	version1 := CreateTestDocumentVersion(documentID)
	version1.Version = "1.0.0"
	version1.Content = "第一个版本的内容"

	version2 := CreateTestDocumentVersion(documentID)
	version2.Version = "2.0.0"
	version2.Content = "第二个版本的内容"

	version3 := CreateTestDocumentVersion(documentID)
	version3.Version = "3.0.0"
	version3.Content = "第三个版本的内容"

	err := repo.Create(ctx, version1)
	if err != nil {
		t.Fatalf("创建第一个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version2)
	if err != nil {
		t.Fatalf("创建第二个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version3)
	if err != nil {
		t.Fatalf("创建第三个文档版本失败: %v", err)
	}

	// 测试获取文档的所有版本
	versions, err := repo.GetByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("获取文档版本列表失败: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("预期获取3个版本，实际获取 %d 个", len(versions))
	}

	// 验证版本内容
	versionMap := make(map[string]*model.DocumentVersion)
	for _, v := range versions {
		versionMap[v.Version] = v
	}

	if versionMap["1.0.0"].Content != "第一个版本的内容" {
		t.Errorf("版本1.0.0内容不匹配，预期 '第一个版本的内容', 实际 %s", versionMap["1.0.0"].Content)
	}

	if versionMap["2.0.0"].Content != "第二个版本的内容" {
		t.Errorf("版本2.0.0内容不匹配，预期 '第二个版本的内容', 实际 %s", versionMap["2.0.0"].Content)
	}

	if versionMap["3.0.0"].Content != "第三个版本的内容" {
		t.Errorf("版本3.0.0内容不匹配，预期 '第三个版本的内容', 实际 %s", versionMap["3.0.0"].Content)
	}

	// 测试获取不存在文档的版本
	versions, err = repo.GetByDocumentID(ctx, "non-existent-doc")
	if err != nil {
		t.Fatalf("获取不存在文档的版本失败: %v", err)
	}
	if len(versions) != 0 {
		t.Errorf("不存在文档的版本数量应该为0，实际为 %d", len(versions))
	}
}

// TestDocumentVersionRepository_GetByDocumentIDAndVersion 测试根据文档ID和版本号获取文档版本功能
func TestDocumentVersionRepository_GetByDocumentIDAndVersion(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个文档版本记录
	version1 := CreateTestDocumentVersion(documentID)
	version1.Version = "1.0.0"
	version1.Content = "第一个版本的内容"

	version2 := CreateTestDocumentVersion(documentID)
	version2.Version = "2.0.0"
	version2.Content = "第二个版本的内容"

	err := repo.Create(ctx, version1)
	if err != nil {
		t.Fatalf("创建第一个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version2)
	if err != nil {
		t.Fatalf("创建第二个文档版本失败: %v", err)
	}

	// 测试获取存在的版本
	retrievedVersion, err := repo.GetByDocumentIDAndVersion(ctx, documentID, "1.0.0")
	if err != nil {
		t.Fatalf("获取文档版本失败: %v", err)
	}

	if retrievedVersion.DocumentID != documentID {
		t.Errorf("文档ID不匹配，预期 %s, 实际 %s", documentID, retrievedVersion.DocumentID)
	}

	if retrievedVersion.Version != "1.0.0" {
		t.Errorf("版本号不匹配，预期 1.0.0, 实际 %s", retrievedVersion.Version)
	}

	if retrievedVersion.Content != "第一个版本的内容" {
		t.Errorf("文档内容不匹配，预期 '第一个版本的内容', 实际 %s", retrievedVersion.Content)
	}

	// 测试获取第二个版本
	retrievedVersion, err = repo.GetByDocumentIDAndVersion(ctx, documentID, "2.0.0")
	if err != nil {
		t.Fatalf("获取第二个文档版本失败: %v", err)
	}

	if retrievedVersion.Content != "第二个版本的内容" {
		t.Errorf("第二个版本内容不匹配，预期 '第二个版本的内容', 实际 %s", retrievedVersion.Content)
	}

	// 测试获取不存在的版本
	_, err = repo.GetByDocumentIDAndVersion(ctx, documentID, "3.0.0")
	if err == nil {
		t.Error("获取不存在的版本应该返回错误")
	}

	// 测试获取不存在文档的版本
	_, err = repo.GetByDocumentIDAndVersion(ctx, "non-existent-doc", "1.0.0")
	if err == nil {
		t.Error("获取不存在文档的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_GetLatestVersion 测试获取最新文档版本功能
func TestDocumentVersionRepository_GetLatestVersion(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个文档版本记录（按时间顺序）
	version1 := CreateTestDocumentVersion(documentID)
	version1.Version = "1.0.0"
	version1.Content = "第一个版本的内容"

	version2 := CreateTestDocumentVersion(documentID)
	version2.Version = "2.0.0"
	version2.Content = "第二个版本的内容"

	err := repo.Create(ctx, version1)
	if err != nil {
		t.Fatalf("创建第一个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version2)
	if err != nil {
		t.Fatalf("创建第二个文档版本失败: %v", err)
	}

	// 测试获取最新版本
	latestVersion, err := repo.GetLatestVersion(ctx, documentID)
	if err != nil {
		t.Fatalf("获取最新文档版本失败: %v", err)
	}

	if latestVersion.DocumentID != documentID {
		t.Errorf("文档ID不匹配，预期 %s, 实际 %s", documentID, latestVersion.DocumentID)
	}

	if latestVersion.Version != "2.0.0" {
		t.Errorf("最新版本号不匹配，预期 2.0.0, 实际 %s", latestVersion.Version)
	}

	if latestVersion.Content != "第二个版本的内容" {
		t.Errorf("最新版本内容不匹配，预期 '第二个版本的内容', 实际 %s", latestVersion.Content)
	}

	// 测试获取不存在文档的最新版本
	_, err = repo.GetLatestVersion(ctx, "non-existent-doc")
	if err == nil {
		t.Error("获取不存在文档的最新版本应该返回错误")
	}
}

// TestDocumentVersionRepository_GetVersionsByStatus 测试根据状态获取文档版本功能
func TestDocumentVersionRepository_GetVersionsByStatus(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个不同状态的文档版本记录
	version1 := CreateTestDocumentVersion(documentID)
	version1.Version = "1.0.0"
	version1.Content = "草稿版本的内容"
	version1.Status = model.DocumentStatusProcessing

	version2 := CreateTestDocumentVersion(documentID)
	version2.Version = "2.0.0"
	version2.Content = "已发布版本的内容"
	version2.Status = model.DocumentStatusCompleted

	version3 := CreateTestDocumentVersion(documentID)
	version3.Version = "3.0.0"
	version3.Content = "另一个草稿版本的内容"
	version3.Status = model.DocumentStatusProcessing

	err := repo.Create(ctx, version1)
	if err != nil {
		t.Fatalf("创建第一个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version2)
	if err != nil {
		t.Fatalf("创建第二个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version3)
	if err != nil {
		t.Fatalf("创建第三个文档版本失败: %v", err)
	}

	// 测试获取草稿状态的版本
	processingVersions, err := repo.GetVersionsByStatus(ctx, documentID, model.DocumentStatusProcessing)
	if err != nil {
		t.Fatalf("获取处理中状态版本失败: %v", err)
	}

	if len(processingVersions) != 2 {
		t.Errorf("预期获取2个处理中状态版本，实际获取 %d 个", len(processingVersions))
	}

	// 验证处理中版本内容
	for _, v := range processingVersions {
		if v.Status != model.DocumentStatusProcessing {
			t.Errorf("获取的版本状态不正确，预期 %v, 实际 %v", model.DocumentStatusProcessing, v.Status)
		}
	}

	// 测试获取已完成状态的版本
	completedVersions, err := repo.GetVersionsByStatus(ctx, documentID, model.DocumentStatusCompleted)
	if err != nil {
		t.Fatalf("获取已完成状态版本失败: %v", err)
	}

	if len(completedVersions) != 1 {
		t.Errorf("预期获取1个已完成状态版本，实际获取 %d 个", len(completedVersions))
	}

	if completedVersions[0].Status != model.DocumentStatusCompleted {
		t.Errorf("获取的版本状态不正确，预期 %v, 实际 %v", model.DocumentStatusCompleted, completedVersions[0].Status)
	}

	if completedVersions[0].Content != "已发布版本的内容" {
		t.Errorf("已完成版本内容不匹配，预期 '已发布版本的内容', 实际 %s", completedVersions[0].Content)
	}

	// 测试获取不存在的状态
	failedVersions, err := repo.GetVersionsByStatus(ctx, documentID, model.DocumentStatusFailed)
	if err != nil {
		t.Fatalf("获取失败状态版本失败: %v", err)
	}
	if len(failedVersions) != 0 {
		t.Errorf("失败状态版本数量应该为0，实际为 %d", len(failedVersions))
	}
}

// TestDocumentVersionRepository_Update 测试更新文档版本功能
func TestDocumentVersionRepository_Update(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个文档版本记录
	version := CreateTestDocumentVersion(documentID)
	version.Version = "1.0.0"
	version.Content = "原始版本内容"
	version.Status = model.DocumentStatusProcessing

	err := repo.Create(ctx, version)
	if err != nil {
		t.Fatalf("创建文档版本失败: %v", err)
	}

	// 准备更新数据
	updates := map[string]interface{}{
		"status":  model.DocumentStatusCompleted,
		"content": "更新后的版本内容",
	}

	// 更新文档版本
	err = repo.Update(ctx, version.ID, updates)
	if err != nil {
		t.Fatalf("更新文档版本失败: %v", err)
	}

	// 验证更新是否成功
	updatedVersion, err := repo.GetByID(ctx, version.ID)
	if err != nil {
		t.Fatalf("获取更新后的文档版本失败: %v", err)
	}

	if updatedVersion.Status != model.DocumentStatusCompleted {
		t.Errorf("版本状态更新失败，预期 %v, 实际 %v", model.DocumentStatusCompleted, updatedVersion.Status)
	}

	if updatedVersion.Content != "更新后的版本内容" {
		t.Errorf("版本内容更新失败，预期 '更新后的版本内容', 实际 %s", updatedVersion.Content)
	}

	// 测试更新不存在的版本
	err = repo.Update(ctx, "non-existent-id", updates)
	if err == nil {
		t.Error("更新不存在的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_UpdateContent 测试更新文档版本内容和状态功能
func TestDocumentVersionRepository_UpdateContent(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个文档版本记录
	version := CreateTestDocumentVersion(documentID)
	version.Version = "1.0.0"
	version.Content = "原始版本内容"
	version.Status = model.DocumentStatusProcessing

	err := repo.Create(ctx, version)
	if err != nil {
		t.Fatalf("创建文档版本失败: %v", err)
	}

	// 更新版本内容和状态
	newContent := "更新后的版本内容"
	newStatus := model.DocumentStatusCompleted

	err = repo.UpdateContent(ctx, documentID, "1.0.0", newContent, newStatus)
	if err != nil {
		t.Fatalf("更新文档版本内容和状态失败: %v", err)
	}

	// 验证更新是否成功
	updatedVersion, err := repo.GetByDocumentIDAndVersion(ctx, documentID, "1.0.0")
	if err != nil {
		t.Fatalf("获取更新后的文档版本失败: %v", err)
	}

	if updatedVersion.Status != newStatus {
		t.Errorf("版本状态更新失败，预期 %v, 实际 %v", newStatus, updatedVersion.Status)
	}

	if updatedVersion.Content != newContent {
		t.Errorf("版本内容更新失败，预期 '%s', 实际 '%s'", newContent, updatedVersion.Content)
	}

	// 测试更新不存在版本的文档
	err = repo.UpdateContent(ctx, documentID, "2.0.0", newContent, newStatus)
	if err == nil {
		t.Error("更新不存在版本的文档应该返回错误")
	}

	// 测试更新不存在文档的版本
	err = repo.UpdateContent(ctx, "non-existent-doc", "1.0.0", newContent, newStatus)
	if err == nil {
		t.Error("更新不存在文档的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_UpdateStatus 测试更新文档版本状态功能
func TestDocumentVersionRepository_UpdateStatus(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个文档版本记录
	version := CreateTestDocumentVersion(documentID)
	version.Version = "1.0.0"
	version.Content = "测试版本内容"
	version.Status = model.DocumentStatusProcessing

	err := repo.Create(ctx, version)
	if err != nil {
		t.Fatalf("创建文档版本失败: %v", err)
	}

	// 更新版本状态
	newStatus := model.DocumentStatusCompleted

	err = repo.UpdateStatus(ctx, documentID, "1.0.0", newStatus)
	if err != nil {
		t.Fatalf("更新文档版本状态失败: %v", err)
	}

	// 验证更新是否成功
	updatedVersion, err := repo.GetByDocumentIDAndVersion(ctx, documentID, "1.0.0")
	if err != nil {
		t.Fatalf("获取更新后的文档版本失败: %v", err)
	}

	if updatedVersion.Status != newStatus {
		t.Errorf("版本状态更新失败，预期 %v, 实际 %v", newStatus, updatedVersion.Status)
	}

	// 验证内容未被修改
	if updatedVersion.Content != "测试版本内容" {
		t.Errorf("版本内容被意外修改，预期 '测试版本内容', 实际 '%s'", updatedVersion.Content)
	}

	// 测试更新不存在版本的文档
	err = repo.UpdateStatus(ctx, documentID, "2.0.0", newStatus)
	if err == nil {
		t.Error("更新不存在版本的文档应该返回错误")
	}

	// 测试更新不存在文档的版本
	err = repo.UpdateStatus(ctx, "non-existent-doc", "1.0.0", newStatus)
	if err == nil {
		t.Error("更新不存在文档的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_Delete 测试删除文档版本功能
func TestDocumentVersionRepository_Delete(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建一个文档版本记录
	version := CreateTestDocumentVersion(documentID)
	version.Version = "1.0.0"
	version.Content = "测试版本内容"

	err := repo.Create(ctx, version)
	if err != nil {
		t.Fatalf("创建文档版本失败: %v", err)
	}

	// 验证版本已创建
	retrievedVersion, err := repo.GetByID(ctx, version.ID)
	if err != nil {
		t.Fatalf("获取文档版本失败: %v", err)
	}
	if retrievedVersion == nil {
		t.Fatal("版本创建失败")
	}

	// 删除版本
	err = repo.Delete(ctx, version.ID)
	if err != nil {
		t.Fatalf("删除文档版本失败: %v", err)
	}

	// 验证版本已被删除
	deletedVersion, err := repo.GetByID(ctx, version.ID)
	if err == nil {
		t.Error("删除后的版本仍然可以获取，删除失败")
	}
	if deletedVersion != nil {
		t.Error("删除后的版本不应该存在")
	}

	// 测试删除不存在的版本
	err = repo.Delete(ctx, "non-existent-id")
	if err == nil {
		t.Error("删除不存在的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_DeleteByDocumentID 测试按文档ID删除所有版本功能
func TestDocumentVersionRepository_DeleteByDocumentID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 创建多个文档版本记录
	version1 := CreateTestDocumentVersion(documentID)
	version1.Version = "1.0.0"
	version1.Content = "第一个版本的内容"

	version2 := CreateTestDocumentVersion(documentID)
	version2.Version = "2.0.0"
	version2.Content = "第二个版本的内容"

	version3 := CreateTestDocumentVersion(documentID)
	version3.Version = "3.0.0"
	version3.Content = "第三个版本的内容"

	err := repo.Create(ctx, version1)
	if err != nil {
		t.Fatalf("创建第一个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version2)
	if err != nil {
		t.Fatalf("创建第二个文档版本失败: %v", err)
	}

	err = repo.Create(ctx, version3)
	if err != nil {
		t.Fatalf("创建第三个文档版本失败: %v", err)
	}

	// 验证版本已创建
	versionList, err := repo.GetByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("获取文档版本列表失败: %v", err)
	}
	if len(versionList) != 3 {
		t.Fatalf("预期创建3个版本记录，实际创建 %d 个", len(versionList))
	}

	// 按文档ID删除所有版本
	err = repo.DeleteByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("按文档ID删除版本失败: %v", err)
	}

	// 验证所有版本已被删除
	deletedVersionList, err := repo.GetByDocumentID(ctx, documentID)
	if err != nil {
		t.Fatalf("获取删除后的文档版本列表失败: %v", err)
	}
	if len(deletedVersionList) != 0 {
		t.Errorf("删除后仍有 %d 个版本记录存在，删除失败", len(deletedVersionList))
	}

	// 测试删除不存在文档的版本
	err = repo.DeleteByDocumentID(ctx, "non-existent-doc")
	if err == nil {
		t.Error("删除不存在文档的版本应该返回错误")
	}
}

// TestDocumentVersionRepository_Count 测试统计文档版本数量功能
func TestDocumentVersionRepository_Count(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewDocumentVersionRepository(db)

	ctx := context.Background()
	documentID := "test-doc-id"

	// 初始时应该没有版本记录
	count, err := repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档版本数量失败: %v", err)
	}
	if count != 0 {
		t.Errorf("初始状态下版本数量应该为0，实际为 %d", count)
	}

	// 创建第一个版本记录
	version1 := CreateTestDocumentVersion(documentID)
	version1.Version = "1.0.0"
	version1.Content = "第一个版本的内容"

	err = repo.Create(ctx, version1)
	if err != nil {
		t.Fatalf("创建第一个文档版本失败: %v", err)
	}

	// 验证数量为1
	count, err = repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档版本数量失败: %v", err)
	}
	if count != 1 {
		t.Errorf("创建第一个版本后数量应该为1，实际为 %d", count)
	}

	// 创建第二个版本记录
	version2 := CreateTestDocumentVersion(documentID)
	version2.Version = "2.0.0"
	version2.Content = "第二个版本的内容"

	err = repo.Create(ctx, version2)
	if err != nil {
		t.Fatalf("创建第二个文档版本失败: %v", err)
	}

	// 验证数量为2
	count, err = repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档版本数量失败: %v", err)
	}
	if count != 2 {
		t.Errorf("创建第二个版本后数量应该为2，实际为 %d", count)
	}

	// 创建第三个版本记录
	version3 := CreateTestDocumentVersion(documentID)
	version3.Version = "3.0.0"
	version3.Content = "第三个版本的内容"

	err = repo.Create(ctx, version3)
	if err != nil {
		t.Fatalf("创建第三个文档版本失败: %v", err)
	}

	// 验证数量为3
	count, err = repo.Count(ctx, documentID)
	if err != nil {
		t.Fatalf("统计文档版本数量失败: %v", err)
	}
	if count != 3 {
		t.Errorf("创建第三个版本后数量应该为3，实际为 %d", count)
	}

	// 测试统计不存在文档的版本数量
	count, err = repo.Count(ctx, "non-existent-doc")
	if err != nil {
		t.Fatalf("统计不存在文档的版本数量失败: %v", err)
	}
	if count != 0 {
		t.Errorf("不存在文档的版本数量应该为0，实际为 %d", count)
	}
}
