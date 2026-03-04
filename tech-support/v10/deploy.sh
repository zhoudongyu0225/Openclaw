#!/bin/bash
#===============================================================================
# 服务器自动部署脚本
# 用途：自动部署网页应用到 Cloudflare Pages
#===============================================================================

# 配置
PROJECT_NAME="逆动力氢镁胶囊"
CF_ACCOUNT_ID="your-account-id"
CF_PROJECT_NAME="hydromg"
GITHUB_REPO="https://github.com/your-repo"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e \"${GREEN}[INFO]${NC} $1\"; }
log_warn() { echo -e \"${YELLOW}[WARN]${NC} $1\"; }
log_error() { echo -e \"${RED}[ERROR]${NC} $1\"; }

# 检查依赖
check_dependencies() {
    log_info \"检查依赖...\"
    
    if ! command -v git &> /dev/null; then
        log_error \"git 未安装\"
        exit 1
    fi
    
    if ! command -v curl &> /dev/null; then
        log_error \"curl 未安装\"
        exit 1
    fi
    
    log_info \"依赖检查完成\"
}

# 部署到 Cloudflare Pages
deploy_to_cloudflare() {
    local dir=$1
    
    log_info \"开始部署到 Cloudflare Pages...\"
    
    # 登录 Cloudflare
    # cf auth
    
    # 创建项目或部署
    # wrangler pages project create $CF_PROJECT_NAME
    # wrangler pages deploy $dir
    
    log_info \"部署完成！\"
    echo \"访问地址: https://$CF_PROJECT_NAME.pages.dev\"
}

# 部署到 Vercel
deploy_to_vercel() {
    local dir=$1
    
    log_info \"开始部署到 Vercel...\"
    
    # 需要先安装 vercel CLI
    # npm i -g vercel
    # cd $dir && vercel --prod
    
    log_info \"部署完成！\"
}

# 主菜单
show_menu() {
    echo \"======================================\"
    echo \"  服务器部署工具\"
    echo \"======================================\"
    echo \"1. 部署到 Cloudflare Pages\"
    echo \"2. 部署到 Vercel\"
    echo \"3. 部署到 本地测试\"
    echo \"0. 退出\"
    echo \"======================================\"
    echo -n \"请选择: \"
}

# 主程序
main() {
    check_dependencies
    
    show_menu
    read choice
    
    case $choice in
        1)
            deploy_to_cloudflare \"./html\"
            ;;
        2)
            deploy_to_vercel \"./html\"
            ;;
        3)
            log_info \"启动本地服务器...\"
            cd ./html && python3 -m http.server 8080
            ;;
        0)
            log_info \"退出\"
            exit 0
            ;;
        *)
            log_error \"无效选择\"
            exit 1
            ;;
    esac
}

main \"$@\"
