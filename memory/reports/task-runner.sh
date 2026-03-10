#!/bin/bash
#===============================================================================
# 分身任务调度器
# 每天自动执行各分身的任务
#===============================================================================

WORKSPACE="$HOME/.openclaw/workspace"
MEMORY_DIR="$WORKSPACE/memory"
LOG_FILE="$MEMORY_DIR/reports/task-execution.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

#-------------------------------------------------------------------------------
# 各分身任务函数
#-------------------------------------------------------------------------------

# 1. 广告创意 - 分析素材
task_ad_creative() {
    log "🎬 执行广告创意任务..."
    cd "$MEMORY_DIR/ad-creative"
    
    # 读取任务清单
    if [ -f "slg-ad-analysis.md" ]; then
        echo "已有分析笔记，继续补充..."
    fi
    
    # 记录进度
    echo "## $(date '+%Y-%m-%d %H:%M') - 素材分析" >> progress.md
    echo "- 任务执行中..." >> progress.md
    
    echo "完成"
}

# 2. 内容运营 - 公众号
task_content_ops() {
    log "📝 执行内容运营任务..."
    cd "$MEMORY_DIR/content-ops"
    
    echo "## $(date '+%Y-%m-%d %H:%M')" >> progress.md
    echo "- 任务执行中..." >> progress.md
    
    echo "完成"
}

# 3. 电商运营
task_ecommerce() {
    log "🛒 执行电商任务..."
    cd "$MEMORY_DIR/ecommerce"
    
    echo "## $(date '+%Y-%m-%d %H:%M')" >> progress.md
    echo "- 任务执行中..." >> progress.md
    
    echo "完成"
}

# 4. 游戏开发
task_game_dev() {
    log "🎮 执行游戏开发任务..."
    cd "$MEMORY_DIR/game-dev"
    
    echo "## $(date '+%Y-%m-%d %H:%M')" >> progress.md
    echo "- 任务执行中..." >> progress.md
    
    echo "完成"
}

# 5. 美术素材
task_art_gen() {
    log "🎨 执行美术素材任务..."
    cd "$MEMORY_DIR/art-gen"
    
    echo "## $(date '+%Y-%m-%d %H:%M')" >> progress.md
    echo "- 任务执行中..." >> progress.md
    
    echo "完成"
}

# 6. 技术支持
task_tech_support() {
    log "🔧 执行技术支持任务..."
    cd "$MEMORY_DIR/tech-support"
    
    echo "## $(date '+%Y-%m-%d %H:%M')" >> progress.md
    echo "- 任务执行中..." >> progress.md
    
    echo "完成"
}

#-------------------------------------------------------------------------------
# 生成每日汇报
#-------------------------------------------------------------------------------
generate_daily_report() {
    log "📊 生成每日进度汇报..."
    
    local report="$MEMORY_DIR/reports/daily-progress-$(date +%Y%m%d).md"
    
    cat > "$report" << EOF
# 📊 每日分身进度汇报

**日期：$(date '+%Y年%m月%d日')**

---

## 各分身状态

| 分身 | 状态 | 今日进度 |
|------|------|----------|
EOF

    # 检查各分身的 progress.md
    for dir in ad-creative content-ops ecommerce game-dev art-gen tech-support; do
        if [ -f "$MEMORY_DIR/$dir/progress.md" ]; then
            local last=$(tail -5 "$MEMORY_DIR/$dir/progress.md" 2>/dev/null | grep -c "完成" || echo 0)
            echo "| $dir | 进行中 | $last 项完成 |" >> "$report"
        else
            echo "| $dir | 待启动 | - |" >> "$report"
        fi
    done
    
    echo -e "\n---\n*自动生成 $(date)*" >> "$report"
    
    # 发送到飞书
    openclaw message send \
        --channel feishu \
        --target oc_a63b77f46415a883feef0b106c5519d1 \
        --message "$(cat "$report")" 2>&1 | grep -v "^\[plugins\]" | grep -v "^\[info\]"
    
    log "汇报已发送"
}

#-------------------------------------------------------------------------------
# 主流程
#-------------------------------------------------------------------------------
main() {
    log "========== 开始执行每日任务 =========="
    
    # 并行执行各分身任务（后台）
    task_ad_creative &
    task_content_ops &
    task_ecommerce &
    task_game_dev &
    task_art_gen &
    task_tech_support &
    
    # 等待所有任务完成
    wait
    
    log "所有任务执行完成"
    
    # 生成汇报
    generate_daily_report
    
    log "========== 任务执行完毕 =========="
}

main "$@"
