#!/bin/bash

# Game Server Deployment Script
# Usage: ./deploy.sh [start|stop|restart|status|logs|build]

set -e

IMAGE_NAME="bullet-hell-game"
CONTAINER_NAME="bullet-hell-server"
PORT_HTTP=8080
PORT_WS=8081
PORT_PPROF=6060

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
}

build() {
    log_info "Building Docker image..."
    docker build -t $IMAGE_NAME:latest -f deploy/docker .
    log_info "Build completed"
}

start() {
    check_docker
    log_info "Starting game server..."
    
    # Check if .env exists
    if [ ! -f .env ]; then
        log_warn ".env file not found, using defaults"
        cp .env.example .env 2>/dev/null || true
    fi
    
    # Start services
    docker compose up -d
    log_info "Game server started"
    status
}

stop() {
    log_info "Stopping game server..."
    docker compose down
    log_info "Game server stopped"
}

restart() {
    stop
    start
}

status() {
    log_info "Checking service status..."
    
    # Check container status
    if docker ps --format '{{.Names}}' | grep -q $CONTAINER_NAME; then
        log_info "Container is running"
        docker ps --filter "name=$CONTAINER_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    else
        log_warn "Container is not running"
    fi
    
    # Check ports
    for port in $PORT_HTTP $PORT_WS $PORT_PPROF; do
        if ss -tuln | grep -q ":$port "; then
            log_info "Port $port is LISTENING"
        else
            log_warn "Port $port is NOT listening"
        fi
    done
    
    # Health check
    if curl -sf http://localhost:$PORT_HTTP/health > /dev/null 2>&1; then
        log_info "Health check: OK"
    else
        log_warn "Health check: FAILED"
    fi
}

logs() {
    docker compose logs -f --tail=100
}

backup() {
    BACKUP_DIR="./backups"
    mkdir -p $BACKUP_DIR
    
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.tar.gz"
    
    log_info "Creating backup..."
    
    # Backup MongoDB
    docker exec bullet-hell-mongodb mongodump --archive=/tmp/backup.gz --gzip
    docker cp bullet-hell-mongodb:/tmp/backup.gz $BACKUP_FILE
    
    # Backup Redis
    docker exec bullet-hell-redis redis-cli SAVE
    docker cp bullet-hell-redis:/data/dump.rdb $BACKUP_DIR/redis_$TIMESTAMP.rdb
    
    log_info "Backup created: $BACKUP_FILE"
}

restore() {
    if [ -z "$1" ]; then
        log_error "Usage: $0 restore <backup_file>"
        exit 1
    fi
    
    BACKUP_FILE=$1
    
    if [ ! -f "$BACKUP_FILE" ]; then
        log_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    
    log_info "Restoring from backup..."
    
    docker cp $BACKUP_FILE bullet-hell-mongodb:/tmp/restore.gz
    docker exec bullet-hell-mongodb mongorestore --archive=/tmp/restore.gz --gzip --drop
    
    log_info "Restore completed"
}

scale() {
    if [ -z "$1" ]; then
        log_error "Usage: $0 scale <number_of_instances>"
        exit 1
    fi
    
    log_info "Scaling to $1 instances..."
    docker compose up -d --scale game-server=$1
    log_info "Scaled to $1 instances"
}

cleanup() {
    log_info "Cleaning up..."
    docker system prune -f
    log_info "Cleanup completed"
}

monitor() {
    log_info "Monitoring resources..."
    
    echo "=== Docker Stats ==="
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"
    
    echo ""
    echo "=== Disk Usage ==="
    df -h /var/lib/docker
    
    echo ""
    echo "=== Memory Info ==="
    free -h
}

case "$1" in
    build)
        build
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs
        ;;
    backup)
        backup
        ;;
    restore)
        restore $2
        ;;
    scale)
        scale $2
        ;;
    cleanup)
        cleanup
        ;;
    monitor)
        monitor
        ;;
    *)
        echo "Usage: $0 {build|start|stop|restart|status|logs|backup|restore|scale|cleanup|monitor}"
        exit 1
        ;;
esac
