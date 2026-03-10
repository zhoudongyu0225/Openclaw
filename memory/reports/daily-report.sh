#!/bin/bash
#===============================================================================
# 每日汇报脚本
# 每天定时向白老师汇报各分身工作进度
#===============================================================================

# 修复 nvm 环境
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
export PATH="$NVM_DIR/versions/node/v22.22.0/bin:$PATH"

REPORT_TIME="12:00"
LOG_FILE="$HOME/.openclaw/workspace/memory/reports/daily-report.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# 读取各分身进度
get_progress() {
    local dir=$1
    local name=$2
    local progress_file="$HOME/.openclaw/workspace/memory/$dir/progress.md"
    local latest_file=$(ls -t "$HOME/.openclaw/workspace/memory/$dir"/*.md 2>/dev/null | head -1)
    
    if [ -f "$progress_file" ]; then
        echo "$(cat "$progress_file")"
    else
        echo "暂无进度"
    fi
}

# 生成汇报
generate_report() {
    local date=$(date '+%Y-%m-%d')
    
    cat << EOF
# 📊 每日分身工作汇报 - $date

---

## 🏥 眠心诊所
$(get_progress "content-ops" "诊所")

---

## 🛒 逆动力电商
$(get_progress "ecommerce" "电商")

---

## 🎮 弹幕游戏
$(get_progress "game-dev" "游戏")

---

## 🎬 游戏广告创意
$(get_progress "ad-creative" "广告")

---

## 🎨 美术素材
$(get_progress "art-gen" "美术")

---

## 🔧 技术支持
$(get_progress "tech-support" "技术")

---

## 📁 产出文件

所有文件位置：~/.openclaw/workspace/memory/

- 眠心诊所：content-ops/
- 逆动力：ecommerce/
- 弹幕游戏：game-dev/
- 游戏广告：ad-creative/
- 美术素材：art-gen/
- 技术支持：tech-support/

---

*汇报时间：$(date)*
EOF
}

# 发送汇报
send_report() {
    local report=$(generate_report)
    
    # 发送到飞书群
    openclaw message send \
        --channel feishu \
        --target oc_a63b77f46415a883feef0b106c5519d1 \
        --message "$report" 2>&1 | grep -v "^\[plugins\]" | grep -v "^\[info\]"
    
    log "汇报已发送"
}

# 主程序
main() {
    log "开始生成每日汇报..."
    
    # 保存汇报到文件
    local report_file="$HOME/.openclaw/workspace/memory/reports/daily-$date.md"
    generate_report > "$report_file"
    
    # 发送汇报
    send_report
    
    log "汇报完成"
}

main "$@"
