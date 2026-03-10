#!/usr/bin/env python3
"""
游戏UI素材处理工具
功能：拼图、裁剪、缩放、格式转换、添加文字
"""
import os
import sys
from PIL import Image, ImageDraw, ImageFont

def create_grid(images, rows, cols, cell_size=300, padding=10):
    """创建图片网格拼图"""
    w, h = cell_size * cols + padding * (cols + 1), cell_size * rows + padding * (rows + 1)
    result = Image.new('RGB', (w, h), (20, 20, 30))
    for i, img_path in enumerate(images[:rows*cols]):
        if os.path.exists(img_path):
            img = Image.open(img_path).convert('RGB')
            img = img.resize((cell_size, cell_size), Image.Resampling.LANCZOS)
            x = (i % cols) * (cell_size + padding) + padding
            y = (i // cols) * (cell_size + padding) + padding
            result.paste(img, (x, y))
    return result

def add_watermark(img, text="游戏素材", opacity=128):
    """添加水印"""
    overlay = Image.new('RGBA', img.size, (0, 0, 0, 0))
    draw = ImageDraw.Draw(overlay)
    try:
        font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 36)
    except:
        font = ImageFont.load_default()
    
    # 简单文字水印
    draw.text((10, 10), text, fill=(255, 255, 255, opacity))
    
    if img.mode != 'RGBA':
        img = img.convert('RGBA')
    return Image.alpha_composite(img, overlay).convert('RGB')

def resize_keep_ratio(img, target_size):
    """保持比例缩放"""
    img.thumbnail(target_size, Image.Resampling.LANCZOS)
    return img

def create_thumbnail_grid(image_dir, output_path, cols=4, thumb_size=200):
    """为目录中所有图片创建缩略图网格"""
    images = []
    for f in os.listdir(image_dir):
        if f.lower().endswith(('.png', '.jpg', '.jpeg', '.gif')):
            images.append(os.path.join(image_dir, f))
    
    if not images:
        print("No images found")
        return
    
    rows = (len(images) + cols - 1) // cols
    grid = create_grid(images, rows, cols, thumb_size)
    grid.save(output_path)
    print(f"Created: {output_path}")

# CLI入口
if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python image_tool.py <command> [args]")
        print("Commands:")
        print("  grid <dir> <output> - Create thumbnail grid")
        print("  watermark <input> <output> <text> - Add watermark")
        sys.exit(1)
    
    cmd = sys.argv[1]
    if cmd == "grid" and len(sys.argv) >= 4:
        create_thumbnail_grid(sys.argv[2], sys.argv[3])
    else:
        print("Unknown command")
