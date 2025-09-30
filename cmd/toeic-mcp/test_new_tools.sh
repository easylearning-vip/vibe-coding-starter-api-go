#!/bin/bash

# Test script for new TOEIC MCP tools
# This script demonstrates how to test the three new tools using curl

# Configuration
MCP_SERVER="http://localhost:6275"
AUTH_TOKEN="9f9f39514d4936a9f58df89584370f5d70b133ee321e7795a0ab430f4917843b"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "TOEIC MCP Server - New Tools Test Script"
echo "=========================================="
echo ""

# Function to make MCP tool call
call_tool() {
    local tool_name=$1
    local arguments=$2
    
    echo -e "${YELLOW}Testing: ${tool_name}${NC}"
    echo "Arguments: ${arguments}"
    echo ""
    
    # Note: This is a simplified example. Actual MCP protocol requires proper JSON-RPC format
    # For real testing, use the MCP Inspector web interface
    
    curl -s -X POST "${MCP_SERVER}" \
        -H "Content-Type: application/json" \
        -H "X-User-Token: ${AUTH_TOKEN}" \
        -d "{
            \"jsonrpc\": \"2.0\",
            \"id\": 1,
            \"method\": \"tools/call\",
            \"params\": {
                \"name\": \"${tool_name}\",
                \"arguments\": ${arguments}
            }
        }" | jq '.'
    
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Test 1: Auto-generate a practice set
echo -e "${GREEN}Test 1: Auto-generate Practice Set${NC}"
echo "This will create a new practice set with 10 questions"
echo ""

call_tool "auto_generate_part2_set" '{
    "count": 10
}'

# Prompt user for set_id
echo -e "${YELLOW}Please enter the set_id from the response above:${NC}"
read -p "Set ID: " SET_ID

if [ -z "$SET_ID" ]; then
    echo -e "${RED}Error: Set ID is required to continue${NC}"
    exit 1
fi

echo ""

# Test 2: Get practice set details
echo -e "${GREEN}Test 2: Get Practice Set Details${NC}"
echo "This will retrieve all questions in the practice set"
echo ""

call_tool "get_part2_set_details" "{
    \"set_id\": ${SET_ID}
}"

# Prompt user for item_id
echo -e "${YELLOW}Please enter an item_id from the response above:${NC}"
read -p "Item ID: " ITEM_ID

if [ -z "$ITEM_ID" ]; then
    echo -e "${RED}Error: Item ID is required to continue${NC}"
    exit 1
fi

echo ""

# Test 3: Submit an answer
echo -e "${GREEN}Test 3: Submit Answer${NC}"
echo "This will submit an answer for the selected question"
echo ""

echo -e "${YELLOW}Please enter your answer (A, B, or C):${NC}"
read -p "Answer: " ANSWER

if [ -z "$ANSWER" ]; then
    echo -e "${RED}Error: Answer is required${NC}"
    exit 1
fi

call_tool "submit_part2_answer" "{
    \"set_id\": ${SET_ID},
    \"item_id\": ${ITEM_ID},
    \"user_answer\": \"${ANSWER}\"
}"

# Test 4: Check updated progress
echo -e "${GREEN}Test 4: Check Updated Progress${NC}"
echo "This will show the updated practice set statistics"
echo ""

call_tool "get_part2_set_details" "{
    \"set_id\": ${SET_ID}
}"

echo ""
echo "=========================================="
echo "Testing Complete!"
echo "=========================================="
echo ""
echo "Summary:"
echo "- Created practice set with ID: ${SET_ID}"
echo "- Submitted answer for item ID: ${ITEM_ID}"
echo "- Answer submitted: ${ANSWER}"
echo ""
echo "For more comprehensive testing, use the MCP Inspector:"
echo "URL: http://localhost:6274/?MCP_PROXY_AUTH_TOKEN=f871086194aa906d024c43c7032fe0219a77399705fba92b090af9559d3487f9"
echo ""

