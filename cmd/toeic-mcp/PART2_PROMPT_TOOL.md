# part2_prompt Tool Documentation

## Overview

The `part2_prompt` tool provides AI assistants with comprehensive instructions on how to effectively use the TOEIC MCP tools for Part2 practice. This tool returns the AI prompt stored in the database, which contains guidelines, best practices, and usage patterns.

## Tool Details

### Tool Name
`part2_prompt`

### Description
Get the latest Part2 AI assistant prompt for using TOEIC MCP tools. This prompt provides instructions on how to use the MCP tools effectively.

### Parameters

**No parameters required.** The tool automatically returns the latest prompt (highest ID) from the database.

### Response Format

Returns the complete Part2 prompt text from the latest database record.

## Use Cases

### Use Case: Get Latest Prompt
Retrieve the latest AI assistant prompt:

```json
{
  "name": "part2_prompt",
  "arguments": {}
}
```

**Note:** The tool always returns the most recent prompt based on the highest ID in the database.

## Implementation Details

### Architecture

```
part2_prompt Tool
    ↓
HandlePart2Prompt (Handler)
    ↓
ToeicAiPromptService (Service Layer)
    ↓
ToeicAiPromptRepository (Repository Layer)
    ↓
Database (toeic_ai_prompts table)
```

### Code Structure

#### 1. Parameter Model (`internal/model/mcp/params.go`)

```go
type Part2PromptParams struct {
    // No parameters needed - always returns the latest prompt
}
```

#### 2. Handler (`internal/handler/mcp/part2_prompt.go`)

```go
func Part2PromptTool() *mcp.Tool {
    return &mcp.Tool{
        Name:        "part2_prompt",
        Description: "Get the latest Part2 AI assistant prompt for using TOEIC MCP tools.",
    }
}

func HandlePart2Prompt(user *model.User, promptService service.ToeicAiPromptService)
    func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part2PromptParams)
    (*mcp.CallToolResult, any, error)
```

#### 3. Service Integration

The tool uses the existing `ToeicAiPromptService` which provides:
- `List(ctx, opts)` - List prompts with pagination and sorting

### Latest Prompt Logic

The tool automatically retrieves the latest prompt:

1. **Query with sorting**: Query the database with `Sort: "id"` and `Order: "desc"`
2. **Get first result**: Retrieve only the first record (highest ID)
3. **Error if none available**: Return error if no prompts exist in the database

```go
// Get the latest prompt (ordered by ID DESC, get the first one)
prompts, total, err := promptService.List(ctx, &service.ListToeicAiPromptOptions{
    Page:     1,
    PageSize: 1,
    Sort:     "id",
    Order:    "desc",
})
if err != nil || total == 0 {
    return nil, nil, fmt.Errorf("failed to get Part2 prompt: no prompts available")
}

// Return the Part2 prompt from the latest record
return &mcp.CallToolResult{
    Content: []mcp.Content{&mcp.TextContent{Text: prompts[0].Part2Prompt}},
}, nil, nil
```

## Database Schema

### Table: `toeic_ai_prompts`

| Column | Type | Description |
|--------|------|-------------|
| `id` | uint | Primary key |
| `part2_prompt` | text | Part2 AI prompt content |
| `part3_prompt` | text | Part3 AI prompt content |
| `part4_prompt` | text | Part4 AI prompt content |
| `created_at` | timestamp | Creation time |
| `updated_at` | timestamp | Last update time |

## Response Example

```
You are a TOEIC Part2 practice assistant. Your role is to help students practice TOEIC Part2 questions using the available MCP tools.

## Available Tools

1. **part2_prompt** - Get this prompt
2. **get_user_info** - Get current user information
3. **list_part2_sets** - List practice sets
4. **auto_generate_part2_set** - Generate new practice set
5. **get_part2_set_details** - View practice set details
6. **submit_part2_answer** - Submit and evaluate answers

## Practice Workflow

1. Generate or select a practice set
2. Present questions one by one
3. Collect user's answer
4. Evaluate and provide feedback
5. Track progress

## Best Practices

- Always validate user input
- Provide clear feedback
- Track completion progress
- Encourage learning from mistakes

[... more content ...]
```

## Integration with MCP Handler

### Dependency Injection

The tool is integrated into the MCP server through dependency injection:

**main.go:**
```go
// Repository layer
fx.Provide(
    repository.NewToeicAiPromptRepository,
)

// Service layer
fx.Provide(
    service.NewToeicAiPromptService,
)

// Handler layer
fx.Provide(
    mcpHandler.NewMCPHandler,
)
```

**server.go:**
```go
type MCPHandler struct {
    promptService service.ToeicAiPromptService
}

func (h *MCPHandler) registerTools(server *mcp.Server, user *model.User) {
    mcp.AddTool(server, Part2PromptTool(), HandlePart2Prompt(user, h.promptService))
}
```

## Error Handling

### Error Cases

| Error | Cause | Message |
|-------|-------|---------|
| No prompts available | Database has no prompt records | `"failed to get Part2 prompt: no prompts available"` |
| Database error | Connection or query failure | `"failed to get Part2 prompt: [error details]"` |

### Error Response Example

```json
{
  "error": "failed to get Part2 prompt: no prompts available"
}
```

## Testing

### Test 1: Get Latest Prompt

**Request:**
```json
{
  "name": "part2_prompt",
  "arguments": {}
}
```

**Expected:** Returns the prompt with the highest ID (most recent)

### Test 2: No Prompts Available

**Scenario:** Database has no prompt records

**Request:**
```json
{
  "name": "part2_prompt",
  "arguments": {}
}
```

**Expected:** Error message "failed to get Part2 prompt: no prompts available"

### Test 3: Multiple Prompts

**Scenario:** Database has prompts with IDs: 1, 2, 5, 10

**Request:**
```json
{
  "name": "part2_prompt",
  "arguments": {}
}
```

**Expected:** Returns the prompt with ID=10 (highest ID)

## Benefits

### 1. Centralized Prompt Management
- Prompts stored in database
- Easy to update without code changes
- Version control through database records

### 2. Consistent AI Behavior
- All AI assistants use the latest prompt
- Ensures consistent user experience
- Standardized tool usage patterns

### 3. Dynamic Updates
- Update prompts without redeploying
- New prompts automatically used
- Instant rollout of new instructions

### 4. Automatic Version Management
- Always uses the latest version
- No manual version selection needed
- Simplified prompt management

## AI Assistant Integration

### Recommended Usage Pattern

```
1. On initialization:
   - Call part2_prompt to get instructions
   - Parse and understand available tools
   - Prepare for user interaction

2. During practice session:
   - Follow prompt guidelines
   - Use tools as instructed
   - Provide feedback per prompt specifications

3. On updates:
   - Periodically refresh prompt
   - Adapt to new instructions
   - Maintain consistency
```

### Example AI Workflow

```python
# Initialize
prompt = call_tool("part2_prompt")
parse_instructions(prompt)

# Practice session
while practicing:
    question = get_next_question()
    user_answer = get_user_input()
    result = submit_answer(question, user_answer)
    provide_feedback(result)
```

## Future Enhancements

### Potential Improvements

1. **Prompt Versioning**: Track prompt versions and changes
2. **A/B Testing**: Support multiple active prompts for testing
3. **Localization**: Multi-language prompt support
4. **Personalization**: User-specific prompt customization
5. **Analytics**: Track which prompts perform better
6. **Templates**: Prompt templates with variable substitution

## Related Tools

- **get_user_info**: Get user context for personalization
- **list_part2_sets**: List available practice sets
- **auto_generate_part2_set**: Generate practice content
- **get_part2_set_details**: View practice questions
- **submit_part2_answer**: Submit and evaluate answers

## Conclusion

The `part2_prompt` tool provides a simple and maintainable way to deliver AI assistant instructions. By automatically retrieving the latest prompt from the database, the system ensures all AI assistants use the most up-to-date instructions without requiring any configuration or version management.

