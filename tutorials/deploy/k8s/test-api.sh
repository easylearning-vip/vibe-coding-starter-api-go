#!/bin/bash

# Vibe Coding Starter API - 测试脚本
# 用于验证 k8s 部署的 API 服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
API_BASE_URL="http://api.vibe-dev.com:8000"
TIMEOUT=30

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 测试 HTTP 请求
test_http() {
    local url=$1
    local expected_status=${2:-200}
    local description=$3
    
    log_info "测试: $description"
    log_info "URL: $url"
    
    local response
    local status_code
    
    response=$(curl -s -w "\n%{http_code}" --max-time $TIMEOUT "$url" 2>/dev/null || echo -e "\n000")
    status_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n -1)
    
    if [[ "$status_code" == "$expected_status" ]]; then
        log_success "✓ 状态码: $status_code"
        if [[ -n "$body" && "$body" != "null" ]]; then
            echo "响应内容: $body"
        fi
        echo
        return 0
    else
        log_error "✗ 状态码: $status_code (期望: $expected_status)"
        if [[ -n "$body" ]]; then
            echo "响应内容: $body"
        fi
        echo
        return 1
    fi
}

# 测试 JSON API
test_json_api() {
    local url=$1
    local expected_status=${2:-200}
    local description=$3
    
    log_info "测试: $description"
    log_info "URL: $url"
    
    local response
    local status_code
    
    response=$(curl -s -w "\n%{http_code}" --max-time $TIMEOUT \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        "$url" 2>/dev/null || echo -e "\n000")
    
    status_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n -1)
    
    if [[ "$status_code" == "$expected_status" ]]; then
        log_success "✓ 状态码: $status_code"
        if [[ -n "$body" && "$body" != "null" ]]; then
            echo "JSON 响应:"
            echo "$body" | python3 -m json.tool 2>/dev/null || echo "$body"
        fi
        echo
        return 0
    else
        log_error "✗ 状态码: $status_code (期望: $expected_status)"
        if [[ -n "$body" ]]; then
            echo "响应内容: $body"
        fi
        echo
        return 1
    fi
}

# 检查前置条件
check_prerequisites() {
    log_info "检查前置条件..."
    
    # 检查 curl
    if ! command -v curl &> /dev/null; then
        log_error "curl 未安装"
        exit 1
    fi
    
    # 检查域名解析
    if ! grep -q "api.vibe-dev.com" /etc/hosts; then
        log_warning "hosts 文件中未找到 api.vibe-dev.com"
        log_info "请运行: echo '127.0.0.1 api.vibe-dev.com' | sudo tee -a /etc/hosts"
    fi
    
    # 检查 k8s 部署状态
    if command -v kubectl &> /dev/null; then
        if kubectl get deployment vibe-api-deployment -n vibe-dev &> /dev/null; then
            local ready_replicas=$(kubectl get deployment vibe-api-deployment -n vibe-dev -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
            local desired_replicas=$(kubectl get deployment vibe-api-deployment -n vibe-dev -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")
            
            if [[ "$ready_replicas" == "$desired_replicas" && "$ready_replicas" != "0" ]]; then
                log_success "k8s 部署状态正常 ($ready_replicas/$desired_replicas)"
            else
                log_warning "k8s 部署状态异常 ($ready_replicas/$desired_replicas)"
            fi
        else
            log_warning "未找到 k8s 部署"
        fi
    fi
    
    log_success "前置条件检查完成"
    echo
}

# 基础连通性测试
test_connectivity() {
    log_info "=== 基础连通性测试 ==="
    
    # 测试健康检查端点
    test_http "$API_BASE_URL/health" 200 "健康检查端点"
    
    # 测试根路径
    test_http "$API_BASE_URL/" 404 "根路径访问"
}

# API 端点测试
test_api_endpoints() {
    log_info "=== API 端点测试 ==="
    
    # 测试 API 健康检查
    test_json_api "$API_BASE_URL/api/v1/health" 200 "API 健康检查"
    
    # 测试用户相关端点（可能需要认证）
    test_json_api "$API_BASE_URL/api/v1/users" 401 "用户列表（未认证）"
    
    # 测试认证端点
    test_json_api "$API_BASE_URL/api/v1/auth/login" 400 "登录端点（无参数）"
    
    # 测试注册端点
    test_json_api "$API_BASE_URL/api/v1/auth/register" 400 "注册端点（无参数）"
}

# 监控端点测试
test_monitoring_endpoints() {
    log_info "=== 监控端点测试 ==="
    
    # 测试指标端点
    test_http "$API_BASE_URL/metrics" 200 "Prometheus 指标"
    
    # 测试健康检查详情
    test_json_api "$API_BASE_URL/health" 200 "健康检查详情"
}

# 性能测试
test_performance() {
    log_info "=== 性能测试 ==="
    
    log_info "测试响应时间..."
    local start_time=$(date +%s%N)
    
    if curl -s --max-time $TIMEOUT "$API_BASE_URL/health" > /dev/null; then
        local end_time=$(date +%s%N)
        local duration=$(( (end_time - start_time) / 1000000 ))
        
        if [[ $duration -lt 1000 ]]; then
            log_success "✓ 响应时间: ${duration}ms (优秀)"
        elif [[ $duration -lt 3000 ]]; then
            log_success "✓ 响应时间: ${duration}ms (良好)"
        else
            log_warning "⚠ 响应时间: ${duration}ms (较慢)"
        fi
    else
        log_error "✗ 性能测试失败"
    fi
    
    echo
}

# 错误处理测试
test_error_handling() {
    log_info "=== 错误处理测试 ==="
    
    # 测试不存在的端点
    test_http "$API_BASE_URL/api/v1/nonexistent" 404 "不存在的端点"
    
    # 测试错误的方法
    log_info "测试: 错误的 HTTP 方法"
    local response=$(curl -s -w "\n%{http_code}" --max-time $TIMEOUT -X DELETE "$API_BASE_URL/health" 2>/dev/null || echo -e "\n000")
    local status_code=$(echo "$response" | tail -n1)
    
    if [[ "$status_code" == "405" || "$status_code" == "404" ]]; then
        log_success "✓ 正确处理错误的 HTTP 方法: $status_code"
    else
        log_warning "⚠ HTTP 方法处理: $status_code"
    fi
    echo
}

# 生成测试报告
generate_report() {
    local total_tests=$1
    local passed_tests=$2
    local failed_tests=$((total_tests - passed_tests))
    
    echo
    log_info "=== 测试报告 ==="
    echo "总测试数: $total_tests"
    echo "通过: $passed_tests"
    echo "失败: $failed_tests"
    
    if [[ $failed_tests -eq 0 ]]; then
        log_success "所有测试通过！🎉"
        return 0
    else
        log_warning "有 $failed_tests 个测试失败"
        return 1
    fi
}

# 主函数
main() {
    echo "Vibe Coding Starter API - 测试脚本"
    echo "=================================="
    echo
    
    check_prerequisites
    
    local total_tests=0
    local passed_tests=0
    
    # 运行测试
    test_connectivity && ((passed_tests++)) || true
    ((total_tests++))
    
    test_api_endpoints && ((passed_tests++)) || true
    ((total_tests++))
    
    test_monitoring_endpoints && ((passed_tests++)) || true
    ((total_tests++))
    
    test_performance && ((passed_tests++)) || true
    ((total_tests++))
    
    test_error_handling && ((passed_tests++)) || true
    ((total_tests++))
    
    # 生成报告
    generate_report $total_tests $passed_tests
}

# 执行主函数
main "$@"
