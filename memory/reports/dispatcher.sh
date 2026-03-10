#!/bin/bash
#===============================================================================
# 分身任务分发器
# 每小时自动执行各分身的任务
#===============================================================================

WORKSPACE="$HOME/.openclaw/workspace"
MEMORY_DIR="$WORKSPACE/memory"
QUEUE_DIR="$MEMORY_DIR/task-queue"
LOG_FILE="$MEMORY_DIR/reports/dispatcher.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# 读取任务内容并转换为 subagent prompt
get_task_prompt() {
    local queue_file="$1"
    local agent_name="$2"
    
    # 读取任务列表
    local tasks=$(grep -E "^\- \[ \]" "$queue_file" 2>/dev/null | sed 's/^- \[ \] //')
    
    if [ -z "$tasks" ]; then
        echo ""
        return
    fi
    
    # 根据分身类型生成不同的 prompt
    case "$agent_name" in
        ad-creative)
            echo "你是广告创意专家（SLG买量）。今日任务：
$tasks

产出要求：
- 必须有真实的素材链接
- 不要编造数据
- 完成后记录到 $MEMORY_DIR/ad-creative/daily-work-$(date +%Y-%m-%d).md"
            ;;
        content-ops)
            echo "你是内容运营专家（眠心诊所）。今日任务：
$tasks

产出要求：
- 内容要有价值
- 基于真实洞察
- 完成后记录到 $MEMORY_DIR/content-ops/daily-work-$(date +%Y-%m-%d).md"
            ;;
        ecommerce)
            echo "你是电商运营专家（氢美健康）。今日任务：
$tasks

产出要求：
- 基于数据分析
- 完成后记录到 $MEMORY_DIR/ecommerce/daily-work-$(date +%Y-%m-%d).md"
            ;;
        game-dev)
            echo "你是游戏开发专家（弹幕游戏）。今日任务：
$tasks

产出要求：
- 写真正能用的代码
- 完成后记录到 $MEMORY_DIR/game-dev/daily-work-$(date +%Y-%m-%d).md"
            ;;
        art-gen)
            echo "你是美术素材专家（AI生成）。今日任务：
$tasks

产出要求：
- 提示词要具体可执行
- 完成后记录到 $MEMORY_DIR/art-gen/daily-work-$(date +%Y-%m-%d).md"
            ;;
        tech-support)
            echo "你是技术支持专家。今日任务：
$tasks

产出要求：
- 给可执行的代码
- 完成后记录到 $MEMORY_DIR/tech-support/daily-work-$(date +%Y-%m-%d).md"
            ;;
    esac
}

# 执行分身任务
run_agent() {
    local agent_name="$1"
    local queue_file="$2"
    
    # 获取任务 prompt
    local prompt=$(get_task_prompt "$queue_file" "$agent_name")
    
    if [ -z "$prompt" ]; then
        log "  [$agent_name] 无待执行任务"
        return 0
    fi
    
    log "  [$agent_name] 开始执行任务..."
    
    # 调用 OpenClaw 执行 subagent
    # 使用 openclaw CLI 执行任务
    local result=$(openclaw sessions spawn \
        --label "$agent_name" \
        --mode run \
        --runtime subagent \
        --timeout 600 \
        "$prompt" 2>&1)
    
    if echo "$result" | grep -q "accepted"; then
        log "  [$agent_name] 任务已提交"
        # 标记任务完成（删除已完成的任务）
        sed -i 's/^- \[ \] /- [x] /' "$queue_file"
        return 0
    else
        log "  [$agent_name] 任务提交失败: $result"
        return 1
    fi
}

# 主流程
main() {
    log "========== 任务分发器启动 =========="
    
    # 检查各分身任务队列
    for agent in ad-creative content-ops ecommerce game-dev art-gen tech-support; do
        queue_file="$QUEUE_DIR/$agent.md"
        
        if [ -f "$queue_file" ]; then
            # 检查是否有待执行任务
            if grep -q "^\- \[ \] " "$queue_file"; then
                run_agent "$agent" "$queue_file"
            else
                log "  [$agent] 队列已清空，补充新任务..."
                # 重新填充任务队列
                case "$agent" in
                    ad-creative)
                        cat > "$queue_file" << 'EOF'
# 广告创意 - 今日任务
## 待执行
- [ ] 分析今日爆款素材
- [ ] 产出新脚本创意
- [ ] 更新 case-studies.md
EOF
                        ;;
                    content-ops)
                        cat > "$queue_file" << 'EOF'
# 内容运营 - 今日任务
## 待执行
- [ ] 撰写公众号文章
- [ ] 产出抖音脚本
- [ ] 分析竞品
EOF
                        ;;
                    ecommerce)
                        cat > "$queue_file" << 'EOF'
# 电商运营 - 今日任务
## 待执行
- [ ] 分析转化漏斗
- [ ] 产出推广文案
- [ ] 竞品监控
EOF
                        ;;
                    game-dev)
                        cat > "$queue_file" << 'EOF'
# 游戏开发 - 今日任务
## 待执行
- [ ] 完善核心模块
- [ ] 产出代码
- [ ] 更新文档
EOF
                        ;;
                    art-gen)
                        cat > "$queue_file" << 'EOF'
# 美术素材 - 今日任务
## 待执行
- [ ] 优化 AI 提示词
- [ ] 产出新模板
- [ ] 更新提示词库
EOF
                        ;;
                    tech-support)
                        cat > "$queue_file" << 'EOF'
# 技术支持 - 今日任务
## 待执行
- [ ] 检查服务状态
- [ ] 技术优化
- [ ] 更新文档
EOF
                        ;;
                esac
            fi
        fi
    done
    
    log "========== 任务分发完成 =========="
}

main "$@"
