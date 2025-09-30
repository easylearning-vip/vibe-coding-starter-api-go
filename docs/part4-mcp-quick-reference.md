# Part 4 MCP Tools - Quick Reference

## Available Tools

### 1. auto_generate_part4_set
Generate a new Part 4 practice set.

```json
{
  "count": 15,                    // Required: 10-30 questions
  "difficulty_level_id": 1,       // Optional: filter by difficulty
  "scenario_id": 2                // Optional: filter by scenario
}
```

**Response:**
```
Practice set generated successfully

Set ID: 123
Total Questions: 15
Scenario ID: 2
Difficulty Level ID: 1

Added 15 questions to practice set

Use get_part4_set_details to view question details
```

---

### 2. list_part4_sets
List all practice sets for the current user.

```json
{
  "page": 1,                      // Optional: default 1
  "page_size": 20                 // Optional: default 20, max 100
}
```

**Response:**
```
total=5 page=1 size=20
#123 total=15 done=10 correct=8 acc=80.00%
#122 total=20 done=20 correct=16 acc=80.00%
#121 total=10 done=5 correct=4 acc=80.00%
```

---

### 3. get_part4_set_details
View detailed information about a practice set.

```json
{
  "set_id": 123                   // Required: practice set ID
}
```

**Response:**
```
Practice Set Details

Set ID: 123
Total Questions: 15
Completed: 10
Correct: 8
Accuracy: 80.00%

Question List:

(Item ID: 456) Talk 1:
Title: Office Announcement
Content: Good morning, everyone. This is a reminder that...

Question 1: What is the main purpose of this announcement?
A. To announce a meeting
B. To remind about a deadline
C. To introduce a new policy
D. To cancel an event

(Item ID: 457) Talk 2:
...
```

---

### 4. submit_part4_answer
Submit an answer for a practice question.

```json
{
  "set_id": 123,                  // Required: practice set ID
  "item_id": 456,                 // Required: item ID from get_part4_set_details
  "user_answer": "B",             // Required: A, B, C, or D
  "is_correct": true              // Optional: AI can judge correctness
}
```

**Response (Correct):**
```
Correct answer!

Your answer: B
Correct answer: B
```

**Response (Incorrect):**
```
Incorrect answer

Your answer: A
Correct answer: B
```

---

### 5. part4_prompt
Get the AI assistant prompt for Part 4 practice.

```json
{}
```

**Response:**
Returns the latest Part 4 AI prompt text with instructions for effective practice.

---

## Typical Practice Flow

### Step 1: Generate a Practice Set
```json
{
  "tool": "auto_generate_part4_set",
  "params": {
    "count": 15,
    "difficulty_level_id": 1
  }
}
```
→ Note the `Set ID` from the response

### Step 2: Get Practice Set Details
```json
{
  "tool": "get_part4_set_details",
  "params": {
    "set_id": 123
  }
}
```
→ Note the `Item ID` for each question

### Step 3: Practice Questions
For each question:
1. Read the talk content
2. Read the specific question
3. Review the answer options
4. Submit your answer:

```json
{
  "tool": "submit_part4_answer",
  "params": {
    "set_id": 123,
    "item_id": 456,
    "user_answer": "B"
  }
}
```

### Step 4: Review Progress
```json
{
  "tool": "list_part4_sets",
  "params": {
    "page": 1,
    "page_size": 10
  }
}
```

---

## Error Handling

### Common Errors

**Invalid count:**
```
parameter error: count must be between 10 and 30
```

**Missing set_id:**
```
parameter error: set_id is required
```

**Invalid answer format:**
```
parameter error: user_answer must be A, B, C, or D
```

**Item doesn't belong to set:**
```
question does not belong to the specified practice set
```

---

## Tips for AI Assistants

### 1. Reading Questions Aloud
When simulating audio practice:
- Read the talk content first
- Then read the specific question
- Present all four answer options
- Wait for user response

### 2. Judging Answers
Use the `is_correct` parameter when:
- You want to provide nuanced feedback
- The database answer might be ambiguous
- You're teaching reasoning skills

### 3. Progress Tracking
- Check `list_part4_sets` regularly to show progress
- Celebrate milestones (50%, 75%, 100% completion)
- Point out accuracy trends

### 4. Adaptive Practice
- Start with easier difficulty levels
- Increase difficulty as accuracy improves
- Focus on specific scenarios where user struggles

---

## Database Schema Reference

### Practice Set Structure
```
Part4PracticeSet (Master)
├── user_id
├── scenario_id (optional)
├── difficulty_level_id (optional)
├── total_questions
├── completed_count
├── correct_count
└── accuracy

Part4PracticeSetItem (Detail)
├── set_id → Part4PracticeSet
├── talk_id → Part4Talk
├── question_index (1, 2, or 3)
├── order_index
├── selected_answer (A/B/C/D)
└── is_correct

Part4Talk
├── test_id
├── talk_number
├── title
├── content
├── question1, question2, question3
├── scenario_id
├── difficulty_level_id
└── answer_options[] → Part4AnswerOption

Part4AnswerOption
├── talk_id
├── question_number (1, 2, or 3)
├── option_a, option_b, option_c, option_d
└── correct_answer (A/B/C/D)
```

---

## Integration with Claude Desktop

Add to your Claude Desktop MCP configuration:

```json
{
  "mcpServers": {
    "toeic-part4": {
      "command": "/path/to/toeic-mcp",
      "args": ["--host", "0.0.0.0", "--port", "6275"],
      "env": {
        "CONFIG_FILE": "/path/to/config.yaml"
      }
    }
  }
}
```

Then use the tools in your conversations:
- "Generate a Part 4 practice set with 15 questions"
- "Show me my Part 4 practice sets"
- "Let's practice Part 4 questions from set 123"
- "Submit my answer: B"

---

## Troubleshooting

### Build Issues
```bash
cd workspace/vocabulary/vibe-coding-starter-api-go
go build -o toeic-mcp ./cmd/toeic-mcp/
```

### Run Server
```bash
./toeic-mcp --host 0.0.0.0 --port 6275 -c config/dev.yaml
```

### Check Logs
Look for:
- "Starting TOEIC MCP server" - Server started successfully
- "invalid user token" - Authentication issues
- "generation failed" - Database or service layer issues

### Database Connection
Ensure:
- MySQL is running
- Database credentials are correct in config
- Tables exist (run migrations if needed)
- Part4 data is imported

---

## Support

For issues or questions:
1. Check the main documentation: `docs/part4-mcp-implementation.md`
2. Review existing Part 2/Part 3 implementations for patterns
3. Check service layer logs for detailed error messages
4. Verify database has Part 4 content loaded

