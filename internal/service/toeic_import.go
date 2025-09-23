package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
)

// importPart2Questions 导入Part 2问题
func (s *toeicImporterService) importPart2Questions(ctx context.Context, testID uint, questions []ParsedPart2Question) error {
	for _, q := range questions {
		// 创建Part2Question模型
		part2Question := &model.Part2Question{
			TestId:         int32(testID),
			QuestionNumber: int32(q.Number),
			QuestionText:   q.Text,
			OptionA:        q.OptionA,
			OptionB:        q.OptionB,
			OptionC:        q.OptionC,
		}

		// 处理正确答案
		if q.CorrectAnswer != "" {
			part2Question.CorrectAnswer = sql.NullString{
				String: q.CorrectAnswer,
				Valid:  true,
			}
		}

		// 处理场景ID
		if q.ScenarioID != nil {
			part2Question.ScenarioId = sql.NullInt32{
				Int32: int32(*q.ScenarioID),
				Valid: true,
			}
		}

		// 处理难度级别ID
		if q.DifficultyLevelID != nil {
			part2Question.DifficultyLevelId = sql.NullInt32{
				Int32: int32(*q.DifficultyLevelID),
				Valid: true,
			}
		}

		// 保存到数据库
		if err := s.part2QuestionRepo.Create(ctx, part2Question); err != nil {
			return fmt.Errorf("failed to create Part 2 question %d: %w", q.Number, err)
		}
	}

	s.logger.Info("Imported Part 2 questions", "count", len(questions))
	return nil
}

// importPart3Conversations 导入Part 3对话
func (s *toeicImporterService) importPart3Conversations(ctx context.Context, testID uint, conversations []ParsedConversation) error {
	for _, conv := range conversations {
		// 创建Part3Conversation模型
		part3Conversation := &model.Part3Conversation{
			TestId:             int32(testID),
			ConversationNumber: int32(conv.Number),
			Content:            conv.Content,
		}

		// 处理标题
		if conv.Title != "" {
			part3Conversation.Title = sql.NullString{
				String: conv.Title,
				Valid:  true,
			}
		}

		// 处理问题文本
		if len(conv.Questions) > 0 {
			part3Conversation.Question1 = s.getQuestionText(conv.Questions, 1)
		}
		if len(conv.Questions) > 1 {
			part3Conversation.Question2 = s.getQuestionText(conv.Questions, 2)
		}
		if len(conv.Questions) > 2 {
			part3Conversation.Question3 = s.getQuestionText(conv.Questions, 3)
		}

		// 处理场景ID
		if conv.ScenarioID != nil {
			part3Conversation.ScenarioId = sql.NullInt32{
				Int32: int32(*conv.ScenarioID),
				Valid: true,
			}
		}

		// 处理难度级别ID
		if conv.DifficultyLevelID != nil {
			part3Conversation.DifficultyLevelId = sql.NullInt32{
				Int32: int32(*conv.DifficultyLevelID),
				Valid: true,
			}
		}

		// 保存对话到数据库
		if err := s.part3ConversationRepo.Create(ctx, part3Conversation); err != nil {
			return fmt.Errorf("failed to create Part 3 conversation %d: %w", conv.Number, err)
		}

		s.logger.Info("Created Part 3 conversation", "number", conv.Number, "id", part3Conversation.ID)

		// 导入答案选项
		for i, question := range conv.Questions {
			if i >= 3 { // 只处理前3个问题
				break
			}

			part3AnswerOption := &model.Part3AnswerOption{
				ConversationId: int32(part3Conversation.ID),
				QuestionNumber: int32(i + 1),
				OptionA:        question.OptionA,
				OptionB:        question.OptionB,
				OptionC:        question.OptionC,
				OptionD:        question.OptionD,
			}

			// 处理正确答案
			if question.CorrectAnswer != "" {
				part3AnswerOption.CorrectAnswer = sql.NullString{
					String: question.CorrectAnswer,
					Valid:  true,
				}
			}

			if err := s.part3AnswerOptionRepo.Create(ctx, part3AnswerOption); err != nil {
				return fmt.Errorf("failed to create Part 3 answer option for conversation %d, question %d: %w", conv.Number, i+1, err)
			}
		}
	}

	s.logger.Info("Imported Part 3 conversations", "count", len(conversations))
	return nil
}

// importPart4Talks 导入Part 4演讲
func (s *toeicImporterService) importPart4Talks(ctx context.Context, testID uint, talks []ParsedTalk) error {
	for _, talk := range talks {
		// 创建Part4Talk模型
		part4Talk := &model.Part4Talk{
			TestId:     int32(testID),
			TalkNumber: int32(talk.Number),
			Content:    talk.Content,
		}

		// 处理标题
		if talk.Title != "" {
			part4Talk.Title = sql.NullString{
				String: talk.Title,
				Valid:  true,
			}
		}

		// 处理问题文本
		if len(talk.Questions) > 0 {
			part4Talk.Question1 = s.getQuestionText(talk.Questions, 1)
		}
		if len(talk.Questions) > 1 {
			part4Talk.Question2 = s.getQuestionText(talk.Questions, 2)
		}
		if len(talk.Questions) > 2 {
			part4Talk.Question3 = s.getQuestionText(talk.Questions, 3)
		}

		// 处理场景ID
		if talk.ScenarioID != nil {
			part4Talk.ScenarioId = sql.NullInt32{
				Int32: int32(*talk.ScenarioID),
				Valid: true,
			}
		}

		// 处理难度级别ID
		if talk.DifficultyLevelID != nil {
			part4Talk.DifficultyLevelId = sql.NullInt32{
				Int32: int32(*talk.DifficultyLevelID),
				Valid: true,
			}
		}

		// 保存演讲到数据库
		if err := s.part4TalkRepo.Create(ctx, part4Talk); err != nil {
			return fmt.Errorf("failed to create Part 4 talk %d: %w", talk.Number, err)
		}

		s.logger.Info("Created Part 4 talk", "number", talk.Number, "id", part4Talk.ID)

		// 导入答案选项
		for i, question := range talk.Questions {
			if i >= 3 { // 只处理前3个问题
				break
			}

			part4AnswerOption := &model.Part4AnswerOption{
				TalkId:         int32(part4Talk.ID),
				QuestionNumber: int32(i + 1),
				OptionA:        question.OptionA,
				OptionB:        question.OptionB,
				OptionC:        question.OptionC,
				OptionD:        question.OptionD,
			}

			// 处理正确答案
			if question.CorrectAnswer != "" {
				part4AnswerOption.CorrectAnswer = sql.NullString{
					String: question.CorrectAnswer,
					Valid:  true,
				}
			}

			if err := s.part4AnswerOptionRepo.Create(ctx, part4AnswerOption); err != nil {
				return fmt.Errorf("failed to create Part 4 answer option for talk %d, question %d: %w", talk.Number, i+1, err)
			}
		}
	}

	s.logger.Info("Imported Part 4 talks", "count", len(talks))
	return nil
}

// getQuestionText 获取指定编号的问题文本
func (s *toeicImporterService) getQuestionText(questions []ParsedQuestion, number int) string {
	for _, q := range questions {
		if q.Number%10 == number || (len(questions) >= number && questions[number-1].Text != "") {
			if len(questions) >= number {
				return questions[number-1].Text
			}
		}
	}
	return ""
}
