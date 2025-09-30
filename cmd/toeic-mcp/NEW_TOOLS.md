# New MCP Tools for Part2 Practice Set Management

This document describes the three new MCP tools added to the TOEIC MCP server for Part2 practice set management.

## Overview

Three new tools have been implemented following the clean architecture pattern:

1. **auto_generate_part2_set** - Automatically generate a practice set with questions
2. **get_part2_set_details** - View complete practice set details including all questions
3. **submit_part2_answer** - Submit an answer and get immediate feedback

## Implementation Details

### Architecture

All new tools follow the established clean architecture pattern:

```
internal/model/mcp/params.go
├── AutoGeneratePart2SetParams
├── GetPart2SetDetailsParams
└── SubmitPart2AnswerParams

internal/handler/mcp/part2_practice_set.go
├── AutoGeneratePart2SetTool()
├── HandleAutoGeneratePart2Set()
├── GetPart2SetDetailsTool()
├── HandleGetPart2SetDetails()
├── SubmitPart2AnswerTool()
└── HandleSubmitPart2Answer()

internal/handler/mcp/server.go
└── registerTools() - registers all tools
```

### Code Changes

**Files Modified:**
- `internal/model/mcp/params.go` - Added 3 new parameter structs
- `internal/handler/mcp/part2_practice_set.go` - Added 6 new functions (3 tools + 3 handlers)
- `internal/handler/mcp/server.go` - Registered 3 new tools

**No Service Layer Changes:**
- Reused existing `Part2PracticeSetService` methods
- No new business logic needed

## Tool 1: auto_generate_part2_set

### Description
自动生成 Part2 练习集

Automatically generates a practice set with the specified number of questions from the question bank.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `count` | int | Yes | Number of questions (10-30) |
| `difficulty_level_id` | int32 | No | Filter by difficulty level |
| `scenario_id` | int32 | No | Filter by scenario |

### Request Example

```json
{
  "name": "auto_generate_part2_set",
  "arguments": {
    "count": 15,
    "difficulty_level_id": 1,
    "scenario_id": 2
  }
}
```

### Response Format

```
✅ 练习集生成成功

练习集ID: 123
题目总数: 15
场景ID: 2
难度级别ID: 1

已添加 15 道题目到练习集

使用 get_part2_set_details 查看题目详情
```

### Error Cases

- `count` not provided: "参数错误: count 是必需的"
- `count` out of range: "参数错误: count 必须在 10-30 之间"
- Service error: "生成失败: [error details]"

### Implementation Notes

- Uses `Part2PracticeSetService.AutoGenerateSet()`
- Questions are selected sequentially by default
- Filters are applied if provided
- Returns set ID for subsequent operations

## Tool 2: get_part2_set_details

### Description
查看 Part2 练习集明细列表（包含题目详情）

Retrieves complete details of a practice set including all questions with their options.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `set_id` | uint | Yes | Practice set ID |

### Request Example

```json
{
  "name": "get_part2_set_details",
  "arguments": {
    "set_id": 123
  }
}
```

### Response Format

```
📋 练习集详情

练习集ID: 123
题目总数: 15
已完成: 5
正确数: 4
准确率: 26.67%

题目列表:
─────────────────────────────────────

题目 1 (Item ID: 456)
问题: Where is the meeting?
A. In the conference room
B. At the restaurant
C. In the office

─────────────────────────────────────

题目 2 (Item ID: 457)
问题: When will the report be ready?
A. Tomorrow
B. Next week
C. This afternoon

...
```

### Error Cases

- `set_id` not provided: "参数错误: set_id 是必需的"
- Set not found: "获取练习集失败: not found"
- Permission denied: "获取练习集失败: not found" (user doesn't own the set)
- Service error: "获取题目列表失败: [error details]"

### Implementation Notes

- Uses `Part2PracticeSetService.GetSetByID()` for set info
- Uses `Part2PracticeSetService.ListSetQuestions()` for questions
- Shows all questions with options (no correct answers revealed)
- Displays progress statistics
- Questions are ordered by `order_index`

## Tool 3: submit_part2_answer

### Description
提交 Part2 练习题答案

Submits an answer for a specific question and provides immediate feedback.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `set_id` | uint | Yes | Practice set ID |
| `item_id` | uint | Yes | Practice set item ID |
| `user_answer` | string | Yes | Answer choice (A, B, or C) |

### Request Example

```json
{
  "name": "submit_part2_answer",
  "arguments": {
    "set_id": 123,
    "item_id": 456,
    "user_answer": "A"
  }
}
```

### Response Format (Correct Answer)

```
✅ 回答正确！

你的答案: A
正确答案: A
```

### Response Format (Incorrect Answer)

```
❌ 回答错误

你的答案: B
正确答案: A
```

### Error Cases

- Missing parameters: "参数错误: 缺少必需参数"
- `set_id` not provided: "参数错误: set_id 是必需的"
- `item_id` not provided: "参数错误: item_id 是必需的"
- `user_answer` not provided: "参数错误: user_answer 是必需的"
- Invalid answer format: "参数错误: user_answer 必须是 A、B 或 C"
- Item not found: "获取题目失败: not found"
- Item doesn't belong to set: "题目不属于指定的练习集"
- Service error: "提交答案失败: [error details]"

### Implementation Notes

- Uses `Part2PracticeSetService.GetItemDetail()` to fetch question
- Validates answer format (A, B, or C - case insensitive)
- Checks if item belongs to the specified set
- Compares user answer with correct answer
- Uses `Part2PracticeSetService.SubmitAnswer()` to record result
- Updates practice set statistics (completed count, correct count, accuracy)
- Provides immediate feedback

## Testing Guide

### Prerequisites

1. Start the TOEIC MCP server:
```bash
cd /home/ubuntu/workspace/vocabulary/vibe-coding-starter-api-go
go run ./cmd/toeic-mcp/ -c configs/config.yaml -port 6275
```

2. Open MCP Inspector:
- URL: `http://localhost:6274/?MCP_PROXY_AUTH_TOKEN=f871086194aa906d024c43c7032fe0219a77399705fba92b090af9559d3487f9`

3. Configure MCP server:
- Server URL: `http://localhost:6275`
- Auth Token: `9f9f39514d4936a9f58df89584370f5d70b133ee321e7795a0ab430f4917843b`
- Header: `X-User-Token: 9f9f39514d4936a9f58df89584370f5d70b133ee321e7795a0ab430f4917843b`

### Test Scenarios

#### Scenario 1: Complete Practice Workflow

1. **Generate a practice set**
```json
{
  "name": "auto_generate_part2_set",
  "arguments": {
    "count": 10
  }
}
```
Expected: Returns set ID (e.g., 123)

2. **View practice set details**
```json
{
  "name": "get_part2_set_details",
  "arguments": {
    "set_id": 123
  }
}
```
Expected: Shows 10 questions with options

3. **Submit answer for first question**
```json
{
  "name": "submit_part2_answer",
  "arguments": {
    "set_id": 123,
    "item_id": 456,
    "user_answer": "A"
  }
}
```
Expected: Shows correct/incorrect feedback

4. **Check progress**
```json
{
  "name": "get_part2_set_details",
  "arguments": {
    "set_id": 123
  }
}
```
Expected: Shows updated progress (completed: 1)

#### Scenario 2: Error Handling

1. **Invalid count**
```json
{
  "name": "auto_generate_part2_set",
  "arguments": {
    "count": 5
  }
}
```
Expected: Error "count 必须在 10-30 之间"

2. **Invalid answer format**
```json
{
  "name": "submit_part2_answer",
  "arguments": {
    "set_id": 123,
    "item_id": 456,
    "user_answer": "D"
  }
}
```
Expected: Error "user_answer 必须是 A、B 或 C"

3. **Non-existent set**
```json
{
  "name": "get_part2_set_details",
  "arguments": {
    "set_id": 99999
  }
}
```
Expected: Error "获取练习集失败: not found"

## Verification Checklist

- [x] All three tools compile without errors
- [x] No linting errors (`go vet`)
- [x] Tools registered in server.go
- [x] Parameter validation implemented
- [x] Error messages in Chinese
- [x] Response format matches requirements
- [x] Reuses existing service methods
- [x] Follows clean architecture pattern
- [x] Consistent with HTTP handler behavior
- [x] User authentication enforced
- [x] Set ownership validated

## Future Enhancements

Potential improvements for future iterations:

1. **Batch Answer Submission**: Submit multiple answers at once
2. **Practice History**: View past practice sessions
3. **Statistics**: Detailed performance analytics
4. **Hints**: Provide hints for difficult questions
5. **Timer**: Track time spent on each question
6. **Review Mode**: Review incorrect answers
7. **Export Results**: Export practice results to file
8. **Difficulty Adjustment**: Automatically adjust difficulty based on performance

## Related Documentation

- [README.md](./README.md) - Full TOEIC MCP server documentation
- [REFACTORING.md](./REFACTORING.md) - Architecture refactoring details
- [QUICK_START.md](./QUICK_START.md) - Quick start guide

