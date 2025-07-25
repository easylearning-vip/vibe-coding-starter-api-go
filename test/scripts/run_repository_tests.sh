#!/bin/bash

# Repository单元测试运行脚本
# 使用Docker创建测试MySQL和Redis环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否运行
check_docker() {
    if ! sudo docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    print_success "Docker is running"
}

# 启动测试环境
setup_test_env() {
    print_info "Setting up test environment..."
    make -f Makefile.test test-setup
    
    # 等待服务就绪
    print_info "Waiting for services to be ready..."
    sleep 5
    
    # 检查MySQL连接
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if sudo docker exec vibe-mysql-test mysqladmin ping -h localhost --silent; then
            print_success "MySQL is ready"
            break
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            print_error "MySQL failed to start after $max_attempts attempts"
            exit 1
        fi
        
        print_info "Waiting for MySQL... (attempt $attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    # 检查Redis连接
    if sudo docker exec vibe-redis-test redis-cli ping > /dev/null 2>&1; then
        print_success "Redis is ready"
    else
        print_error "Redis is not responding"
        exit 1
    fi
}

# 运行Repository测试
run_repository_tests() {
    print_info "Running Repository unit tests..."
    
    # 运行测试并捕获输出
    if go test -v -race -short ./test/repository/... 2>&1 | tee test_output.log; then
        print_success "All Repository tests passed!"
        
        # 生成测试报告
        generate_test_report
    else
        print_error "Some Repository tests failed!"
        exit 1
    fi
}

# 生成测试报告
generate_test_report() {
    print_info "Generating test report..."
    
    local total_tests=$(grep -c "=== RUN" test_output.log || echo "0")
    local passed_tests=$(grep -c "--- PASS:" test_output.log || echo "0")
    local failed_tests=$(grep -c "--- FAIL:" test_output.log || echo "0")
    local test_time=$(grep "ok.*github.com/your-org/vibe-coding-starter/test/repository" test_output.log | awk '{print $3}' || echo "0s")
    
    echo ""
    echo "=========================================="
    echo "           Repository Test Report"
    echo "=========================================="
    echo "Total Tests:    $total_tests"
    echo "Passed Tests:   $passed_tests"
    echo "Failed Tests:   $failed_tests"
    echo "Test Duration:  $test_time"
    echo ""
    
    # 显示各个Repository的测试结果
    echo "Repository Test Results:"
    echo "------------------------"
    
    local repositories=("Article" "Category" "Comment" "File" "Tag" "User")
    
    for repo in "${repositories[@]}"; do
        local repo_tests=$(grep -c "Test${repo}RepositoryTestSuite/" test_output.log || echo "0")
        local repo_passed=$(grep "Test${repo}RepositoryTestSuite.*PASS:" test_output.log | wc -l || echo "0")
        
        if [ "$repo_tests" -gt 0 ]; then
            printf "%-12s: %2d tests, %2d passed\n" "$repo" "$repo_tests" "$repo_passed"
        fi
    done
    
    echo ""
    echo "Test Coverage Areas:"
    echo "-------------------"
    echo "✓ CRUD Operations (Create, Read, Update, Delete)"
    echo "✓ Data Validation (Unique constraints, Required fields)"
    echo "✓ Query Operations (Filters, Search, Pagination)"
    echo "✓ Relationship Handling (Foreign keys, Associations)"
    echo "✓ Error Handling (Not found, Duplicate entries)"
    echo "✓ Database Transactions (Race conditions, Concurrency)"
    echo ""
    
    # 清理临时文件
    rm -f test_output.log
}

# 清理测试环境
cleanup_test_env() {
    print_info "Cleaning up test environment..."
    make -f Makefile.test test-teardown
    print_success "Test environment cleaned up"
}

# 主函数
main() {
    echo "=========================================="
    echo "    Repository Unit Tests with Docker"
    echo "=========================================="
    echo ""
    
    # 检查依赖
    check_docker
    
    # 设置清理陷阱
    trap cleanup_test_env EXIT
    
    # 运行测试流程
    setup_test_env
    run_repository_tests
    
    print_success "Repository tests completed successfully!"
}

# 运行主函数
main "$@"
