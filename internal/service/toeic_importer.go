package service

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// toeicImporterService TOEIC文件导入服务实现
type toeicImporterService struct {
	testRepo              repository.TestRepository
	scenarioRepo          repository.ScenarioRepository
	difficultyLevelRepo   repository.DifficultyLevelRepository
	part2QuestionRepo     repository.Part2QuestionRepository
	part3ConversationRepo repository.Part3ConversationRepository
	part3AnswerOptionRepo repository.Part3AnswerOptionRepository
	part4TalkRepo         repository.Part4TalkRepository
	part4AnswerOptionRepo repository.Part4AnswerOptionRepository
	logger                logger.Logger

	// 缓存映射
	scenarioMap   map[string]uint
	difficultyMap map[string]uint
}

// NewToeicImporterService 创建TOEIC导入服务
func NewToeicImporterService(
	testRepo repository.TestRepository,
	scenarioRepo repository.ScenarioRepository,
	difficultyLevelRepo repository.DifficultyLevelRepository,
	part2QuestionRepo repository.Part2QuestionRepository,
	part3ConversationRepo repository.Part3ConversationRepository,
	part3AnswerOptionRepo repository.Part3AnswerOptionRepository,
	part4TalkRepo repository.Part4TalkRepository,
	part4AnswerOptionRepo repository.Part4AnswerOptionRepository,
	logger logger.Logger,
) ToeicImporterService {
	return &toeicImporterService{
		testRepo:              testRepo,
		scenarioRepo:          scenarioRepo,
		difficultyLevelRepo:   difficultyLevelRepo,
		part2QuestionRepo:     part2QuestionRepo,
		part3ConversationRepo: part3ConversationRepo,
		part3AnswerOptionRepo: part3AnswerOptionRepo,
		part4TalkRepo:         part4TalkRepo,
		part4AnswerOptionRepo: part4AnswerOptionRepo,
		logger:                logger,
		scenarioMap:           make(map[string]uint),
		difficultyMap:         make(map[string]uint),
	}
}

// ImportResult 导入结果
type ImportResult struct {
	TestName           string   `json:"test_name"`
	Part2Questions     int      `json:"part2_questions"`
	Part3Conversations int      `json:"part3_conversations"`
	Part4Talks         int      `json:"part4_talks"`
	Errors             []string `json:"errors"`
	Warnings           []string `json:"warnings"`
}

// ParsedQuestion 解析的问题
type ParsedQuestion struct {
	Number        int
	Text          string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectAnswer string
}

// ParsedPart2Question 解析的Part2问题
type ParsedPart2Question struct {
	Number            int
	Text              string
	OptionA           string
	OptionB           string
	OptionC           string
	CorrectAnswer     string
	ScenarioID        *uint
	DifficultyLevelID *uint
}

// ParsedConversation 解析的对话
type ParsedConversation struct {
	Number            int
	Title             string
	Content           string
	ScenarioID        *uint
	DifficultyLevelID *uint
	Questions         []ParsedQuestion
}

// ParsedTalk 解析的演讲
type ParsedTalk struct {
	Number            int
	Title             string
	Content           string
	ScenarioID        *uint
	DifficultyLevelID *uint
	Questions         []ParsedQuestion
}

// ValidateFile 验证文件是否为.md格式的文件
func (s *toeicImporterService) ValidateFile(filename string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// 检查是否为.md文件
	if !strings.HasSuffix(strings.ToLower(filename), ".md") {
		return fmt.Errorf("file must have .md extension: %s", filename)
	}

	return nil
}

// ImportFile 导入TOEIC测试文件
func (s *toeicImporterService) ImportFile(ctx context.Context, filename string) (*ImportResult, error) {
	s.logger.Info("Starting TOEIC file import", "filename", filename)

	// 验证文件
	if err := s.ValidateFile(filename); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// 初始化映射
	if err := s.initializeMappings(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize mappings: %w", err)
	}

	// 提取测试名称
	testName := s.extractTestName(filename)

	// 获取或创建测试
	test, err := s.getOrCreateTest(ctx, testName)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create test: %w", err)
	}

	// 读取文件内容
	lines, err := s.readFileLines(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result := &ImportResult{
		TestName: testName,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 解析并导入Part 2问题
	part2Questions, err := s.parsePart2(lines)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Part 2 parsing error: %v", err))
	} else {
		if err := s.importPart2Questions(ctx, test.ID, part2Questions); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Part 2 import error: %v", err))
		} else {
			result.Part2Questions = len(part2Questions)
		}
	}

	// 解析并导入Part 3对话
	part3Conversations, err := s.parsePart3(lines)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Part 3 parsing error: %v", err))
	} else {
		if err := s.importPart3Conversations(ctx, test.ID, part3Conversations); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Part 3 import error: %v", err))
		} else {
			result.Part3Conversations = len(part3Conversations)
		}
	}

	// 解析并导入Part 4演讲
	part4Talks, err := s.parsePart4(lines)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Part 4 parsing error: %v", err))
	} else {
		if err := s.importPart4Talks(ctx, test.ID, part4Talks); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Part 4 import error: %v", err))
		} else {
			result.Part4Talks = len(part4Talks)
		}
	}

	if len(result.Errors) == 0 {
		s.logger.Info("TOEIC file import completed successfully", "filename", filename)
	} else {
		s.logger.Warn("TOEIC file import completed with errors", "filename", filename, "errors", len(result.Errors))
	}

	return result, nil
}

// initializeMappings 初始化场景和难度级别映射
func (s *toeicImporterService) initializeMappings(ctx context.Context) error {
	// 加载场景映射
	scenarios, _, err := s.scenarioRepo.List(ctx, repository.ListOptions{PageSize: 1000})
	if err != nil {
		return fmt.Errorf("failed to load scenarios: %w", err)
	}

	for _, scenario := range scenarios {
		s.scenarioMap[scenario.Name] = scenario.ID
	}

	// 加载难度级别映射
	difficultyLevels, _, err := s.difficultyLevelRepo.List(ctx, repository.ListOptions{PageSize: 1000})
	if err != nil {
		return fmt.Errorf("failed to load difficulty levels: %w", err)
	}

	for _, level := range difficultyLevels {
		s.difficultyMap[level.Name] = level.ID
	}

	s.logger.Info("Initialized mappings", "scenarios", len(s.scenarioMap), "difficulty_levels", len(s.difficultyMap))
	return nil
}

// extractTestName 从文件名提取测试名称
func (s *toeicImporterService) extractTestName(filename string) string {
	parts := strings.Split(filename, "/")
	name := parts[len(parts)-1]
	if strings.HasSuffix(name, ".md") {
		name = name[:len(name)-3]
	}
	return name
}

// getOrCreateTest 获取或创建测试
func (s *toeicImporterService) getOrCreateTest(ctx context.Context, testName string) (*model.Test, error) {
	// 尝试获取现有测试
	test, err := s.testRepo.GetByName(ctx, testName)
	if err == nil {
		return test, nil
	}

	// 创建新测试
	newTest := &model.Test{
		Name: testName,
		Description: sql.NullString{
			String: fmt.Sprintf("TOEIC Test: %s", testName),
			Valid:  true,
		},
	}

	if err := s.testRepo.Create(ctx, newTest); err != nil {
		return nil, fmt.Errorf("failed to create test: %w", err)
	}

	s.logger.Info("Created new test", "name", testName, "id", newTest.ID)
	return newTest, nil
}

// readFileLines 读取文件所有行
func (s *toeicImporterService) readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// CleanupTestData 清理测试数据
func (s *toeicImporterService) CleanupTestData(ctx context.Context, testName string) error {
	test, err := s.testRepo.GetByName(ctx, testName)
	if err != nil {
		s.logger.Info("Test not found for cleanup", "testName", testName)
		return nil // 测试不存在，无需清理
	}

	s.logger.Info("Cleaning up test data", "testName", testName, "testID", test.ID)

	// 获取Part3对话ID列表，用于删除相关的答案选项
	part3Conversations, _, err := s.part3ConversationRepo.List(ctx, repository.ListOptions{
		Filters:  map[string]interface{}{"test_id": int32(test.ID)},
		PageSize: 1000,
	})
	if err != nil {
		s.logger.Warn("Failed to get Part3 conversations for cleanup", "error", err)
	} else if len(part3Conversations) > 0 {
		var conversationIDs []uint
		for _, conv := range part3Conversations {
			conversationIDs = append(conversationIDs, conv.ID)
		}

		// 删除Part3答案选项
		if err := s.part3AnswerOptionRepo.DeleteByConversationIDs(ctx, conversationIDs); err != nil {
			s.logger.Error("Failed to delete Part3 answer options", "error", err)
			return fmt.Errorf("failed to delete Part3 answer options: %w", err)
		}
		s.logger.Info("Deleted Part3 answer options", "count", len(conversationIDs))
	}

	// 获取Part4演讲ID列表，用于删除相关的答案选项
	part4Talks, _, err := s.part4TalkRepo.List(ctx, repository.ListOptions{
		Filters:  map[string]interface{}{"test_id": int32(test.ID)},
		PageSize: 1000,
	})
	if err != nil {
		s.logger.Warn("Failed to get Part4 talks for cleanup", "error", err)
	} else if len(part4Talks) > 0 {
		var talkIDs []uint
		for _, talk := range part4Talks {
			talkIDs = append(talkIDs, talk.ID)
		}

		// 删除Part4答案选项
		if err := s.part4AnswerOptionRepo.DeleteByTalkIDs(ctx, talkIDs); err != nil {
			s.logger.Error("Failed to delete Part4 answer options", "error", err)
			return fmt.Errorf("failed to delete Part4 answer options: %w", err)
		}
		s.logger.Info("Deleted Part4 answer options", "count", len(talkIDs))
	}

	// 删除Part2问题
	if err := s.part2QuestionRepo.DeleteByTestID(ctx, test.ID); err != nil {
		s.logger.Error("Failed to delete Part2 questions", "error", err)
		return fmt.Errorf("failed to delete Part2 questions: %w", err)
	}
	s.logger.Info("Deleted Part2 questions for test", "testID", test.ID)

	// 删除Part3对话
	if err := s.part3ConversationRepo.DeleteByTestID(ctx, test.ID); err != nil {
		s.logger.Error("Failed to delete Part3 conversations", "error", err)
		return fmt.Errorf("failed to delete Part3 conversations: %w", err)
	}
	s.logger.Info("Deleted Part3 conversations for test", "testID", test.ID)

	// 删除Part4演讲
	if err := s.part4TalkRepo.DeleteByTestID(ctx, test.ID); err != nil {
		s.logger.Error("Failed to delete Part4 talks", "error", err)
		return fmt.Errorf("failed to delete Part4 talks: %w", err)
	}
	s.logger.Info("Deleted Part4 talks for test", "testID", test.ID)

	s.logger.Info("Successfully cleaned up test data", "testName", testName, "testID", test.ID)
	return nil
}

// parseScenarioAndDifficulty 解析场景和难度级别
func (s *toeicImporterService) parseScenarioAndDifficulty(line string) (*uint, *uint, error) {
	// Pattern 1: *1|Workplace/Business* - **Difficulty: Intermediate**
	// Pattern 2: *1|Workplace/Business* (without difficulty)
	// Pattern 3: [1|Workplace/Business] - **Difficulty: Intermediate**
	// Pattern 4: [1|Workplace/Business] (without difficulty)

	// Try asterisk pattern with difficulty first
	scenarioRegex := regexp.MustCompile(`\*(\d+)\|([^*]+)\*\s*-\s*\*\*Difficulty:\s*([^*]+)\*\*`)
	matches := scenarioRegex.FindStringSubmatch(line)

	if len(matches) == 4 {
		return s.processScenarioMatch(matches[1], matches[2], matches[3])
	}

	// Try square bracket pattern with difficulty
	scenarioRegex = regexp.MustCompile(`\[(\d+)\|([^\]]+)\]\s*-\s*\*\*Difficulty:\s*([^*]+)\*\*`)
	matches = scenarioRegex.FindStringSubmatch(line)

	if len(matches) == 4 {
		return s.processScenarioMatch(matches[1], matches[2], matches[3])
	}

	// Try square bracket pattern with "Difficulty: Level" format
	scenarioRegex = regexp.MustCompile(`\[(\d+)\|([^\]]+)\]\s*-\s*Difficulty:\s*([A-Za-z]+)`)
	matches = scenarioRegex.FindStringSubmatch(line)

	if len(matches) == 4 {
		return s.processScenarioMatch(matches[1], matches[2], matches[3])
	}

	// Try square bracket pattern with simple difficulty
	scenarioRegex = regexp.MustCompile(`\[(\d+)\|([^\]]+)\]\s*-\s*([A-Za-z]+)`)
	matches = scenarioRegex.FindStringSubmatch(line)

	if len(matches) == 4 {
		return s.processScenarioMatch(matches[1], matches[2], matches[3])
	}

	// Try asterisk pattern without difficulty
	scenarioOnlyRegex := regexp.MustCompile(`\*(\d+)\|([^*]+)\*`)
	matches = scenarioOnlyRegex.FindStringSubmatch(line)

	if len(matches) == 3 {
		return s.processScenarioMatch(matches[1], matches[2], "")
	}

	// Try square bracket pattern without difficulty
	scenarioOnlyRegex = regexp.MustCompile(`\[(\d+)\|([^\]]+)\]`)
	matches = scenarioOnlyRegex.FindStringSubmatch(line)

	if len(matches) == 3 {
		return s.processScenarioMatch(matches[1], matches[2], "")
	}

	return nil, nil, fmt.Errorf("failed to parse scenario and difficulty from: %s", line)
}

// processScenarioMatch 处理匹配的场景和难度字符串
func (s *toeicImporterService) processScenarioMatch(scenarioIDStr, scenarioName, difficultyName string) (*uint, *uint, error) {
	// Convert scenario ID
	scenarioID, err := strconv.Atoi(scenarioIDStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid scenario ID: %s", scenarioIDStr)
	}

	// Map scenario name to ID (fallback to parsed ID if name not found)
	var finalScenarioID *uint
	if id, exists := s.scenarioMap[strings.TrimSpace(scenarioName)]; exists {
		finalScenarioID = &id
	} else {
		id := uint(scenarioID)
		finalScenarioID = &id
		s.logger.Warn("Scenario not found in database, using ID", "scenario", scenarioName, "id", scenarioID)
	}

	// Map difficulty name to ID if provided
	var difficultyID *uint
	if difficultyName != "" {
		if id, exists := s.difficultyMap[strings.TrimSpace(difficultyName)]; exists {
			difficultyID = &id
		} else {
			s.logger.Warn("Difficulty level not found in database", "difficulty", difficultyName)
		}
	}

	return finalScenarioID, difficultyID, nil
}
