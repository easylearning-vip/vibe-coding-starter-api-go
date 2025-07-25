#!/bin/bash

# Vibe Coding Starter 快速启动脚本
# 适用于已安装 Docker 和 k3d 的系统，快速创建开发环境

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

# 显示横幅
show_banner() {
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "    Vibe Coding Starter 快速启动脚本"
    echo "=================================================================="
    echo -e "${NC}"
    echo "本脚本将快速创建和启动开发环境："
    echo "  • 创建 k3d 集群"
    echo "  • 部署 MySQL 和 Redis 服务"
    echo "  • 配置开发环境"
    echo
    echo "前提条件：已安装 Docker、k3d、kubectl"
    echo "预计启动时间：2-3 分钟"
    echo
}

# 检查依赖
check_dependencies() {
    log_step "检查依赖..."
    
    local missing_deps=()
    
    if ! command -v docker &> /dev/null; then
        missing_deps+=("docker")
    fi
    
    if ! command -v k3d &> /dev/null; then
        missing_deps+=("k3d")
    fi
    
    if ! command -v kubectl &> /dev/null; then
        missing_deps+=("kubectl")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "缺少依赖: ${missing_deps[*]}"
        log_info "请先运行完整安装脚本: ./setup-dev-environment.sh"
        exit 1
    fi
    
    # 检查 Docker 是否运行
    if ! docker info &> /dev/null && ! sudo docker info &> /dev/null; then
        log_error "Docker 服务未运行，请启动 Docker"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 创建存储目录
create_storage_directories() {
    log_step "创建存储目录..."
    
    mkdir -p ~/.local/share/k3d/vibe-dev-storage
    chmod 755 ~/.local/share/k3d/vibe-dev-storage
    
    log_success "存储目录创建完成"
}

# 创建或启动 k3d 集群
setup_k3d_cluster() {
    log_step "设置 k3d 集群..."
    
    # 检查集群是否已存在
    if k3d cluster list | grep -q "vibe-dev"; then
        log_info "集群 'vibe-dev' 已存在"
        
        # 检查集群是否运行
        if k3d cluster list | grep "vibe-dev" | grep -q "1/1.*2/2"; then
            log_info "集群已在运行"
        else
            log_info "启动现有集群..."
            k3d cluster start vibe-dev
        fi
    else
        log_info "创建新集群..."
        
        # 进入 k3d 配置目录
        SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        K3D_DIR="$SCRIPT_DIR/../k3d"
        
        if [[ ! -f "$K3D_DIR/k3d-cluster.yaml" ]]; then
            log_error "找不到 k3d 配置文件: $K3D_DIR/k3d-cluster.yaml"
            exit 1
        fi
        
        cd "$K3D_DIR"

        # 创建临时配置文件，替换环境变量
        log_info "准备集群配置..."
        envsubst < k3d-cluster.yaml > k3d-cluster-temp.yaml

        # 使用临时配置文件创建集群
        k3d cluster create --config k3d-cluster-temp.yaml

        # 清理临时文件
        rm -f k3d-cluster-temp.yaml
    fi
    
    # 等待集群就绪
    log_info "等待集群就绪..."
    sleep 10
    
    # 验证集群
    kubectl cluster-info
    kubectl get nodes
    
    log_success "k3d 集群设置完成"
}

# 部署服务
deploy_services() {
    log_step "部署开发服务..."
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    K3D_DIR="$SCRIPT_DIR/../k3d"
    MANIFESTS_DIR="$K3D_DIR/manifests"
    
    cd "$K3D_DIR"
    
    # 检查服务是否已部署
    if kubectl get namespace vibe-dev &> /dev/null; then
        log_info "命名空间 'vibe-dev' 已存在"
    else
        log_info "部署命名空间..."
        kubectl apply -f manifests/namespace.yaml
    fi
    
    # 部署或更新服务
    log_info "部署 MySQL 服务..."
    kubectl apply -f manifests/mysql.yaml
    
    log_info "部署 Redis 服务..."
    kubectl apply -f manifests/redis.yaml
    
    # 等待服务就绪
    log_info "等待服务就绪..."
    kubectl wait --for=condition=ready pod --all -n vibe-dev --timeout=300s
    
    log_success "开发服务部署完成"
}

# 验证服务
verify_services() {
    log_step "验证服务状态..."
    
    echo
    log_info "Pod 状态："
    kubectl get pods -n vibe-dev
    
    echo
    log_info "服务状态："
    kubectl get svc -n vibe-dev
    
    echo
    log_info "测试服务连接..."
    
    # 测试 MySQL
    if kubectl exec -n vibe-dev mysql-0 -- mysqladmin ping -h localhost -u root -prootpassword &> /dev/null; then
        log_success "MySQL 服务正常"
    else
        log_warning "MySQL 服务可能还在启动中"
    fi

    # 测试 Redis
    if kubectl exec -n vibe-dev redis-0 -- redis-cli ping &> /dev/null; then
        log_success "Redis 服务正常"
    else
        log_warning "Redis 服务可能还在启动中"
    fi
    
    log_success "服务验证完成"
}

# 显示连接信息
show_connection_info() {
    echo
    echo -e "${PURPLE}=================================================================="
    echo "                    开发环境就绪"
    echo -e "==================================================================${NC}"
    echo
    
    log_success "开发环境已启动！"
    echo
    
    echo "数据库连接信息："
    echo "  MySQL:"
    echo "    Host: localhost"
    echo "    Port: 3306"
    echo "    Database: vibe_coding_starter"
    echo "    Username: vibe_user"
    echo "    Password: vibe_password"
    echo
    echo "  Redis:"
    echo "    Host: localhost"
    echo "    Port: 6379"
    echo "    Password: (无)"
    echo
    
    echo "快速连接命令："
    echo "  mysql -h localhost -P 3306 -u vibe_user -pvibe_password vibe_coding_starter"
    echo "  redis-cli -h localhost -p 6379"
    echo
    
    echo "常用管理命令："
    echo "  kubectl get all -n vibe-dev          # 查看所有服务"
    echo "  kubectl logs -f statefulset/mysql -n vibe-dev   # MySQL 日志"
    echo "  kubectl logs -f statefulset/redis -n vibe-dev   # Redis 日志"
    echo "  k3d cluster stop vibe-dev            # 停止集群"
    echo "  k3d cluster start vibe-dev           # 启动集群"
    echo
    
    echo -e "${GREEN}🚀 开发环境已就绪！可以开始开发了！${NC}"
    echo
}

# 主函数
main() {
    show_banner
    
    # 确认启动
    echo -e "${YELLOW}是否继续启动开发环境？${NC}"
    read -p "请输入 y/Y 继续，或按任意键取消: " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "启动已取消"
        exit 0
    fi
    
    echo
    log_info "开始启动 Vibe Coding Starter 开发环境..."
    echo
    
    # 执行启动步骤
    check_dependencies
    create_storage_directories
    setup_k3d_cluster
    deploy_services
    verify_services
    show_connection_info
    
    log_success "开发环境启动完成！"
}

# 错误处理
handle_error() {
    local exit_code=$?
    log_error "启动过程中发生错误 (退出码: $exit_code)"
    echo
    echo "故障排除建议："
    echo "1. 检查 Docker 是否运行"
    echo "2. 检查网络连接"
    echo "3. 确保端口 3306 和 6379 未被占用"
    echo "4. 尝试删除现有集群: k3d cluster delete vibe-dev"
    exit $exit_code
}

# 设置错误处理
trap handle_error ERR

# 运行主函数
main "$@"
