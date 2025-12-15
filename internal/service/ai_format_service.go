package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

// AIFriendlyFormatService AI友好格式服务接口
type AIFriendlyFormatService interface {
	// StructuredContent 结构化文档内容
	StructuredContent(ctx context.Context, documentID, version, content string, docType model.DocumentType) (*model.AIStructuredContent, error)

	// GenerateLLMFormat 生成LLM优化格式
	GenerateLLMFormat(ctx context.Context, structuredContent *model.AIStructuredContent, options *model.LLMFormatOptions) (*model.LLMOptimizedContent, error)

	// GenerateMultiGranularityRepresentation 生成多粒度文档表示
	GenerateMultiGranularityRepresentation(ctx context.Context, structuredContent *model.AIStructuredContent) (*model.MultiGranularityRepresentation, error)

	// InjectContext 注入上下文
	InjectContext(ctx context.Context, documentID, version, query string, options *model.ContextInjectionOptions) (*model.ContextInjectionResult, error)
}

// aiFriendlyFormatService AI友好格式服务实现
type aiFriendlyFormatService struct {
	documentService DocumentService
}

// NewAIFriendlyFormatService 创建AI友好格式服务实例
func NewAIFriendlyFormatService(documentService DocumentService) AIFriendlyFormatService {
	return &aiFriendlyFormatService{
		documentService: documentService,
	}
}

// StructuredContent 结构化文档内容
func (s *aiFriendlyFormatService) StructuredContent(ctx context.Context, documentID, version, content string, docType model.DocumentType) (*model.AIStructuredContent, error) {
	// 根据文档类型选择结构化策略
	structuredContent := &model.AIStructuredContent{
		DocumentID: documentID,
		Version:    version,
		DocType:    docType,
		CreatedAt:  time.Now(),
	}

	// 内容分段
	segments := s.segmentContent(content)
	structuredContent.Segments = segments

	// 代码示例提取
	codeExamples := s.extractCodeExamples(content, docType)
	structuredContent.CodeExamples = codeExamples

	// 语义标注
	annotations := s.generateSemanticAnnotations(content, docType)
	structuredContent.Annotations = annotations

	// 建立关联关系
	structuredContent.Relations = s.buildRelations(segments, codeExamples, annotations)

	return structuredContent, nil
}

// GenerateLLMFormat 生成LLM优化格式
func (s *aiFriendlyFormatService) GenerateLLMFormat(ctx context.Context, structuredContent *model.AIStructuredContent, options *model.LLMFormatOptions) (*model.LLMOptimizedContent, error) {
	if options == nil {
		options = &model.LLMFormatOptions{
			MaxTokens:       4000,
			PreserveCode:    true,
			SummaryLevel:    model.SummaryLevelMedium,
			IncludeMetadata: true,
		}
	}

	optimizedContent := &model.LLMOptimizedContent{
		DocumentID: structuredContent.DocumentID,
		Version:    structuredContent.Version,
		Options:    *options,
		CreatedAt:  time.Now(),
	}

	// 根据上下文限制优化内容
	optimizedContent.Header = s.generateHeader(structuredContent, options)
	optimizedContent.Content = s.optimizeContent(structuredContent, options)
	optimizedContent.Footer = s.generateFooter(structuredContent, options)
	optimizedContent.Metadata = s.generateMetadata(structuredContent, options)

	return optimizedContent, nil
}

// GenerateMultiGranularityRepresentation 生成多粒度文档表示
func (s *aiFriendlyFormatService) GenerateMultiGranularityRepresentation(ctx context.Context, structuredContent *model.AIStructuredContent) (*model.MultiGranularityRepresentation, error) {
	representation := &model.MultiGranularityRepresentation{
		DocumentID: structuredContent.DocumentID,
		Version:    structuredContent.Version,
		CreatedAt:  time.Now(),
	}

	// 概览表示
	representation.Overview = s.generateOverview(structuredContent)

	// 章节表示
	representation.Sections = s.generateSections(structuredContent)

	// 段落表示
	representation.Paragraphs = s.generateParagraphs(structuredContent)

	// 代码片段表示
	representation.CodeSnippets = s.generateCodeSnippets(structuredContent)

	return representation, nil
}

// InjectContext 注入上下文
func (s *aiFriendlyFormatService) InjectContext(ctx context.Context, documentID, version, query string, options *model.ContextInjectionOptions) (*model.ContextInjectionResult, error) {
	if options == nil {
		options = &model.ContextInjectionOptions{
			MaxContextSize: 3000,
			IncludeCode:    true,
			PriorityLevel:  model.PriorityMedium,
			Format:         model.FormatMarkdown,
		}
	}

	result := &model.ContextInjectionResult{
		DocumentID: documentID,
		Version:    version,
		Query:      query,
		Options:    *options,
		CreatedAt:  time.Now(),
	}

	// 获取文档内容
	documentVersion, err := s.documentService.GetDocumentByVersion(ctx, documentID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get document version: %v", err)
	}

	// 获取文档信息以获取类型
	document, err := s.documentService.GetDocument(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %v", err)
	}

	// 结构化文档内容
	structuredContent, err := s.StructuredContent(ctx, documentID, version, documentVersion.Content, document.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to structure content: %v", err)
	}

	// 根据查询和选项选择相关内容
	result.SelectedContent = s.selectRelevantContent(structuredContent, query, options)

	// 格式化选中的内容
	result.FormattedContext = s.formatContext(result.SelectedContent, options)

	return result, nil
}

// segmentContent 内容分段
func (s *aiFriendlyFormatService) segmentContent(content string) []*model.ContentSegment {
	var segments []*model.ContentSegment

	// 按标题分段
	lines := strings.Split(content, "\n")
	var currentSegment strings.Builder
	var currentTitle string
	var segmentIndex int

	for _, line := range lines {
		// 检测标题
		if isTitleLine(line) {
			// 保存当前段落
			if currentSegment.Len() > 0 {
				segments = append(segments, &model.ContentSegment{
					ID:       fmt.Sprintf("seg_%d", segmentIndex),
					Title:    currentTitle,
					Content:  currentSegment.String(),
					Type:     s.detectSegmentType(currentTitle),
					Position: segmentIndex,
				})
				segmentIndex++
			}

			// 开始新段落
			currentTitle = strings.TrimSpace(line)
			currentSegment.Reset()
			currentSegment.WriteString(line + "\n")
		} else {
			currentSegment.WriteString(line + "\n")
		}
	}

	// 保存最后一个段落
	if currentSegment.Len() > 0 {
		segments = append(segments, &model.ContentSegment{
			ID:       fmt.Sprintf("seg_%d", segmentIndex),
			Title:    currentTitle,
			Content:  currentSegment.String(),
			Type:     s.detectSegmentType(currentTitle),
			Position: segmentIndex,
		})
	}

	return segments
}

// extractCodeExamples 提取代码示例
func (s *aiFriendlyFormatService) extractCodeExamples(content string, docType model.DocumentType) []*model.CodeExample {
	var examples []*model.CodeExample

	// 提取代码块
	codeBlockRegex := regexp.MustCompile("```(\\w+)?\\n([\\s\\S]*?)```")
	matches := codeBlockRegex.FindAllStringSubmatch(content, -1)

	for i, match := range matches {
		language := "unknown"
		if len(match[1]) > 0 {
			language = match[1]
		}

		examples = append(examples, &model.CodeExample{
			ID:          fmt.Sprintf("code_%d", i),
			Language:    language,
			Code:        match[2],
			Description: s.generateCodeDescription(match[2], language),
			Position:    i,
		})
	}

	// 提取内联代码
	inlineCodeRegex := regexp.MustCompile("`([^`]+)`")
	inlineMatches := inlineCodeRegex.FindAllStringSubmatch(content, -1)

	for i, match := range inlineMatches {
		examples = append(examples, &model.CodeExample{
			ID:          fmt.Sprintf("inline_%d", i),
			Language:    "inline",
			Code:        match[1],
			Description: "Inline code example",
			Position:    i,
			IsInline:    true,
		})
	}

	return examples
}

// generateSemanticAnnotations 生成语义标注
func (s *aiFriendlyFormatService) generateSemanticAnnotations(content string, docType model.DocumentType) []*model.SemanticAnnotation {
	var annotations []*model.SemanticAnnotation

	// 定义概念
	conceptRegex := regexp.MustCompile(`\b(定义|概念|解释|说明|含义)\b[:：]\s*([^\n]+)`)
	matches := conceptRegex.FindAllStringSubmatch(content, -1)

	for i, match := range matches {
		annotations = append(annotations, &model.SemanticAnnotation{
			ID:       fmt.Sprintf("concept_%d", i),
			Type:     "concept",
			Value:    match[2],
			Context:  match[0],
			Position: i,
		})
	}

	// 定义步骤
	stepRegex := regexp.MustCompile(`\b(步骤|流程|操作|方法)\b[:：]\s*([^\n]+)`)
	stepMatches := stepRegex.FindAllStringSubmatch(content, -1)

	for i, match := range stepMatches {
		annotations = append(annotations, &model.SemanticAnnotation{
			ID:       fmt.Sprintf("step_%d", i),
			Type:     "step",
			Value:    match[2],
			Context:  match[0],
			Position: i,
		})
	}

	// 识别API端点
	if docType == model.DocumentTypeSwagger || docType == model.DocumentTypeOpenAPI {
		apiRegex := regexp.MustCompile(`"path":\s*"([^"]+)"`)
		apiMatches := apiRegex.FindAllStringSubmatch(content, -1)

		for i, match := range apiMatches {
			annotations = append(annotations, &model.SemanticAnnotation{
				ID:       fmt.Sprintf("api_%d", i),
				Type:     "api_endpoint",
				Value:    match[1],
				Context:  match[0],
				Position: i,
			})
		}
	}

	return annotations
}

// buildRelations 建立关联关系
func (s *aiFriendlyFormatService) buildRelations(segments []*model.ContentSegment, codeExamples []*model.CodeExample, annotations []*model.SemanticAnnotation) []*model.ContentRelation {
	var relations []*model.ContentRelation

	// 建立代码示例与段落的关联
	for _, code := range codeExamples {
		for _, segment := range segments {
			if strings.Contains(segment.Content, code.Code) {
				relations = append(relations, &model.ContentRelation{
					SourceID:   segment.ID,
					TargetID:   code.ID,
					Type:       "contains",
					Confidence: 0.9,
				})
			}
		}
	}

	// 建立语义标注与段落的关联
	for _, annotation := range annotations {
		for _, segment := range segments {
			if strings.Contains(segment.Content, annotation.Context) {
				relations = append(relations, &model.ContentRelation{
					SourceID:   segment.ID,
					TargetID:   annotation.ID,
					Type:       "annotates",
					Confidence: 0.8,
				})
			}
		}
	}

	return relations
}

// isTitleLine 判断是否为标题行
func isTitleLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Markdown标题
	if strings.HasPrefix(trimmed, "#") {
		return true
	}

	// 其他可能的标题模式
	return false
}

// detectSegmentType 检测段落类型
func (s *aiFriendlyFormatService) detectSegmentType(title string) string {
	lowerTitle := strings.ToLower(title)

	if strings.Contains(lowerTitle, "简介") || strings.Contains(lowerTitle, "介绍") || strings.Contains(lowerTitle, "概述") {
		return "introduction"
	}

	if strings.Contains(lowerTitle, "示例") || strings.Contains(lowerTitle, "例子") {
		return "example"
	}

	if strings.Contains(lowerTitle, "步骤") || strings.Contains(lowerTitle, "流程") {
		return "procedure"
	}

	if strings.Contains(lowerTitle, "结论") || strings.Contains(lowerTitle, "总结") {
		return "conclusion"
	}

	return "content"
}

// generateCodeDescription 生成代码描述
func (s *aiFriendlyFormatService) generateCodeDescription(code, language string) string {
	// 简单的代码描述生成逻辑
	lines := strings.Split(strings.TrimSpace(code), "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])

		// 如果第一行是注释，提取作为描述
		if strings.HasPrefix(firstLine, "//") || strings.HasPrefix(firstLine, "#") {
			return strings.TrimSpace(firstLine[2:])
		}

		// 否则使用语言和代码长度生成简单描述
		return fmt.Sprintf("%s代码示例 (%d行)", language, len(lines))
	}

	return fmt.Sprintf("%s代码示例", language)
}

// generateHeader 生成LLM优化格式的头部
func (s *aiFriendlyFormatService) generateHeader(structuredContent *model.AIStructuredContent, options *model.LLMFormatOptions) string {
	var header strings.Builder

	header.WriteString(fmt.Sprintf("# 文档：%s (版本: %s)\n\n", structuredContent.DocumentID, structuredContent.Version))

	if options.IncludeMetadata {
		header.WriteString("## 文档信息\n")
		header.WriteString(fmt.Sprintf("- 文档类型: %s\n", structuredContent.DocType))
		header.WriteString(fmt.Sprintf("- 段落数量: %d\n", len(structuredContent.Segments)))
		header.WriteString(fmt.Sprintf("- 代码示例数量: %d\n", len(structuredContent.CodeExamples)))
		header.WriteString(fmt.Sprintf("- 语义标注数量: %d\n\n", len(structuredContent.Annotations)))
	}

	return header.String()
}

// optimizeContent 优化内容
func (s *aiFriendlyFormatService) optimizeContent(structuredContent *model.AIStructuredContent, options *model.LLMFormatOptions) string {
	var content strings.Builder
	currentTokens := 0
	headerTokens := s.estimateTokens(s.generateHeader(structuredContent, options))
	footerTokens := s.estimateTokens(s.generateFooter(structuredContent, options))
	availableTokens := options.MaxTokens - headerTokens - footerTokens - 100 // 预留100个token

	// 根据摘要级别调整内容
	switch options.SummaryLevel {
	case model.SummaryLevelBrief:
		// 只包含标题和关键段落
		for _, segment := range structuredContent.Segments {
			if segment.Type == "introduction" || segment.Type == "conclusion" {
				segmentTokens := s.estimateTokens(segment.Content)
				if currentTokens+segmentTokens <= availableTokens {
					content.WriteString(fmt.Sprintf("## %s\n\n", segment.Title))
					content.WriteString(segment.Content)
					content.WriteString("\n\n")
					currentTokens += segmentTokens
				}
			}
		}

	case model.SummaryLevelMedium:
		// 包含所有段落标题和前几个段落
		for i, segment := range structuredContent.Segments {
			segmentTokens := s.estimateTokens(segment.Content)
			if currentTokens+segmentTokens <= availableTokens {
				content.WriteString(fmt.Sprintf("## %s\n\n", segment.Title))
				if i < 3 { // 前三个段落包含完整内容
					content.WriteString(segment.Content)
				} else { // 后续段落只包含摘要
					content.WriteString(s.summarizeContent(segment.Content, 100))
				}
				content.WriteString("\n\n")
				currentTokens += segmentTokens
			}
		}

	case model.SummaryLevelDetailed:
		// 包含所有内容
		for _, segment := range structuredContent.Segments {
			segmentTokens := s.estimateTokens(segment.Content)
			if currentTokens+segmentTokens <= availableTokens {
				content.WriteString(fmt.Sprintf("## %s\n\n", segment.Title))
				content.WriteString(segment.Content)
				content.WriteString("\n\n")
				currentTokens += segmentTokens
			}
		}
	}

	// 添加代码示例（如果启用）
	if options.PreserveCode && len(structuredContent.CodeExamples) > 0 {
		content.WriteString("## 代码示例\n\n")
		for _, example := range structuredContent.CodeExamples {
			exampleTokens := s.estimateTokens(example.Code)
			if currentTokens+exampleTokens+50 <= availableTokens { // 预留50个token用于格式
				content.WriteString(fmt.Sprintf("### %s\n\n", example.Description))
				content.WriteString(fmt.Sprintf("```%s\n", example.Language))
				content.WriteString(example.Code)
				content.WriteString("\n```\n\n")
				currentTokens += exampleTokens + 50
			}
		}
	}

	return content.String()
}

// generateFooter 生成LLM优化格式的尾部
func (s *aiFriendlyFormatService) generateFooter(structuredContent *model.AIStructuredContent, options *model.LLMFormatOptions) string {
	var footer strings.Builder

	if options.IncludeMetadata {
		footer.WriteString("## 关键概念\n\n")
		for _, annotation := range structuredContent.Annotations {
			if annotation.Type == "concept" {
				footer.WriteString(fmt.Sprintf("- %s: %s\n", annotation.Type, annotation.Value))
			}
		}
		footer.WriteString("\n")
	}

	footer.WriteString("---\n")
	footer.WriteString(fmt.Sprintf("本文档由AI友好格式服务生成 (时间: %s)\n", time.Now().Format("2006-01-02 15:04:05")))

	return footer.String()
}

// generateMetadata 生成元数据
func (s *aiFriendlyFormatService) generateMetadata(structuredContent *model.AIStructuredContent, options *model.LLMFormatOptions) map[string]interface{} {
	metadata := make(map[string]interface{})

	metadata["document_id"] = structuredContent.DocumentID
	metadata["version"] = structuredContent.Version
	metadata["doc_type"] = structuredContent.DocType
	metadata["segment_count"] = len(structuredContent.Segments)
	metadata["code_example_count"] = len(structuredContent.CodeExamples)
	metadata["annotation_count"] = len(structuredContent.Annotations)
	metadata["relation_count"] = len(structuredContent.Relations)
	metadata["options"] = options
	metadata["generated_at"] = time.Now()

	return metadata
}

// generateOverview 生成概览表示
func (s *aiFriendlyFormatService) generateOverview(structuredContent *model.AIStructuredContent) string {
	var overview strings.Builder

	overview.WriteString(fmt.Sprintf("# 文档概览: %s\n\n", structuredContent.DocumentID))

	// 添加标题和简介
	for _, segment := range structuredContent.Segments {
		if segment.Type == "introduction" {
			overview.WriteString(fmt.Sprintf("## %s\n\n", segment.Title))
			overview.WriteString(s.summarizeContent(segment.Content, 200))
			overview.WriteString("\n\n")
		}
	}

	// 添加章节列表
	overview.WriteString("## 章节列表\n\n")
	for _, segment := range structuredContent.Segments {
		overview.WriteString(fmt.Sprintf("- %s\n", segment.Title))
	}

	return overview.String()
}

// generateSections 生成章节表示
func (s *aiFriendlyFormatService) generateSections(structuredContent *model.AIStructuredContent) []*model.SectionRepresentation {
	var sections []*model.SectionRepresentation

	for _, segment := range structuredContent.Segments {
		section := &model.SectionRepresentation{
			ID:       segment.ID,
			Title:    segment.Title,
			Type:     segment.Type,
			Content:  s.summarizeContent(segment.Content, 300),
			Position: segment.Position,
		}

		// 查找相关代码示例
		var relatedCodes []string
		for _, relation := range structuredContent.Relations {
			if relation.SourceID == segment.ID && relation.Type == "contains" {
				relatedCodes = append(relatedCodes, relation.TargetID)
			}
		}
		section.RelatedCodeExamples = relatedCodes

		// 查找相关语义标注
		var relatedAnnotations []string
		for _, relation := range structuredContent.Relations {
			if relation.SourceID == segment.ID && relation.Type == "annotates" {
				relatedAnnotations = append(relatedAnnotations, relation.TargetID)
			}
		}
		section.RelatedAnnotations = relatedAnnotations

		sections = append(sections, section)
	}

	return sections
}

// generateParagraphs 生成段落表示
func (s *aiFriendlyFormatService) generateParagraphs(structuredContent *model.AIStructuredContent) []*model.ParagraphRepresentation {
	var paragraphs []*model.ParagraphRepresentation

	// 将每个段落拆分为更小的单元
	paragraphIndex := 0
	for _, segment := range structuredContent.Segments {
		// 按句子拆分段落
		sentences := s.splitIntoSentences(segment.Content)

		for i := 0; i < len(sentences); i += 3 { // 每个段落包含3个句子
			end := i + 3
			if end > len(sentences) {
				end = len(sentences)
			}

			if i < end {
				paragraph := &model.ParagraphRepresentation{
					ID:          fmt.Sprintf("para_%d", paragraphIndex),
					Content:     strings.Join(sentences[i:end], " "),
					SegmentID:   segment.ID,
					SegmentType: segment.Type,
					Position:    paragraphIndex,
				}
				paragraphs = append(paragraphs, paragraph)
				paragraphIndex++
			}
		}
	}

	return paragraphs
}

// generateCodeSnippets 生成代码片段表示
func (s *aiFriendlyFormatService) generateCodeSnippets(structuredContent *model.AIStructuredContent) []*model.CodeSnippetRepresentation {
	var snippets []*model.CodeSnippetRepresentation

	for _, example := range structuredContent.CodeExamples {
		snippet := &model.CodeSnippetRepresentation{
			ID:          example.ID,
			Language:    example.Language,
			Code:        example.Code,
			Description: example.Description,
			Position:    example.Position,
			IsInline:    example.IsInline,
		}

		// 查找相关段落
		var relatedSegments []string
		for _, relation := range structuredContent.Relations {
			if relation.TargetID == example.ID && relation.Type == "contains" {
				relatedSegments = append(relatedSegments, relation.SourceID)
			}
		}
		snippet.RelatedSegments = relatedSegments

		snippets = append(snippets, snippet)
	}

	return snippets
}

// selectRelevantContent 选择相关内容
func (s *aiFriendlyFormatService) selectRelevantContent(structuredContent *model.AIStructuredContent, query string, options *model.ContextInjectionOptions) []*model.ContentSelection {
	var selections []*model.ContentSelection

	// 简单的关键词匹配逻辑
	queryTerms := strings.Fields(strings.ToLower(query))

	// 评分每个段落
	for _, segment := range structuredContent.Segments {
		score := s.calculateRelevanceScore(segment.Content, queryTerms)
		if score > 0.1 { // 只选择相关性大于0.1的段落
			selection := &model.ContentSelection{
				ID:       segment.ID,
				Type:     "segment",
				Content:  segment.Content,
				Title:    segment.Title,
				Score:    score,
				Position: segment.Position,
			}
			selections = append(selections, selection)
		}
	}

	// 评分每个代码示例
	if options.IncludeCode {
		for _, example := range structuredContent.CodeExamples {
			score := s.calculateRelevanceScore(example.Code, queryTerms)
			if score > 0.1 { // 只选择相关性大于0.1的代码示例
				selection := &model.ContentSelection{
					ID:       example.ID,
					Type:     "code",
					Content:  example.Code,
					Title:    example.Description,
					Score:    score,
					Position: example.Position,
				}
				selections = append(selections, selection)
			}
		}
	}

	// 按分数排序并选择前N个
	selections = s.sortAndLimitSelections(selections, options.MaxContextSize)

	return selections
}

// formatContext 格式化上下文
func (s *aiFriendlyFormatService) formatContext(selections []*model.ContentSelection, options *model.ContextInjectionOptions) string {
	var context strings.Builder

	switch options.Format {
	case model.FormatMarkdown:
		context.WriteString("# 相关文档内容\n\n")
		for _, selection := range selections {
			if selection.Type == "segment" {
				context.WriteString(fmt.Sprintf("## %s\n\n", selection.Title))
				context.WriteString(selection.Content)
				context.WriteString("\n\n")
			} else if selection.Type == "code" {
				context.WriteString(fmt.Sprintf("## %s\n\n", selection.Title))
				// 这里需要从代码示例中获取语言信息，暂时使用默认值
				context.WriteString("```text\n")
				context.WriteString(selection.Content)
				context.WriteString("\n```\n\n")
			}
		}

	case model.FormatJSON:
		jsonData := make(map[string]interface{})
		jsonData["query"] = selections
		jsonBytes, _ := json.MarshalIndent(jsonData, "", "  ")
		context.WriteString(string(jsonBytes))

	case model.FormatPlainText:
		for _, selection := range selections {
			context.WriteString(fmt.Sprintf("%s: %s\n\n", selection.Title, selection.Content))
		}
	}

	return context.String()
}

// summarizeContent 内容摘要
func (s *aiFriendlyFormatService) summarizeContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}

	// 简单截断，实际应用中可以使用更智能的摘要算法
	summary := content[:maxChars]

	// 尝试在句子边界截断
	lastPeriod := strings.LastIndex(summary, "。")
	if lastPeriod > maxChars/2 {
		summary = summary[:lastPeriod+1]
	}

	return summary + "..."
}

// splitIntoSentences 拆分为句子
func (s *aiFriendlyFormatService) splitIntoSentences(content string) []string {
	// 简单的句子拆分逻辑
	sentences := regexp.MustCompile(`[。！？.!?]\s*`).Split(content, -1)
	var result []string

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			result = append(result, sentence)
		}
	}

	return result
}

// estimateTokens 估算token数量
func (s *aiFriendlyFormatService) estimateTokens(text string) int {
	// 简单的token估算，实际应用中可以使用更精确的方法
	// 假设每个汉字或每个单词对应1.3个token
	chineseCount := regexp.MustCompile(`[\p{Han}]`).FindAllStringIndex(text, -1)
	wordCount := len(strings.Fields(text))

	return int(float64(len(chineseCount)+wordCount) * 1.3)
}

// calculateRelevanceScore 计算相关性分数
func (s *aiFriendlyFormatService) calculateRelevanceScore(content string, queryTerms []string) float64 {
	if len(queryTerms) == 0 {
		return 0
	}

	lowerContent := strings.ToLower(content)
	matches := 0

	for _, term := range queryTerms {
		if strings.Contains(lowerContent, term) {
			matches++
		}
	}

	return float64(matches) / float64(len(queryTerms))
}

// sortAndLimitSelections 排序并限制选择数量
func (s *aiFriendlyFormatService) sortAndLimitSelections(selections []*model.ContentSelection, maxContextSize int) []*model.ContentSelection {
	// 按分数降序排序
	for i := 0; i < len(selections)-1; i++ {
		for j := i + 1; j < len(selections); j++ {
			if selections[i].Score < selections[j].Score {
				selections[i], selections[j] = selections[j], selections[i]
			}
		}
	}

	// 限制选择数量，考虑上下文大小限制
	totalTokens := 0
	var result []*model.ContentSelection

	for _, selection := range selections {
		selectionTokens := s.estimateTokens(selection.Content)
		if totalTokens+selectionTokens <= maxContextSize {
			result = append(result, selection)
			totalTokens += selectionTokens
		} else {
			break
		}
	}

	return result
}
