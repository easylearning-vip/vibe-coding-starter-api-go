#!/bin/bash

# k3d 和相关工具自动安装脚本
# 支持 Linux 和 macOS

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# 检测操作系统和架构
detect_system() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="darwin"
    else
        log_error "不支持的操作系统: $OSTYPE"
        exit 1
    fi
    
    ARCH=$(uname -m)
    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH="amd64"
    elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
        ARCH="arm64"
    else
        log_error "不支持的架构: $ARCH"
        exit 1
    fi
    
    log_info "检测到系统: $OS-$ARCH"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        log_info "运行安装脚本: ./install-docker.sh"
        exit 1
    fi

    # 检查 Docker 是否运行，首先尝试不使用 sudo
    if docker info &> /dev/null; then
        log_success "Docker 已安装并运行"
        return
    fi

    # 如果不使用 sudo 失败，尝试使用 sudo
    if sudo docker info &> /dev/null; then
        log_success "Docker 已安装并运行 (需要 sudo 权限)"
        log_warning "当前用户需要 sudo 权限访问 Docker"
        log_info "建议将当前用户添加到 docker 组: sudo usermod -aG docker \$USER"
        log_info "然后重新登录或运行: newgrp docker"
        return
    fi

    # 如果两种方式都失败，说明 Docker 服务未运行
    log_error "Docker 服务未运行，请启动 Docker"
    log_info "尝试启动 Docker 服务:"
    log_info "  - systemd: sudo systemctl start docker"
    log_info "  - service: sudo service docker start"
    log_info "  - Docker Desktop: 启动 Docker Desktop 应用"
    exit 1
}

# 安装 kubectl
install_kubectl() {
    if command -v kubectl &> /dev/null; then
        log_info "kubectl 已安装: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
        return
    fi
    
    log_info "安装 kubectl..."
    
    if [[ "$OS" == "darwin" ]] && command -v brew &> /dev/null; then
        brew install kubectl
    else
        # 获取最新版本
        KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)
        
        # 下载 kubectl
        curl -LO "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/$OS/$ARCH/kubectl"
        
        # 设置执行权限
        chmod +x kubectl
        
        # 移动到系统路径
        sudo mv kubectl /usr/local/bin/
    fi
    
    log_success "kubectl 安装完成"
}

# 安装 k3d
install_k3d() {
    if command -v k3d &> /dev/null; then
        log_info "k3d 已安装: $(k3d version)"
        return
    fi
    
    log_info "安装 k3d..."
    
    if [[ "$OS" == "darwin" ]] && command -v brew &> /dev/null; then
        brew install k3d
    else
        # 使用官方安装脚本
        curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
    fi
    
    log_success "k3d 安装完成"
}

# 安装 helm (可选)
install_helm() {
    if command -v helm &> /dev/null; then
        log_info "Helm 已安装: $(helm version --short)"
        return
    fi
    
    log_info "安装 Helm..."
    
    if [[ "$OS" == "darwin" ]] && command -v brew &> /dev/null; then
        brew install helm
    else
        curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    fi
    
    log_success "Helm 安装完成"
}

# 验证安装
verify_installation() {
    log_info "验证安装..."
    
    local missing_tools=()
    
    # 检查 kubectl
    if ! command -v kubectl &> /dev/null; then
        missing_tools+=("kubectl")
    fi
    
    # 检查 k3d
    if ! command -v k3d &> /dev/null; then
        missing_tools+=("k3d")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "安装失败，缺少工具: ${missing_tools[*]}"
        exit 1
    fi
    
    log_success "所有工具安装成功！"
    echo
    echo "已安装的工具版本："
    # 获取 Docker 版本，优先尝试不使用 sudo
    if docker --version &> /dev/null; then
        echo "- Docker: $(docker --version)"
    elif sudo docker --version &> /dev/null; then
        echo "- Docker: $(sudo docker --version) (需要 sudo)"
    else
        echo "- Docker: 版本获取失败"
    fi
    echo "- kubectl: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
    echo "- k3d: $(k3d version)"
    if command -v helm &> /dev/null; then
        echo "- Helm: $(helm version --short)"
    fi
}

# 配置 shell 补全 (可选)
setup_completion() {
    log_info "设置 shell 补全..."
    
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
    
    # 创建补全目录
    mkdir -p "$COMPLETION_DIR"
    
    # 生成补全脚本
    if command -v kubectl &> /dev/null; then
        kubectl completion $CURRENT_SHELL > "$COMPLETION_DIR/kubectl"
    fi
    
    if command -v k3d &> /dev/null; then
        k3d completion $CURRENT_SHELL > "$COMPLETION_DIR/k3d"
    fi
    
    if command -v helm &> /dev/null; then
        helm completion $CURRENT_SHELL > "$COMPLETION_DIR/helm"
    fi
    
    # 添加到 profile 文件
    if [[ "$CURRENT_SHELL" == "bash" ]]; then
        echo "# Load kubectl, k3d, helm completions" >> "$PROFILE_FILE"
        echo "for f in ~/.bash_completion.d/*; do source \$f; done" >> "$PROFILE_FILE"
    elif [[ "$CURRENT_SHELL" == "zsh" ]]; then
        echo "# Load kubectl, k3d, helm completions" >> "$PROFILE_FILE"
        echo "fpath=(~/.zsh/completions \$fpath)" >> "$PROFILE_FILE"
        echo "autoload -U compinit && compinit" >> "$PROFILE_FILE"
    fi
    
    log_success "Shell 补全设置完成，请重新加载 shell 或重新登录"
}

# 显示后续步骤
show_next_steps() {
    echo
    log_success "k3d 工具链安装完成！"
    echo
    echo "后续步骤："
    echo "1. 重新加载 shell 以启用补全功能:"
    echo "   source ~/.bashrc  # 或 source ~/.zshrc"
    echo
    echo "2. 测试 k3d 安装:"
    echo "   k3d version"
    echo
    echo "3. 创建测试集群:"
    echo "   k3d cluster create test --agents 1"
    echo "   kubectl get nodes"
    echo "   k3d cluster delete test"
    echo
    echo "4. 开始使用开发环境:"
    echo "   cd vibe-coding-starter-api-go/tutorials/develop/k3d"
    echo "   k3d cluster create --config k3d-cluster.yaml"
    echo "   kubectl apply -f manifests/"
    echo
    echo "5. 学习资源:"
    echo "   - k3d 文档: https://k3d.io/"
    echo "   - Kubernetes 文档: https://kubernetes.io/docs/"
    echo "   - kubectl 备忘单: https://kubernetes.io/docs/reference/kubectl/cheatsheet/"
}

# 主函数
main() {
    log_info "开始安装 k3d 工具链..."
    
    detect_system
    check_docker
    install_kubectl
    install_k3d
    install_helm
    verify_installation
    setup_completion
    show_next_steps
    
    log_success "k3d 工具链安装完成！"
}

# 运行主函数
main "$@"
