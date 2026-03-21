# TOOLS.md - Local Notes

Skills define _how_ tools work. This file is for _your_ specifics — the stuff that's unique to your setup.

## What Goes Here

Things like:

- Camera names and locations
- SSH hosts and aliases
- Preferred voices for TTS
- Speaker/room names
- Device nicknames
- Anything environment-specific

## Examples

```markdown
### Cameras

- living-room → Main area, 180° wide angle
- front-door → Entrance, motion-triggered

### SSH

- home-server → 192.168.1.100, user: admin

### TTS

- Preferred voice: "Nova" (warm, slightly British)
- Default speaker: Kitchen HomePod
```

## Why Separate?

Skills are shared. Your setup is yours. Keeping them apart means you can update skills without losing your notes, and share skills without leaking your infrastructure.

---

Add whatever helps you do your job. This is your cheat sheet.

## Git 操作规范 ⚠️ 重要

### Clone - 必须全量
```bash
git clone https://github.com/yangkai258/mdm-iot-platform.git --全量克隆，不要用 --depth=1
```

### Fetch - 必须全量
```bash
git fetch origin --unshallow --如果已经是浅克隆，补全完整历史
```

### 问题排查
如果发现 GitHub 上有但本地没有的文件：
```bash
git fetch origin --unshallow
git checkout FETCH_HEAD -- docs/  # 提取缺失的文件
```

### 记住
- **Clone 和 Fetch 永远用全量**（不用 --depth=1）
- 如果误用了浅克隆，用 `git fetch --unshallow` 补救
