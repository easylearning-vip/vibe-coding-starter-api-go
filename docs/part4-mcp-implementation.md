# Part 4 MCP Tools Implementation

## Overview

This document describes the implementation of Part 4 practice tools for the TOEIC MCP server, following the existing Part 2 and Part 3 implementation patterns.

## Implemented Tools

### 1. `auto_generate_part4_set`
Auto-generates a Part 4 practice set with multiple questions.

**Parameters:**
- `count` (required): Number of questions to include (10-30)
- `difficulty_level_id` (optional): Filter by difficulty level
- `scenario_id` (optional): Filter by scenario

**Functionality:**
- Creates a new practice set for the user
- Randomly or sequentially selects talks from the question bank
- Each talk has 3 questions, and items are created for individual questions
- Returns set ID and summary information

### 2. `get_part4_set_details`
Views Part 4 practice set details, including all question details and answer options.

**Parameters:**
- `set_id` (required): Practice set ID

**Functionality:**
- Retrieves practice set metadata (total questions, completed count, accuracy)
- Lists all questions in the set with:
  - Talk content and title
  - Specific question text for each item
  - Answer options (A, B, C, D) for each question
- Displays questions in order

### 3. `list_part4_sets`
Lists user's personal Part 4 practice sets with pagination support.

**Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Page size (default: 20, max: 100)
- `sort` (optional): Sort column
- `order` (optional): Sort order (asc/desc)

**Functionality:**
- Returns paginated list of practice sets
- Shows summary for each set: ID, total questions, completed count, correct count, accuracy

### 4. `part4_prompt`
Gets the latest Part 4 AI assistant prompt that provides instructions on how to use the TOEIC MCP tools effectively for Part 4 practice.

**Parameters:**
- None

**Functionality:**
- Retrieves the latest AI prompt from the database
- Returns the Part4Prompt field from the toeic_ai_prompts table

### 5. `submit_part4_answer`
Submits Part 4 practice answer with optional AI-judged correctness.

**Parameters:**
- `set_id` (required): Practice set ID
- `item_id` (required): Practice set item ID
- `user_answer` (required): User's answer choice (A, B, C, or D)
- `is_correct` (optional): Whether the answer is correct (if not provided, compares with database)

**Functionality:**
- Validates the answer format (must be A, B, C, or D)
- Determines correctness either by:
  - AI assistant judgment (if `is_correct` parameter provided)
  - Database comparison (if `is_correct` not provided)
- Updates practice set item with selected answer and correctness
- Updates practice set statistics (completed count, correct count, accuracy)
- Returns feedback with correct answer

## Architecture

### Files Created/Modified

#### New Files:
1. `internal/handler/mcp/part4_practice_set.go` - MCP tool handlers for Part 4 practice sets
2. `internal/handler/mcp/part4_prompt.go` - MCP tool handler for Part 4 AI prompt
3. `docs/part4-mcp-implementation.md` - This documentation

#### Modified Files:
1. `internal/model/mcp/params.go` - Added Part 4 parameter definitions
2. `cmd/toeic-mcp/main.go` - Registered Part 4 repositories and services
3. `internal/handler/mcp/server.go` - Added Part 4 service and registered Part 4 tools

### Existing Infrastructure Used

The implementation leverages existing infrastructure:

**Repository Layer:**
- `Part4PracticeSetRepository` - Already exists
- `Part4PracticeSetItemRepository` - Already exists
- `Part4TalkRepository` - Already exists (with `ListWithAnswerOptions` method)

**Service Layer:**
- `Part4PracticeSetService` - Already exists with all required methods:
  - `AutoGenerateSet`
  - `GetSetByID`
  - `ListSets`
  - `ListSetQuestions`
  - `GetItemDetail`
  - `SubmitAnswer`

**Model Layer:**
- `Part4PracticeSet` - Practice set master table
- `Part4PracticeSetItem` - Practice set items (references talk_id and question_index)
- `Part4Talk` - Talk content with 3 questions
- `Part4AnswerOption` - Answer options for each question in a talk

## Data Model

### Part4PracticeSet
- `user_id`: Owner of the practice set
- `scenario_id`: Optional scenario filter
- `difficulty_level_id`: Optional difficulty filter
- `total_questions`: Total number of questions
- `completed_count`: Number of completed questions
- `correct_count`: Number of correct answers
- `accuracy`: Calculated accuracy percentage

### Part4PracticeSetItem
- `set_id`: Reference to practice set
- `talk_id`: Reference to Part4Talk
- `question_index`: Which question (1, 2, or 3) from the talk
- `order_index`: Display order in the practice set
- `selected_answer`: User's selected answer (A/B/C/D)
- `is_correct`: Whether the answer was correct

### Part4Talk
- `test_id`: Reference to test
- `talk_number`: Talk number in test
- `title`: Optional talk title
- `content`: Talk content/transcript
- `question1`, `question2`, `question3`: Three questions for this talk
- `scenario_id`: Scenario classification
- `difficulty_level_id`: Difficulty level
- `answer_options`: Related answer options (1:N relationship)

### Part4AnswerOption
- `talk_id`: Reference to talk
- `question_number`: Which question (1, 2, or 3)
- `option_a`, `option_b`, `option_c`, `option_d`: Answer choices
- `correct_answer`: Correct answer (A/B/C/D)

## Key Design Decisions

### 1. Question-Level Granularity
Unlike Part 2 (single question per item), Part 4 items reference a specific question within a talk. Each talk has 3 questions, and practice set items are created at the question level, not the talk level.

### 2. Answer Options Structure
Part 4 uses a separate `Part4AnswerOption` table with `question_number` to distinguish between the three questions in each talk. The handler matches the question_index from the item with the question_number in the answer options.

### 3. AI-Assisted Grading
The `submit_part4_answer` tool supports an optional `is_correct` parameter, allowing AI assistants to judge answer correctness based on context rather than strict database comparison. This is useful for:
- Handling ambiguous questions
- Providing more nuanced feedback
- Supporting practice scenarios where database answers may not be set

### 4. Consistent API Design
All Part 4 tools follow the same patterns as Part 2 and Part 3:
- Same parameter naming conventions
- Same response formats
- Same error handling patterns
- Same pagination approach

## Testing Recommendations

### Manual Testing Steps:

1. **Generate a practice set:**
   ```json
   {
     "tool": "auto_generate_part4_set",
     "params": {
       "count": 15,
       "difficulty_level_id": 1,
       "scenario_id": 1
     }
   }
   ```

2. **List practice sets:**
   ```json
   {
     "tool": "list_part4_sets",
     "params": {
       "page": 1,
       "page_size": 10
     }
   }
   ```

3. **Get set details:**
   ```json
   {
     "tool": "get_part4_set_details",
     "params": {
       "set_id": 1
     }
   }
   ```

4. **Submit an answer:**
   ```json
   {
     "tool": "submit_part4_answer",
     "params": {
       "set_id": 1,
       "item_id": 1,
       "user_answer": "A"
     }
   }
   ```

5. **Get AI prompt:**
   ```json
   {
     "tool": "part4_prompt",
     "params": {}
   }
   ```

## Future Enhancements

Potential improvements for future iterations:

1. **Batch Answer Submission**: Allow submitting multiple answers at once
2. **Progress Tracking**: Add detailed progress tracking per talk
3. **Performance Analytics**: Track time spent per question
4. **Adaptive Difficulty**: Adjust question difficulty based on user performance
5. **Review Mode**: Allow users to review incorrect answers with explanations
6. **Audio Integration**: Support for audio playback of talks (when available)

## Conclusion

The Part 4 MCP tools implementation successfully mirrors the Part 2 and Part 3 patterns while accommodating the unique structure of Part 4 (talks with multiple questions). The implementation is production-ready and follows all established architectural patterns in the codebase.

