package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// parsePart2 解析Part 2问题
func (s *toeicImporterService) parsePart2(lines []string) ([]ParsedPart2Question, error) {
	var questions []ParsedPart2Question
	var inPart2 bool

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "## Part 2" || strings.HasPrefix(line, "## Part 2:") {
			inPart2 = true
			continue
		}

		if inPart2 && strings.HasPrefix(line, "## Part") && line != "## Part 2" {
			break
		}

		if !inPart2 {
			continue
		}

		// Look for question pattern: **8. Are you available for an employee evaluation tomorrow?**
		questionRegex := regexp.MustCompile(`^\*\*(\d+)\.\s*(.+)\*\*$`)
		matches := questionRegex.FindStringSubmatch(line)

		if len(matches) == 3 {
			questionNum, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			questionText := strings.TrimSpace(matches[2])

			// Look for scenario/difficulty on next line
			if i+1 < len(lines) {
				scenarioID, difficultyID, err := s.parseScenarioAndDifficulty(lines[i+1])
				if err != nil {
					s.logger.Warn("Failed to parse scenario/difficulty for question", "question", questionNum, "error", err)
				}

				// Parse options A, B, C (handle both uppercase and lowercase)
				var optionA, optionB, optionC string
				for j := i + 2; j < len(lines) && j < i+5; j++ {
					optLine := strings.TrimSpace(lines[j])
					if strings.HasPrefix(optLine, "A. ") || strings.HasPrefix(optLine, "a. ") {
						optionA = strings.TrimSpace(optLine[3:])
					} else if strings.HasPrefix(optLine, "B. ") || strings.HasPrefix(optLine, "b. ") {
						optionB = strings.TrimSpace(optLine[3:])
					} else if strings.HasPrefix(optLine, "C. ") || strings.HasPrefix(optLine, "c. ") {
						optionC = strings.TrimSpace(optLine[3:])
					}
				}

				if optionA != "" && optionB != "" && optionC != "" {
					questions = append(questions, ParsedPart2Question{
						Number:            questionNum,
						Text:              questionText,
						OptionA:           optionA,
						OptionB:           optionB,
						OptionC:           optionC,
						CorrectAnswer:     "", // Will be determined later or left empty
						ScenarioID:        scenarioID,
						DifficultyLevelID: difficultyID,
					})
				}
			}
		}

		// Look for alternative question pattern: ### Question 7 [1|Workplace/Business] - **Difficulty: Intermediate**
		altQuestionRegex := regexp.MustCompile(`^### Question (\d+) `)
		matches = altQuestionRegex.FindStringSubmatch(line)

		if len(matches) == 2 {
			questionNum, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			// Parse scenario/difficulty from same line
			scenarioID, difficultyID, err := s.parseScenarioAndDifficulty(line)
			if err != nil {
				s.logger.Warn("Failed to parse scenario/difficulty for question", "question", questionNum, "error", err)
			}

			// Look for question text on next line (in bold)
			var questionText string
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if strings.HasPrefix(nextLine, "**") && strings.HasSuffix(nextLine, "**") {
					questionText = strings.TrimSpace(nextLine[2 : len(nextLine)-2])
					// Handle **Q:** format
					if strings.HasPrefix(questionText, "Q: ") {
						questionText = strings.TrimSpace(questionText[3:])
					}
				} else if strings.HasPrefix(nextLine, "**Q:**") {
					// Handle **Q:** format where Q: is bold but question text follows
					questionText = strings.TrimSpace(nextLine[5:]) // Remove **Q:** prefix
					// Remove any leading asterisk that might remain
					if strings.HasPrefix(questionText, "*") {
						questionText = strings.TrimSpace(questionText[1:])
					}
				}
			}

			if questionText == "" {
				continue
			}

			// Parse options A, B, C (handle both uppercase and lowercase, and dash format)
			var optionA, optionB, optionC string
			for j := i + 2; j < len(lines) && j < i+6; j++ {
				optLine := strings.TrimSpace(lines[j])
				if optLine == "" {
					continue
				}
				if strings.HasPrefix(optLine, "A. ") || strings.HasPrefix(optLine, "a. ") {
					optionA = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "B. ") || strings.HasPrefix(optLine, "b. ") {
					optionB = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "C. ") || strings.HasPrefix(optLine, "c. ") {
					optionC = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "- A. ") || strings.HasPrefix(optLine, "- a. ") {
					optionA = strings.TrimSpace(optLine[5:])
				} else if strings.HasPrefix(optLine, "- B. ") || strings.HasPrefix(optLine, "- b. ") {
					optionB = strings.TrimSpace(optLine[5:])
				} else if strings.HasPrefix(optLine, "- C. ") || strings.HasPrefix(optLine, "- c. ") {
					optionC = strings.TrimSpace(optLine[5:])
				} else if strings.HasPrefix(optLine, "###") {
					// Hit next question
					break
				}
			}

			if optionA != "" && optionB != "" && optionC != "" {
				questions = append(questions, ParsedPart2Question{
					Number:            questionNum,
					Text:              questionText,
					OptionA:           optionA,
					OptionB:           optionB,
					OptionC:           optionC,
					CorrectAnswer:     "", // Will be determined later or left empty
					ScenarioID:        scenarioID,
					DifficultyLevelID: difficultyID,
				})
			}
		}

		// Look for **Question 7** [1|Workplace/Business] - Difficulty: Beginner format
		altQuestionRegex2 := regexp.MustCompile(`^\*\*Question (\d+)\*\* `)
		matches = altQuestionRegex2.FindStringSubmatch(line)

		if len(matches) == 2 {
			questionNum, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			// Parse scenario/difficulty from same line
			scenarioID, difficultyID, err := s.parseScenarioAndDifficulty(line)
			if err != nil {
				s.logger.Warn("Failed to parse scenario/difficulty for question", "question", questionNum, "error", err)
			}

			// Look for question text on next line with "- Q:" format
			var questionText string
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if strings.HasPrefix(nextLine, "- Q:") {
					questionText = strings.TrimSpace(nextLine[4:]) // Remove "- Q:" prefix
				}
			}

			if questionText != "" {
				// Parse answer options A, B, C
				var optionA, optionB, optionC string
				j := i + 2
				for ; j < len(lines) && j < i+6; j++ {
					optLine := strings.TrimSpace(lines[j])
					if optLine == "" {
						continue
					}

					if strings.HasPrefix(optLine, "- A:") {
						optionA = strings.TrimSpace(optLine[4:])
					} else if strings.HasPrefix(optLine, "- B:") {
						optionB = strings.TrimSpace(optLine[4:])
					} else if strings.HasPrefix(optLine, "- C:") {
						optionC = strings.TrimSpace(optLine[4:])
					} else if strings.HasPrefix(optLine, "**Question") || strings.HasPrefix(optLine, "###") {
						// Hit next question or section
						break
					}
				}

				if optionA != "" && optionB != "" && optionC != "" {
					questions = append(questions, ParsedPart2Question{
						Number:            questionNum,
						Text:              questionText,
						OptionA:           optionA,
						OptionB:           optionB,
						OptionC:           optionC,
						CorrectAnswer:     "", // Will be determined later or left empty
						ScenarioID:        scenarioID,
						DifficultyLevelID: difficultyID,
					})
				}
			}
		}
	}

	s.logger.Info("Parsed Part 2 questions", "count", len(questions))
	return questions, nil
}

// parsePart3 解析Part 3对话
func (s *toeicImporterService) parsePart3(lines []string) ([]ParsedConversation, error) {
	var conversations []ParsedConversation
	var inPart3 bool

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "## Part 3" || strings.HasPrefix(line, "## Part 3:") || line == "## Part 3: Conversations" {
			inPart3 = true
			continue
		}

		if inPart3 && strings.HasPrefix(line, "## Part") && line != "## Part 3" {
			break
		}

		if !inPart3 {
			continue
		}

		// Look for conversation pattern: ### Conversation 1 (Questions 32-34)
		convRegex := regexp.MustCompile(`^### Conversation (\d+) \(Questions (\d+)-(\d+)\)$`)
		matches := convRegex.FindStringSubmatch(line)

		// Also look for alternative pattern: ### Questions 32-34: Bookstore Conversation [1|Workplace/Business] - **Difficulty: Intermediate**
		if len(matches) == 0 {
			altConvRegex := regexp.MustCompile(`^### Questions (\d+)-(\d+): (.+) \[`)
			altMatches := altConvRegex.FindStringSubmatch(line)
			if len(altMatches) == 4 {
				// Convert to standard format for processing
				startQ, _ := strconv.Atoi(altMatches[1])
				convNum := (startQ-32)/3 + 1 // Calculate conversation number based on question range
				matches = []string{line, strconv.Itoa(convNum), altMatches[1], altMatches[2]}
			}
		}

		// Also look for pattern without space before bracket
		if len(matches) == 0 {
			altConvRegex2 := regexp.MustCompile(`^### Questions (\d+)-(\d+): (.+)\[`)
			altMatches := altConvRegex2.FindStringSubmatch(line)
			if len(altMatches) == 4 {
				startQ, _ := strconv.Atoi(altMatches[1])
				convNum := (startQ-32)/3 + 1
				matches = []string{line, strconv.Itoa(convNum), altMatches[1], altMatches[2]}
			}
		}

		// Also look for pattern without title
		if len(matches) == 0 {
			altConvRegex3 := regexp.MustCompile(`^### Questions (\d+)-(\d+) \[`)
			altMatches := altConvRegex3.FindStringSubmatch(line)
			if len(altMatches) == 3 {
				startQ, _ := strconv.Atoi(altMatches[1])
				convNum := (startQ-32)/3 + 1
				matches = []string{line, strconv.Itoa(convNum), altMatches[1], altMatches[2]}
			}
		}

		// Also look for simple format: ### Questions 32-34
		if len(matches) == 0 {
			altConvRegex4 := regexp.MustCompile(`^### Questions (\d+)-(\d+)$`)
			altMatches := altConvRegex4.FindStringSubmatch(line)
			if len(altMatches) == 3 {
				startQ, _ := strconv.Atoi(altMatches[1])
				convNum := (startQ-32)/3 + 1
				matches = []string{line, strconv.Itoa(convNum), altMatches[1], altMatches[2]}
			}
		}

		// Also look for format with title: ### Questions 32-34: Employee Performance Discussion
		if len(matches) == 0 {
			altConvRegex5 := regexp.MustCompile(`^### Questions (\d+)-(\d+): (.+)$`)
			altMatches := altConvRegex5.FindStringSubmatch(line)
			if len(altMatches) == 4 {
				startQ, _ := strconv.Atoi(altMatches[1])
				convNum := (startQ-32)/3 + 1
				matches = []string{line, strconv.Itoa(convNum), altMatches[1], altMatches[2]}
			}
		}

		if len(matches) == 4 {
			convNum, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			startQ, _ := strconv.Atoi(matches[2])
			endQ, _ := strconv.Atoi(matches[3])

			// Parse scenario/difficulty from current line or next line
			var scenarioID, difficultyID *uint
			scenarioID, difficultyID, err = s.parseScenarioAndDifficulty(line)
			if err != nil && i+1 < len(lines) {
				scenarioID, difficultyID, err = s.parseScenarioAndDifficulty(lines[i+1])
				if err != nil {
					s.logger.Warn("Failed to parse scenario/difficulty for conversation", "conversation", convNum, "error", err)
				}
			}

			// Parse conversation content and questions
			conversation := s.parseConversationContent(lines, i, convNum, startQ, endQ, scenarioID, difficultyID)
			if conversation != nil {
				conversations = append(conversations, *conversation)
			}
		}
	}

	s.logger.Info("Parsed Part 3 conversations", "count", len(conversations))
	return conversations, nil
}

// parseConversationContent 解析对话内容和问题
func (s *toeicImporterService) parseConversationContent(lines []string, startIndex, convNum, startQ, endQ int, scenarioID, difficultyID *uint) *ParsedConversation {
	var content strings.Builder
	var questions []ParsedQuestion
	j := startIndex + 2

	// Skip scenario/difficulty line if it exists
	if j < len(lines) && strings.Contains(lines[j], "*") && strings.Contains(lines[j], "Difficulty:") {
		j++
	}

	// Parse conversation dialogue and embedded questions
	var allContent strings.Builder
	for ; j < len(lines); j++ {
		line := strings.TrimSpace(lines[j])
		if line == "" {
			continue
		}

		// Stop if we hit next conversation or part
		if strings.HasPrefix(line, "###") || strings.HasPrefix(line, "## Part") {
			break
		}

		allContent.WriteString(line + "\n")
	}

	// Now parse the content to separate dialogue from questions
	contentStr := allContent.String()

	// Check if questions are embedded in the content (format: **Q32:** or **Question 32:**)
	if strings.Contains(contentStr, "**Q"+fmt.Sprintf("%d", startQ)+":**") ||
		strings.Contains(contentStr, "**Question "+fmt.Sprintf("%d", startQ)+":**") {
		// Parse embedded questions format
		questions = s.parseEmbeddedQuestions(contentStr, startQ, endQ)
		// Extract only dialogue content (everything before first question)
		content.WriteString(s.extractDialogueFromEmbeddedContent(contentStr, startQ))
	} else {
		// Traditional format - questions are separate from dialogue
		// First, separate dialogue from questions in the content
		contentLines := strings.Split(contentStr, "\n")
		var dialogueLines []string
		var questionStartIdx = -1

		for i, line := range contentLines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Check if this line starts a question (multiple patterns)
			questionPatterns := []string{
				`^\*\*Question (\d+):\*\*`,
				`^\*\*Question (\d+)\*\*\s*$`,
				`^\*\*(\d+)\.\s*(.+)\*\*$`,
			}

			isQuestion := false
			for _, pattern := range questionPatterns {
				if matched, _ := regexp.MatchString(pattern, line); matched {
					questionStartIdx = i
					isQuestion = true
					break
				}
			}

			if isQuestion {
				break
			}

			// Add to dialogue if it's conversation content
			if strings.HasPrefix(line, "**") && strings.Contains(line, ":**") {
				// Speaker lines like **Woman:** or **Man:**
				dialogueLines = append(dialogueLines, line)
			} else if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") && !strings.Contains(line, "Question") {
				// Other bold text that's not a question
				dialogueLines = append(dialogueLines, line)
			} else if len(dialogueLines) > 0 && !strings.HasPrefix(line, "A. ") && !strings.HasPrefix(line, "B. ") && !strings.HasPrefix(line, "C. ") && !strings.HasPrefix(line, "D. ") {
				// Regular dialogue content (but not answer options)
				dialogueLines = append(dialogueLines, line)
			}
		}

		// Set dialogue content
		content.WriteString(strings.Join(dialogueLines, "\n"))

		// Parse questions from the remaining content lines
		if questionStartIdx >= 0 {
			questionLines := contentLines[questionStartIdx:]
			for qNum := startQ; qNum <= endQ; qNum++ {
				question, _ := s.parseQuestion(questionLines, 0, qNum)
				if question != nil {
					questions = append(questions, *question)
				}
			}
		}
	}

	return &ParsedConversation{
		Number:            convNum,
		Title:             fmt.Sprintf("Conversation %d (Questions %d-%d)", convNum, startQ, endQ),
		Content:           strings.TrimSpace(content.String()),
		ScenarioID:        scenarioID,
		DifficultyLevelID: difficultyID,
		Questions:         questions,
	}
}

// parsePart4 解析Part 4演讲
func (s *toeicImporterService) parsePart4(lines []string) ([]ParsedTalk, error) {
	var talks []ParsedTalk
	var inPart4 bool

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "## Part 4" || strings.HasPrefix(line, "## Part 4:") {
			inPart4 = true
			continue
		}

		if inPart4 && strings.HasPrefix(line, "## Part") && line != "## Part 4" {
			break
		}

		if !inPart4 {
			continue
		}

		// Look for talk pattern: ### Talk 1 (Questions 71-73)
		talkRegex := regexp.MustCompile(`^### Talk (\d+) \(Questions (\d+)-(\d+)\)$`)
		matches := talkRegex.FindStringSubmatch(line)

		// Also look for alternative pattern: ### Questions 71-73: Radio Broadcast [5|Media/Educational] - **Difficulty: Advanced**
		if len(matches) == 0 {
			altTalkRegex := regexp.MustCompile(`^### Questions (\d+)-(\d+): (.+) \[`)
			altMatches := altTalkRegex.FindStringSubmatch(line)
			if len(altMatches) == 4 {
				startQ, _ := strconv.Atoi(altMatches[1])
				talkNum := (startQ-71)/3 + 1
				matches = []string{line, strconv.Itoa(talkNum), altMatches[1], altMatches[2]}
			}
		}

		// Also look for pattern without space before bracket
		if len(matches) == 0 {
			altTalkRegex2 := regexp.MustCompile(`^### Questions (\d+)-(\d+): (.+)\[`)
			altMatches := altTalkRegex2.FindStringSubmatch(line)
			if len(altMatches) == 4 {
				startQ, _ := strconv.Atoi(altMatches[1])
				talkNum := (startQ-71)/3 + 1
				matches = []string{line, strconv.Itoa(talkNum), altMatches[1], altMatches[2]}
			}
		}

		// Also look for pattern without title
		if len(matches) == 0 {
			altTalkRegex3 := regexp.MustCompile(`^### Questions (\d+)-(\d+) \[`)
			altMatches := altTalkRegex3.FindStringSubmatch(line)
			if len(altMatches) == 3 {
				startQ, _ := strconv.Atoi(altMatches[1])
				talkNum := (startQ-71)/3 + 1
				matches = []string{line, strconv.Itoa(talkNum), altMatches[1], altMatches[2]}
			}
		}

		if len(matches) == 4 {
			talkNum, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			startQ, _ := strconv.Atoi(matches[2])
			endQ, _ := strconv.Atoi(matches[3])

			// Parse scenario/difficulty from current line or next line
			var scenarioID, difficultyID *uint
			scenarioID, difficultyID, err = s.parseScenarioAndDifficulty(line)
			if err != nil && i+1 < len(lines) {
				scenarioID, difficultyID, err = s.parseScenarioAndDifficulty(lines[i+1])
				if err != nil {
					s.logger.Warn("Failed to parse scenario/difficulty for talk", "talk", talkNum, "error", err)
				}
			}

			// Parse talk content and questions
			talk := s.parseTalkContent(lines, i, talkNum, startQ, endQ, scenarioID, difficultyID)
			if talk != nil {
				talks = append(talks, *talk)
			}
		}
	}

	s.logger.Info("Parsed Part 4 talks", "count", len(talks))
	return talks, nil
}

// parseTalkContent 解析演讲内容和问题
func (s *toeicImporterService) parseTalkContent(lines []string, startIndex, talkNum, startQ, endQ int, scenarioID, difficultyID *uint) *ParsedTalk {
	var content strings.Builder
	var questions []ParsedQuestion
	j := startIndex + 2

	// Skip scenario/difficulty line if it exists
	if j < len(lines) && strings.Contains(lines[j], "*") && strings.Contains(lines[j], "Difficulty:") {
		j++
	}

	// Parse talk content and embedded questions
	var allContent strings.Builder
	for ; j < len(lines); j++ {
		line := strings.TrimSpace(lines[j])
		if line == "" {
			continue
		}

		// Stop if we hit next talk or part
		if strings.HasPrefix(line, "###") || strings.HasPrefix(line, "## Part") {
			break
		}

		allContent.WriteString(line + "\n")
	}

	// Now parse the content to separate talk content from questions
	contentStr := allContent.String()

	// Check if questions are embedded in the content (format: **Q71:** or **Question 71:**)
	if strings.Contains(contentStr, "**Q"+fmt.Sprintf("%d", startQ)+":**") ||
		strings.Contains(contentStr, "**Question "+fmt.Sprintf("%d", startQ)+":**") {
		// Parse embedded questions format
		questions = s.parseEmbeddedQuestions(contentStr, startQ, endQ)
		// Extract only talk content (everything before first question)
		content.WriteString(s.extractDialogueFromEmbeddedContent(contentStr, startQ))
	} else {
		// Traditional format - questions are separate from talk content
		// First, separate talk content from questions
		contentLines := strings.Split(contentStr, "\n")
		var talkLines []string
		var questionStartIdx = -1

		for i, line := range contentLines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Check if this line starts a question (multiple patterns)
			questionPatterns := []string{
				`^\*\*Question (\d+):\*\*`,
				`^\*\*Question (\d+)\*\*\s*$`,
				`^\*\*(\d+)\.\s*(.+)\*\*$`,
			}

			isQuestion := false
			for _, pattern := range questionPatterns {
				if matched, _ := regexp.MatchString(pattern, line); matched {
					questionStartIdx = i
					isQuestion = true
					break
				}
			}

			if isQuestion {
				break
			}

			// Add to talk content
			if strings.HasPrefix(line, "**") && strings.Contains(line, ":**") {
				// Speaker lines like **Host:** or **Announcer:**
				talkLines = append(talkLines, line)
			} else if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") && !strings.Contains(line, "Question") {
				// Other bold text that's not a question
				talkLines = append(talkLines, line)
			} else if len(talkLines) > 0 || (!strings.HasPrefix(line, "**") && !strings.HasPrefix(line, "A. ") && !strings.HasPrefix(line, "B. ") && !strings.HasPrefix(line, "C. ") && !strings.HasPrefix(line, "D. ")) {
				// Regular talk content (but not answer options)
				talkLines = append(talkLines, line)
			}
		}

		// Set talk content
		content.WriteString(strings.Join(talkLines, "\n"))

		// Parse questions from the remaining content lines
		if questionStartIdx >= 0 {
			questionLines := contentLines[questionStartIdx:]
			for qNum := startQ; qNum <= endQ; qNum++ {
				question, _ := s.parseQuestion(questionLines, 0, qNum)
				if question != nil {
					questions = append(questions, *question)
				}
			}
		}
	}

	return &ParsedTalk{
		Number:            talkNum,
		Title:             fmt.Sprintf("Talk %d (Questions %d-%d)", talkNum, startQ, endQ),
		Content:           strings.TrimSpace(content.String()),
		ScenarioID:        scenarioID,
		DifficultyLevelID: difficultyID,
		Questions:         questions,
	}
}

// parseQuestion 解析单个问题及其选项
func (s *toeicImporterService) parseQuestion(lines []string, startIndex int, expectedNum int) (*ParsedQuestion, int) {
	for i := startIndex; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Look for question pattern: **32. What kind of event are the speakers planning?**
		questionRegex := regexp.MustCompile(`^\*\*(\d+)\.\s*(.+)\*\*$`)
		matches := questionRegex.FindStringSubmatch(line)

		// Also look for alternative pattern: **Question 32:** What kind of business do the speakers work for?
		if len(matches) == 0 {
			altQuestionRegex := regexp.MustCompile(`^\*\*Question (\d+):\*\*\s*(.+)$`)
			matches = altQuestionRegex.FindStringSubmatch(line)
		}

		// Also look for pattern: **Question 32:** (question text on same line or next line)
		if len(matches) == 0 {
			altQuestionRegex2 := regexp.MustCompile(`^\*\*Question (\d+):\*\*(.*)$`)
			matches = altQuestionRegex2.FindStringSubmatch(line)
			if len(matches) == 3 && strings.TrimSpace(matches[2]) == "" {
				// Question text is on next line
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					matches[2] = nextLine
				}
			}
		}

		// Also look for pattern: **Question 32** (without colon, used in test09)
		if len(matches) == 0 {
			altQuestionRegex3 := regexp.MustCompile(`^\*\*Question (\d+)\*\*\s*$`)
			matches = altQuestionRegex3.FindStringSubmatch(line)
			if len(matches) == 2 {
				// Question text is on next line
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					if nextLine != "" && !strings.HasPrefix(nextLine, "A. ") {
						matches = []string{line, matches[1], nextLine}
					}
				}
			}
		}

		if len(matches) == 3 {
			questionNum, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			// Skip if not the expected question number
			if questionNum != expectedNum {
				continue
			}

			questionText := strings.TrimSpace(matches[2])

			// Parse options A, B, C, D (handle both uppercase and lowercase)
			var optionA, optionB, optionC, optionD string
			j := i + 1

			for ; j < len(lines) && j < i+6; j++ {
				optLine := strings.TrimSpace(lines[j])
				if optLine == "" {
					continue
				}

				if strings.HasPrefix(optLine, "A. ") || strings.HasPrefix(optLine, "a. ") {
					optionA = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "B. ") || strings.HasPrefix(optLine, "b. ") {
					optionB = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "C. ") || strings.HasPrefix(optLine, "c. ") {
					optionC = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "D. ") || strings.HasPrefix(optLine, "d. ") {
					optionD = strings.TrimSpace(optLine[3:])
				} else if strings.HasPrefix(optLine, "**") || strings.HasPrefix(optLine, "###") {
					// Hit next question or section
					break
				}
			}

			if optionA != "" && optionB != "" && optionC != "" && optionD != "" {
				return &ParsedQuestion{
					Number:        questionNum,
					Text:          questionText,
					OptionA:       optionA,
					OptionB:       optionB,
					OptionC:       optionC,
					OptionD:       optionD,
					CorrectAnswer: "", // Will be determined later or left empty
				}, j
			}
		}
	}

	return nil, len(lines)
}

// parseEmbeddedQuestions 解析嵌入在内容中的问题
func (s *toeicImporterService) parseEmbeddedQuestions(content string, startQ, endQ int) []ParsedQuestion {
	var questions []ParsedQuestion
	lines := strings.Split(content, "\n")

	for qNum := startQ; qNum <= endQ; qNum++ {
		// Look for question patterns: **Q32:** or **Question 32:**
		var questionText string
		var questionStartIdx = -1

		for i, line := range lines {
			line = strings.TrimSpace(line)

			// Pattern 1: **Q32:** question text
			qPattern1 := regexp.MustCompile(`^\*\*Q(\d+):\*\*\s*(.+)$`)
			matches := qPattern1.FindStringSubmatch(line)

			if len(matches) == 3 {
				if num, err := strconv.Atoi(matches[1]); err == nil && num == qNum {
					questionText = strings.TrimSpace(matches[2])
					questionStartIdx = i
					break
				}
			}

			// Pattern 2: **Question 32:** question text
			qPattern2 := regexp.MustCompile(`^\*\*Question (\d+):\*\*\s*(.+)$`)
			matches = qPattern2.FindStringSubmatch(line)

			if len(matches) == 3 {
				if num, err := strconv.Atoi(matches[1]); err == nil && num == qNum {
					questionText = strings.TrimSpace(matches[2])
					questionStartIdx = i
					break
				}
			}

			// Pattern 3: **Question 32** (without colon, used in test09)
			qPattern3 := regexp.MustCompile(`^\*\*Question (\d+)\*\*\s*$`)
			matches = qPattern3.FindStringSubmatch(line)

			if len(matches) == 2 {
				if num, err := strconv.Atoi(matches[1]); err == nil && num == qNum {
					// Question text is on the next line
					if i+1 < len(lines) {
						nextLine := strings.TrimSpace(lines[i+1])
						if nextLine != "" && !strings.HasPrefix(nextLine, "A. ") {
							questionText = nextLine
							questionStartIdx = i
							break
						}
					}
				}
			}
		}

		if questionStartIdx == -1 {
			s.logger.Warn("Question not found", "questionNumber", qNum)
			continue
		}

		// Parse answer options A, B, C, D
		var optionA, optionB, optionC, optionD string

		// Start looking for options from the line after the question
		startIdx := questionStartIdx + 1
		if questionText == "" && questionStartIdx+2 < len(lines) {
			// If question text was on next line, start from line after that
			startIdx = questionStartIdx + 2
		}

		for i := startIdx; i < len(lines) && i < startIdx+8; i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "A. ") || strings.HasPrefix(line, "A.  ") {
				optionA = strings.TrimSpace(line[3:])
			} else if strings.HasPrefix(line, "B. ") || strings.HasPrefix(line, "B.  ") {
				optionB = strings.TrimSpace(line[3:])
			} else if strings.HasPrefix(line, "C. ") || strings.HasPrefix(line, "C.  ") {
				optionC = strings.TrimSpace(line[3:])
			} else if strings.HasPrefix(line, "D. ") || strings.HasPrefix(line, "D.  ") {
				optionD = strings.TrimSpace(line[3:])
			} else if strings.HasPrefix(line, "**Q") || strings.HasPrefix(line, "**Question") || strings.HasPrefix(line, "###") {
				// Hit next question or section
				break
			}
		}

		if optionA != "" && optionB != "" && optionC != "" && optionD != "" {
			questions = append(questions, ParsedQuestion{
				Number:        qNum,
				Text:          questionText,
				OptionA:       optionA,
				OptionB:       optionB,
				OptionC:       optionC,
				OptionD:       optionD,
				CorrectAnswer: "", // Will be determined later or left empty
			})
		} else {
			s.logger.Warn("Incomplete question options", "questionNumber", qNum, "optionA", optionA, "optionB", optionB, "optionC", optionC, "optionD", optionD)
		}
	}

	return questions
}

// extractDialogueFromEmbeddedContent 从嵌入问题的内容中提取对话部分
func (s *toeicImporterService) extractDialogueFromEmbeddedContent(content string, startQ int) string {
	lines := strings.Split(content, "\n")
	var dialogue strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Stop when we hit the first question
		qPattern1 := regexp.MustCompile(`^\*\*Q(\d+):\*\*`)
		qPattern2 := regexp.MustCompile(`^\*\*Question (\d+):\*\*`)
		qPattern3 := regexp.MustCompile(`^\*\*Question (\d+)\*\*\s*$`)

		if qPattern1.MatchString(line) || qPattern2.MatchString(line) || qPattern3.MatchString(line) {
			break
		}

		// Add dialogue lines (speaker lines and regular content)
		if strings.HasPrefix(line, "**") && strings.Contains(line, ":**") {
			// Speaker lines like **Woman:** or **Man:**
			dialogue.WriteString(line + "\n")
		} else if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") && !strings.Contains(line, ":") {
			// Other bold text that's not a question
			dialogue.WriteString(line + "\n")
		} else if dialogue.Len() > 0 && !strings.HasPrefix(line, "A. ") && !strings.HasPrefix(line, "B. ") && !strings.HasPrefix(line, "C. ") && !strings.HasPrefix(line, "D. ") {
			// Regular dialogue content (but not answer options)
			dialogue.WriteString(line + "\n")
		}
	}

	return strings.TrimSpace(dialogue.String())
}
