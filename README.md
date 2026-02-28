<!-- ©AngelaMos | 2026 -->
<!-- README.md -->

```

           ⠀⠀⠀⠀⠀⠀⢀⣤⣚⣛⡶⣒⡁⠐⠦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⡾⠃⣀⠈⠉⡀⠈⠳⡀⠘⢆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⢸⠃⢸⣿⠀⣼⡇⠀⠀⢣⠀⠈⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⣀⠤⠖⠚⠒⠺⢧⡀⠙⠁⠀⠀⡜⠀⠀⠣⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⢀⠞⢣⠀⡠⠄⠀⠀⠀⠈⠦⠤⠤⠊⠀⠀⠀⠀⠈⢆⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⡞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡄⠀⠀⠀⠀⠀⠀⠘⡆⠀⠀⠀⠀⠀⠀⠀⠀
           ⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⣰⠃⠀⠀⠀⠀⠀⠀⠀⠀
           ⠸⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠞⠦⠀⠀⠀⣀⢼⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠈⠓⠤⢄⣀⣀⣀⡀⠤⠒⠁⠀⣀⠔⡶⠻⢥⠼⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⠈⠁⠒⢒⠏⠈⠉⠀⠀⢷⠤⠮⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⠀⠀⠀⡜⠀⠀⠀⢀⠀⢈⢆⣀⣼⡒⠒⠲⢤⣀⠀⣠⠔⢆⠀
           ⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠈⢧⠀⠙⢜⣆⠑⣄⠀⠀⠈⡝⡄⠀⠈⡄
           ⠀⠀⠀⠀⠀⠀⠀⢰⠋⢃⠀⠀⢀⡀⠬⠂⠀⠈⡏⢦⡀⠑⠢⠔⠁⡸⠀⢀⠁
           ⠀⠀⠀⠀⠀⠀⢀⠼⠤⢸⡄⠀⡏⠀⠀⠀⠀⠰⡁⠀⠉⠒⠒⠒⠈⠀⠀⡸⠀
           ⠀⠀⠀⠀⠀⠀⠈⢇⣀⡀⠘⣦⡞⠀⠀⠀⠀⢠⠟⠀⠀⠀⠀⠀⠠⠀⡰⠁⠀
           ⠀⠀⠀⠀⠀⠀⠀⠀⢀⡩⠭⠤⠿⣲⢖⣒⠚⢫⣀⣀⣀⣀⡀⢀⡧⠊⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⠀⢰⠉⠀⠀⠀⠀⠈⡠⠤⠛⠛⠦⡀⠀⠀⠈⣹⠀⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⠀⠘⣄⠀⠀⠀⠀⡎⠀⠀⠀⠀⠀⠀⠀⠀⠈⠑⡆⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⠀⠀⠸⡉⠒⠒⠒⠳⣄⠀⠀⠀⠀⠀⠀⠀⢀⣠⠃⠀⠀⠀⠀
           ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠉⢧⣈⣉⣀⣈⣉⠭⠝⠒⠊⠀⠀⠀⠀⠀⠀⠀
            ╭────────────────────────────╮
            │  YOSHI-AUDIT  v0.1.0       │
            │  System sanity checker     │
            ╰────────────────────────────╯
```

# yoshi-audit

System monitor, Docker prune manager, and deep audit tool -- all in one TUI.

## Features

**Dashboard**
- Live CPU, memory, disk, and GPU usage with progress bars
- Top-N process table sorted by memory, CPU, GPU, or name
- Zombie and orphan process detection

**Docker**
- Visual disk usage breakdown for images, containers, and volumes
- Multi-select prune with safety categorization and protection rules
- Preset-based cleanup profiles (gentle, moderate, aggressive, nuclear)

**Audit**
- Deep scan for zombies, orphans, runaway daemons, and memory leaks
- Baseline system for learned process whitelisting
- Per-finding actions: inspect, add to baseline, or ignore

## Installation

```sh
go install github.com/CarterPerez-dev/yoshi-audit/cmd@latest
```

## Usage

```sh
yoshi-audit
```

## Keybindings

### Global

| Key     | Action          |
|---------|-----------------|
| `1`     | Dashboard tab   |
| `2`     | Docker tab      |
| `3`     | Audit tab       |
| `Tab`   | Cycle tabs      |
| `P`     | Pause/resume    |
| `Q`     | Quit            |

### Dashboard

| Key | Action          |
|-----|-----------------|
| `M` | Sort by memory  |
| `C` | Sort by CPU     |
| `G` | Sort by GPU     |
| `N` | Sort by name    |

### Docker

| Key       | Action              |
|-----------|---------------------|
| `Space`   | Select/deselect     |
| `A`       | Select all          |
| `P`       | Toggle protect      |
| `D`       | Delete selected     |
| `I`       | Images sub-tab      |
| `C`       | Containers sub-tab  |
| `V`       | Volumes sub-tab     |
| `N`       | Networks sub-tab    |
| `R`       | Refresh data        |
| `1`-`4`   | Apply preset        |

### Audit

| Key     | Action              |
|---------|---------------------|
| `R`     | Rescan              |
| `Enter` | Inspect finding     |
| `A`     | Add to baseline     |
| `I`     | Ignore finding      |

## Configuration

Config file location:

```
~/.config/yoshi-audit/config.yaml
```

## Requirements

- Go 1.24+
- Docker (for Docker tab)
- nvidia-smi (optional, for GPU stats)
