#!/bin/bash

# Docker 和 Docker Compose 自动安装脚本
# 支持 Ubuntu/Debian、CentOS/RHEL、macOS

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

# 检测操作系统
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            OS=$ID
            VERSION=$VERSION_ID
        else
            log_error "无法检测 Linux 发行版"
            exit 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        log_error "不支持的操作系统: $OSTYPE"
        exit 1
    fi
    
    log_info "检测到操作系统: $OS"
}

# 检查是否已安装 Docker
check_docker_installed() {
    if command -v docker &> /dev/null; then
        log_warning "Docker 已安装: $(docker --version)"
        return 0
    else
        return 1
    fi
}

# Ubuntu/Debian 安装 Docker
install_docker_ubuntu() {
    log_info "在 Ubuntu/Debian 上安装 Docker..."
    
    # 更新包索引
    sudo apt update
    
    # 安装必要的包
    sudo apt install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    
    # 添加 Docker 官方 GPG 密钥
    curl -fsSL https://download.docker.com/linux/$OS/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    
    # 添加 Docker 仓库
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/$OS $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # 安装 Docker Engine
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    log_success "Docker 安装完成"
}

# CentOS/RHEL 安装 Docker
install_docker_centos() {
    log_info "在 CentOS/RHEL 上安装 Docker..."
    
    # 安装必要的包
    sudo yum install -y yum-utils
    
    # 添加 Docker 仓库
    sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    
    # 安装 Docker Engine
    sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    log_success "Docker 安装完成"
}

# macOS 安装 Docker
install_docker_macos() {
    log_info "在 macOS 上安装 Docker..."
    
    if command -v brew &> /dev/null; then
        log_info "使用 Homebrew 安装 Docker Desktop..."
        brew install --cask docker
    else
        log_warning "未检测到 Homebrew"
        log_info "请手动下载并安装 Docker Desktop for Mac:"
        log_info "https://docs.docker.com/desktop/mac/install/"
        exit 1
    fi
    
    log_success "Docker Desktop 安装完成"
    log_warning "请启动 Docker Desktop 应用程序"
}

# 配置 Docker
configure_docker() {
    log_info "配置 Docker..."
    
    if [[ "$OS" != "macos" ]]; then
        # 启动 Docker 服务
        sudo systemctl start docker
        sudo systemctl enable docker
        
        # 将当前用户添加到 docker 组
        sudo usermod -aG docker $USER
        
        log_success "Docker 服务已启动并设置为开机自启"
        log_warning "请重新登录或运行 'newgrp docker' 使组更改生效"
    fi
}

# 验证安装
verify_installation() {
    log_info "验证 Docker 安装..."
    
    # 等待 Docker 服务启动
    sleep 5
    
    if [[ "$OS" == "macos" ]]; then
        log_warning "请确保 Docker Desktop 已启动"
        log_info "您可以运行 'docker --version' 来验证安装"
    else
        # 测试 Docker 是否正常工作
        if sudo docker run --rm hello-world &> /dev/null; then
            log_success "Docker 安装验证成功！"
        else
            log_error "Docker 安装验证失败"
            exit 1
        fi
    fi
    
    # 显示版本信息
    echo
    log_success "安装完成！版本信息："
    docker --version
    docker compose version
}

# 显示后续步骤
show_next_steps() {
    echo
    log_success "Docker 安装完成！"
    echo
    echo "后续步骤："
    echo "1. 如果是 Linux 系统，请重新登录或运行: newgrp docker"
    echo "2. 如果是 macOS 系统，请启动 Docker Desktop 应用程序"
    echo "3. 验证安装: docker --version"
    echo "4. 运行测试: docker run --rm hello-world"
    echo "5. 开始使用开发环境:"
    echo "   cd vibe-coding-starter-api-go/tutorials/develop/docker-compose"
    echo "   docker compose -f docker-compose.dev.yml up -d"
}

# 主函数
main() {
    log_info "开始安装 Docker..."
    
    detect_os
    
    if check_docker_installed; then
        log_info "Docker 已安装，跳过安装步骤"
    else
        case $OS in
            "ubuntu"|"debian")
                install_docker_ubuntu
                ;;
            "centos"|"rhel"|"fedora")
                install_docker_centos
                ;;
            "macos")
                install_docker_macos
                ;;
            *)
                log_error "不支持的操作系统: $OS"
                exit 1
                ;;
        esac
        
        configure_docker
    fi
    
    verify_installation
    show_next_steps
}

# 运行主函数
main "$@"
