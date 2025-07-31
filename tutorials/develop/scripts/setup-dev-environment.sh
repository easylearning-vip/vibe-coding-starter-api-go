#!/bin/bash

# Vibe Coding Starter 开发环境一键安装脚本
# 适用于全新的 Ubuntu 系统，自动安装和配置所有依赖
# 包括：Docker、k3d、kubectl、k3d集群、MySQL、Redis

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
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

log_substep() {
    echo -e "${CYAN}  →${NC} $1"
}

# 显示横幅
show_banner() {
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "    Vibe Coding Starter 开发环境一键安装脚本"
    echo "=================================================================="
    echo -e "${NC}"
    echo "本脚本将为您安装和配置以下组件："
    echo "  • Docker Engine"
    echo "  • k3d (Kubernetes in Docker)"
    echo "  • kubectl (Kubernetes CLI)"
    echo "  • k3d 开发集群"
    echo "  • MySQL 8.0 数据库"
    echo "  • Redis 7 缓存服务"
    echo
    echo "适用系统：Ubuntu 18.04+ / Debian 10+"
    echo "预计安装时间：5-10 分钟"
    echo
}

# 检查系统要求
check_system_requirements() {
    log_step "检查系统要求..."
    
    # 检查操作系统
    if [[ ! -f /etc/os-release ]]; then
        log_error "无法检测操作系统版本"
        exit 1
    fi
    
    . /etc/os-release
    log_substep "操作系统: $PRETTY_NAME"
    
    # 检查是否为 Ubuntu/Debian
    if [[ "$ID" != "ubuntu" && "$ID" != "debian" ]]; then
        log_error "此脚本仅支持 Ubuntu 和 Debian 系统"
        exit 1
    fi
    
    # 检查内存
    MEMORY_GB=$(free -g | awk '/^Mem:/{print $2}')
    if [[ $MEMORY_GB -lt 4 ]]; then
        log_warning "系统内存少于 4GB，可能影响性能"
    fi
    log_substep "可用内存: ${MEMORY_GB}GB"
    
    # 检查磁盘空间
    DISK_SPACE_GB=$(df / | awk 'NR==2{print int($4/1024/1024)}')
    if [[ $DISK_SPACE_GB -lt 10 ]]; then
        log_error "磁盘空间不足，至少需要 10GB 可用空间"
        exit 1
    fi
    log_substep "可用磁盘空间: ${DISK_SPACE_GB}GB"
    
    # 检查网络连接
    if ! ping -c 1 google.com &> /dev/null; then
        log_error "网络连接失败，请检查网络设置"
        exit 1
    fi
    log_substep "网络连接: 正常"
    
    log_success "系统要求检查通过"
}

# 更新系统包
update_system() {
    log_step "更新系统包..."
    
    log_substep "更新包索引..."
    sudo apt update -qq
    
    log_substep "安装基础工具..."
    sudo apt install -y \
        curl \
        wget \
        gnupg \
        lsb-release \
        ca-certificates \
        apt-transport-https \
        software-properties-common \
        unzip \
        git \
        vim \
        htop \
        tree
    
    log_success "系统包更新完成"
}

# 安装 Docker
install_docker() {
    log_step "安装 Docker..."
    
    if command -v docker &> /dev/null; then
        log_substep "Docker 已安装: $(docker --version)"
        return 0
    fi
    
    log_substep "添加 Docker 官方 GPG 密钥..."
    curl -fsSL https://download.docker.com/linux/$ID/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    
    log_substep "添加 Docker 仓库..."
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/$ID $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    log_substep "安装 Docker Engine..."
    sudo apt update -qq
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    log_substep "配置 Docker 服务..."
    sudo systemctl start docker
    sudo systemctl enable docker
    
    log_substep "将当前用户添加到 docker 组..."
    sudo usermod -aG docker $USER
    
    log_success "Docker 安装完成"
}

# 安装 kubectl
install_kubectl() {
    log_step "安装 kubectl..."
    
    if command -v kubectl &> /dev/null; then
        log_substep "kubectl 已安装: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
        return 0
    fi
    
    log_substep "下载 kubectl..."
    KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)
    curl -LO "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl"
    
    log_substep "安装 kubectl..."
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
    
    log_success "kubectl 安装完成"
}

# 安装 k3d
install_k3d() {
    log_step "安装 k3d..."
    
    if command -v k3d &> /dev/null; then
        log_substep "k3d 已安装: $(k3d version)"
        return 0
    fi
    
    log_substep "下载并安装 k3d..."
    curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
    
    log_success "k3d 安装完成"
}

# 验证 Docker 安装
verify_docker() {
    log_step "验证 Docker 安装..."
    
    # 等待 Docker 服务启动
    sleep 3
    
    # 测试 Docker 是否正常工作
    if sudo docker run --rm hello-world &> /dev/null; then
        log_substep "Docker 运行测试通过"
    else
        log_error "Docker 运行测试失败"
        exit 1
    fi
    
    log_success "Docker 验证完成"
}

# 创建存储目录
create_storage_directories() {
    log_step "创建存储目录..."
    
    log_substep "创建 k3d 存储目录..."
    mkdir -p ~/.local/share/k3d/vibe-dev-storage
    chmod 755 ~/.local/share/k3d/vibe-dev-storage
    
    log_success "存储目录创建完成"
}

# 创建 k3d 集群
create_k3d_cluster() {
    log_step "创建 k3d 集群..."
    
    # 检查集群是否已存在
    if k3d cluster list | grep -q "vibe-dev"; then
        log_substep "k3d 集群 'vibe-dev' 已存在"
        return 0
    fi
    
    # 进入 k3d 配置目录
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    K3D_DIR="$SCRIPT_DIR/../k3d"
    
    if [[ ! -f "$K3D_DIR/k3d-cluster.yaml" ]]; then
        log_error "找不到 k3d 配置文件: $K3D_DIR/k3d-cluster.yaml"
        exit 1
    fi
    
    log_substep "使用配置文件创建集群..."
    cd "$K3D_DIR"

    # 创建临时配置文件，替换环境变量
    log_substep "准备集群配置..."
    envsubst < k3d-cluster.yaml > k3d-cluster-temp.yaml

    # 使用临时配置文件创建集群
    k3d cluster create --config k3d-cluster-temp.yaml

    # 清理临时文件
    rm -f k3d-cluster-temp.yaml
    
    log_substep "等待集群就绪..."
    sleep 10
    
    log_substep "验证集群状态..."
    kubectl cluster-info
    kubectl get nodes
    
    log_success "k3d 集群创建完成"
}

# 部署服务
deploy_services() {
    log_step "部署开发服务..."
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    K3D_DIR="$SCRIPT_DIR/../k3d"
    MANIFESTS_DIR="$K3D_DIR/manifests"
    
    if [[ ! -d "$MANIFESTS_DIR" ]]; then
        log_error "找不到 manifests 目录: $MANIFESTS_DIR"
        exit 1
    fi
    
    cd "$K3D_DIR"
    
    log_substep "部署命名空间..."
    kubectl apply -f manifests/namespace.yaml
    
    log_substep "部署 MySQL 服务..."
    kubectl apply -f manifests/mysql.yaml
    
    log_substep "部署 Redis 服务..."
    kubectl apply -f manifests/redis.yaml
    
    log_substep "等待服务就绪..."
    kubectl wait --for=condition=ready pod --all -n vibe-dev --timeout=300s
    
    log_substep "检查服务状态..."
    kubectl get all -n vibe-dev
    
    log_success "开发服务部署完成"
}

# 验证服务
verify_services() {
    log_step "验证服务状态..."

    log_substep "检查 Pod 状态..."
    kubectl get pods -n vibe-dev

    log_substep "检查服务端点..."
    kubectl get svc -n vibe-dev

    log_substep "测试 MySQL 连接..."
    if kubectl exec -n vibe-dev mysql-0 -- mysqladmin ping -h localhost -u root -prootpassword &> /dev/null; then
        log_substep "MySQL 服务正常"
    else
        log_warning "MySQL 服务可能还在启动中"
    fi

    log_substep "测试 Redis 连接..."
    if kubectl exec -n vibe-dev redis-0 -- redis-cli ping &> /dev/null; then
        log_substep "Redis 服务正常"
    else
        log_warning "Redis 服务可能还在启动中"
    fi

    log_success "服务验证完成"
}

# 配置 shell 补全
setup_shell_completion() {
    log_step "配置 shell 补全..."

    # 检测当前 shell
    CURRENT_SHELL=$(basename "$SHELL")

    case $CURRENT_SHELL in
        "bash")
            COMPLETION_DIR="$HOME/.bash_completion.d"
            PROFILE_FILE="$HOME/.bashrc"
            ;;
        "zsh")
            COMPLETION_DIR="$HOME/.zsh/completions"
            PROFILE_FILE="$HOME/.zshrc"
            ;;
        *)
            log_warning "不支持的 shell: $CURRENT_SHELL，跳过补全设置"
            return
            ;;
    esac

    log_substep "创建补全目录..."
    mkdir -p "$COMPLETION_DIR"

    log_substep "生成补全脚本..."
    kubectl completion $CURRENT_SHELL > "$COMPLETION_DIR/kubectl" 2>/dev/null || true
    k3d completion $CURRENT_SHELL > "$COMPLETION_DIR/k3d" 2>/dev/null || true

    # 添加到 profile 文件
    if [[ "$CURRENT_SHELL" == "bash" ]]; then
        if ! grep -q "bash_completion.d" "$PROFILE_FILE" 2>/dev/null; then
            echo "" >> "$PROFILE_FILE"
            echo "# Load kubectl, k3d completions" >> "$PROFILE_FILE"
            echo "for f in ~/.bash_completion.d/*; do [[ -r \$f ]] && source \$f; done" >> "$PROFILE_FILE"
        fi
    elif [[ "$CURRENT_SHELL" == "zsh" ]]; then
        if ! grep -q "zsh/completions" "$PROFILE_FILE" 2>/dev/null; then
            echo "" >> "$PROFILE_FILE"
            echo "# Load kubectl, k3d completions" >> "$PROFILE_FILE"
            echo "fpath=(~/.zsh/completions \$fpath)" >> "$PROFILE_FILE"
            echo "autoload -U compinit && compinit" >> "$PROFILE_FILE"
        fi
    fi

    log_success "Shell 补全配置完成"
}

# 创建便捷别名
create_aliases() {
    log_step "创建便捷别名..."

    CURRENT_SHELL=$(basename "$SHELL")
    case $CURRENT_SHELL in
        "bash")
            PROFILE_FILE="$HOME/.bashrc"
            ;;
        "zsh")
            PROFILE_FILE="$HOME/.zshrc"
            ;;
        *)
            log_warning "不支持的 shell: $CURRENT_SHELL，跳过别名设置"
            return
            ;;
    esac

    # 添加别名到 profile 文件
    if ! grep -q "# Vibe Coding Starter aliases" "$PROFILE_FILE" 2>/dev/null; then
        cat >> "$PROFILE_FILE" << 'EOF'

# Vibe Coding Starter aliases
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgn='kubectl get nodes'
alias kns='kubectl config set-context --current --namespace'
alias vibe-dev='kubectl config set-context --current --namespace=vibe-dev'
alias vibe-logs-mysql='kubectl logs -f statefulset/mysql -n vibe-dev'
alias vibe-logs-redis='kubectl logs -f statefulset/redis -n vibe-dev'
alias vibe-mysql='kubectl exec -it mysql-0 -n vibe-dev -- mysql -u vibe_user -pvibe_password vibe_coding_starter'
alias vibe-redis='kubectl exec -it redis-0 -n vibe-dev -- redis-cli'
alias vibe-status='kubectl get all -n vibe-dev'
EOF
    fi

    log_success "便捷别名创建完成"
}

# 显示安装总结
show_installation_summary() {
    echo
    echo -e "${PURPLE}=================================================================="
    echo "                    安装完成总结"
    echo -e "==================================================================${NC}"
    echo

    log_success "所有组件安装成功！"
    echo
    echo "已安装的组件版本："
    echo "  • Docker: $(docker --version 2>/dev/null || echo '安装失败')"
    echo "  • kubectl: $(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null || echo '安装失败')"
    echo "  • k3d: $(k3d version 2>/dev/null || echo '安装失败')"
    echo

    echo "已部署的服务："
    echo "  • k3d 集群: vibe-dev"
    echo "  • MySQL 8.0: localhost:3306"
    echo "  • Redis 7: localhost:6379"
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
}

# 显示后续步骤
show_next_steps() {
    echo "后续步骤："
    echo "1. 重新加载 shell 配置以启用别名和补全："
    echo "   source ~/.bashrc  # 或 source ~/.zshrc"
    echo
    echo "2. 验证安装："
    echo "   k3d cluster list"
    echo "   kubectl get nodes"
    echo "   vibe-status  # 查看所有服务状态"
    echo
    echo "3. 连接数据库："
    echo "   # 使用别名快速连接"
    echo "   vibe-mysql   # 连接 MySQL"
    echo "   vibe-redis   # 连接 Redis"
    echo
    echo "   # 或使用本地客户端"
    echo "   mysql -h localhost -P 3306 -u vibe_user -pvibe_password vibe_coding_starter"
    echo "   redis-cli -h localhost -p 6379"
    echo
    echo "4. 查看服务日志："
    echo "   vibe-logs-mysql  # MySQL 日志"
    echo "   vibe-logs-redis  # Redis 日志"
    echo
    echo "5. 管理集群："
    echo "   k3d cluster stop vibe-dev    # 停止集群"
    echo "   k3d cluster start vibe-dev   # 启动集群"
    echo "   k3d cluster delete vibe-dev  # 删除集群"
    echo
    echo "6. 开始开发："
    echo "   cd vibe-coding-starter-api-go"
    echo "   go run cmd/server/main.go -config configs/config-k3d.yaml"
    echo

    echo "有用的别名："
    echo "  k          = kubectl"
    echo "  kgp        = kubectl get pods"
    echo "  kgs        = kubectl get svc"
    echo "  vibe-dev   = 切换到 vibe-dev 命名空间"
    echo "  vibe-status = 查看所有服务状态"
    echo

    echo "文档和资源："
    echo "  • k3d 文档: https://k3d.io/"
    echo "  • Kubernetes 文档: https://kubernetes.io/docs/"
    echo "  • kubectl 备忘单: https://kubernetes.io/docs/reference/kubectl/cheatsheet/"
    echo

    echo -e "${GREEN}🎉 开发环境已准备就绪！祝您开发愉快！${NC}"
    echo
}

# 错误处理函数
handle_error() {
    local exit_code=$?
    log_error "安装过程中发生错误 (退出码: $exit_code)"
    echo
    echo "故障排除建议："
    echo "1. 检查网络连接"
    echo "2. 确保有足够的磁盘空间"
    echo "3. 检查系统权限"
    echo "4. 查看详细错误信息"
    echo
    echo "如需帮助，请查看日志或联系技术支持"
    exit $exit_code
}

# 清理函数
cleanup_on_exit() {
    local exit_code=$?
    if [[ $exit_code -ne 0 ]]; then
        log_warning "检测到异常退出，正在清理..."
        # 这里可以添加清理逻辑
    fi
}

# 主函数
main() {
    # 设置错误处理
    trap handle_error ERR
    trap cleanup_on_exit EXIT

    # 显示横幅
    show_banner

    # 确认安装
    echo -e "${YELLOW}是否继续安装？这将安装 Docker、k3d、kubectl 并创建开发集群。${NC}"
    read -p "请输入 y/Y 继续，或按任意键取消: " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "安装已取消"
        exit 0
    fi

    echo
    log_info "开始安装 Vibe Coding Starter 开发环境..."
    echo

    # 执行安装步骤
    check_system_requirements
    update_system
    install_docker
    install_kubectl
    install_k3d
    verify_docker
    create_storage_directories
    create_k3d_cluster
    deploy_services
    verify_services
    setup_shell_completion
    create_aliases

    # 显示总结和后续步骤
    show_installation_summary
    show_next_steps

    log_success "Vibe Coding Starter 开发环境安装完成！"
}

# 检查是否以 root 用户运行
if [[ $EUID -eq 0 ]]; then
    log_error "请不要以 root 用户运行此脚本"
    log_info "正确用法: ./setup-dev-environment.sh"
    exit 1
fi

# 运行主函数
main "$@"
