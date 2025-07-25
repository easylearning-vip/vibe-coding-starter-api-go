#!/bin/bash

# Vibe Coding Starter 环境清理脚本
# 用于停止和清理开发环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

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

log_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Vibe Coding Starter 环境清理脚本"
    echo
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  stop      停止开发环境（保留集群和数据）"
    echo "  clean     删除集群但保留工具和配置"
    echo "  reset     重置环境（删除集群和数据，保留工具）"
    echo "  purge     完全清理（删除所有组件和数据）"
    echo "  help      显示此帮助信息"
    echo
    echo "示例:"
    echo "  $0 stop     # 停止服务但保留数据"
    echo "  $0 clean    # 删除集群"
    echo "  $0 reset    # 重置环境"
    echo "  $0 purge    # 完全清理"
}

# 检查依赖
check_dependencies() {
    local missing_deps=()
    
    if ! command -v k3d &> /dev/null; then
        missing_deps+=("k3d")
    fi
    
    if ! command -v kubectl &> /dev/null; then
        missing_deps+=("kubectl")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_warning "缺少工具: ${missing_deps[*]}"
        log_info "某些清理操作可能无法执行"
    fi
}

# 停止开发环境
stop_environment() {
    log_step "停止开发环境..."
    
    if ! command -v k3d &> /dev/null; then
        log_error "k3d 未安装，无法停止集群"
        return 1
    fi
    
    # 检查集群是否存在
    if ! k3d cluster list | grep -q "vibe-dev"; then
        log_info "集群 'vibe-dev' 不存在"
        return 0
    fi
    
    # 停止集群
    log_info "停止 k3d 集群..."
    k3d cluster stop vibe-dev
    
    log_success "开发环境已停止"
    log_info "数据已保留，使用 'k3d cluster start vibe-dev' 可重新启动"
}

# 清理集群
clean_cluster() {
    log_step "清理 k3d 集群..."
    
    if ! command -v k3d &> /dev/null; then
        log_error "k3d 未安装，无法删除集群"
        return 1
    fi
    
    # 检查集群是否存在
    if ! k3d cluster list | grep -q "vibe-dev"; then
        log_info "集群 'vibe-dev' 不存在"
    else
        log_info "删除 k3d 集群..."
        k3d cluster delete vibe-dev
        log_success "k3d 集群已删除"
    fi
    
    # 清理存储目录
    if [[ -d ~/.local/share/k3d/vibe-dev-storage ]]; then
        log_info "清理存储目录..."
        rm -rf ~/.local/share/k3d/vibe-dev-storage
        log_success "存储目录已清理"
    fi
    
    # 清理 kubeconfig
    if command -v kubectl &> /dev/null; then
        log_info "清理 kubeconfig..."
        kubectl config delete-context k3d-vibe-dev &> /dev/null || true
        kubectl config delete-cluster k3d-vibe-dev &> /dev/null || true
        kubectl config delete-user admin@k3d-vibe-dev &> /dev/null || true
        log_success "kubeconfig 已清理"
    fi
}

# 重置环境
reset_environment() {
    log_step "重置开发环境..."
    
    # 先清理集群
    clean_cluster
    
    # 清理 Docker 镜像（可选）
    if command -v docker &> /dev/null; then
        log_info "清理相关 Docker 镜像..."
        
        # 清理 k3s 镜像
        docker images | grep "rancher/k3s" | awk '{print $3}' | xargs -r docker rmi &> /dev/null || true
        
        # 清理数据库镜像（可选，用户可能在其他地方使用）
        read -p "是否清理 MySQL 和 Redis 镜像？(y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker images | grep -E "(mysql|redis)" | awk '{print $3}' | xargs -r docker rmi &> /dev/null || true
            log_info "数据库镜像已清理"
        fi
        
        # 清理未使用的资源
        docker system prune -f &> /dev/null || true
        log_success "Docker 资源已清理"
    fi
    
    log_success "环境重置完成"
}

# 完全清理
purge_all() {
    log_step "完全清理所有组件..."
    
    log_warning "这将删除所有相关组件和数据，包括："
    log_warning "- k3d 集群和数据"
    log_warning "- k3d 工具"
    log_warning "- kubectl 工具"
    log_warning "- Docker 镜像"
    log_warning "- 配置文件和别名"
    echo
    
    read -p "确定要继续吗？这个操作不可逆！(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "清理已取消"
        return 0
    fi
    
    # 先重置环境
    reset_environment
    
    # 卸载 k3d
    if command -v k3d &> /dev/null; then
        log_info "卸载 k3d..."
        sudo rm -f /usr/local/bin/k3d
        log_success "k3d 已卸载"
    fi
    
    # 卸载 kubectl
    if command -v kubectl &> /dev/null; then
        log_info "卸载 kubectl..."
        sudo rm -f /usr/local/bin/kubectl
        log_success "kubectl 已卸载"
    fi
    
    # 清理配置目录
    log_info "清理配置目录..."
    rm -rf ~/.local/share/k3d
    rm -rf ~/.kube
    
    # 清理 shell 配置中的别名
    log_info "清理 shell 别名..."
    for profile in ~/.bashrc ~/.zshrc; do
        if [[ -f "$profile" ]]; then
            # 备份原文件
            cp "$profile" "${profile}.backup.$(date +%Y%m%d_%H%M%S)"
            
            # 删除 Vibe Coding Starter 相关的别名
            sed -i '/# Vibe Coding Starter aliases/,/^$/d' "$profile" 2>/dev/null || true
            sed -i '/# Load kubectl, k3d/,/compinit/d' "$profile" 2>/dev/null || true
        fi
    done
    
    # 清理补全目录
    rm -rf ~/.bash_completion.d/kubectl ~/.bash_completion.d/k3d
    rm -rf ~/.zsh/completions/kubectl ~/.zsh/completions/k3d
    
    log_success "完全清理完成"
    log_info "请重新加载 shell 或重新登录以使更改生效"
}

# 显示状态
show_status() {
    log_step "检查当前状态..."
    
    echo "工具安装状态："
    if command -v docker &> /dev/null; then
        echo "  ✓ Docker: $(docker --version 2>/dev/null || echo '已安装但无法获取版本')"
    else
        echo "  ✗ Docker: 未安装"
    fi
    
    if command -v k3d &> /dev/null; then
        echo "  ✓ k3d: $(k3d version 2>/dev/null || echo '已安装但无法获取版本')"
    else
        echo "  ✗ k3d: 未安装"
    fi
    
    if command -v kubectl &> /dev/null; then
        echo "  ✓ kubectl: $(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null || echo '已安装但无法获取版本')"
    else
        echo "  ✗ kubectl: 未安装"
    fi
    
    echo
    echo "集群状态："
    if command -v k3d &> /dev/null; then
        if k3d cluster list | grep -q "vibe-dev"; then
            echo "  ✓ k3d 集群 'vibe-dev' 存在"
            k3d cluster list | grep "vibe-dev"
        else
            echo "  ✗ k3d 集群 'vibe-dev' 不存在"
        fi
    else
        echo "  ? 无法检查（k3d 未安装）"
    fi
    
    echo
    echo "服务状态："
    if command -v kubectl &> /dev/null && kubectl cluster-info &> /dev/null; then
        if kubectl get namespace vibe-dev &> /dev/null; then
            echo "  ✓ vibe-dev 命名空间存在"
            kubectl get pods -n vibe-dev 2>/dev/null || echo "  ? 无法获取 Pod 状态"
        else
            echo "  ✗ vibe-dev 命名空间不存在"
        fi
    else
        echo "  ? 无法检查（集群未连接）"
    fi
}

# 主函数
main() {
    case "${1:-help}" in
        "stop")
            check_dependencies
            stop_environment
            ;;
        "clean")
            check_dependencies
            clean_cluster
            ;;
        "reset")
            check_dependencies
            reset_environment
            ;;
        "purge")
            check_dependencies
            purge_all
            ;;
        "status")
            show_status
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# 运行主函数
main "$@"
