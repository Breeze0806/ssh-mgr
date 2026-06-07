# ssh连接管理器

`ssh-mgr` 是一个管理 [PuTTY](https://putty.org/) / [WinSCP](https://winscp.net/) / [FileZilla](https://filezilla-project.org/) 连接信息的命令行工具。连接信息用 **SM3 + SM4**（国密）加密存储在本地，主密码解锁 SM4 key。

适用于 **Windows**（主用场景）和 **Linux**（GNOME / KDE 等桌面环境）。同一个二进制不交叉：Windows 上是 `mgr.exe`，Linux 上是 `ssh-mgr`，两边各管各的 `config.json`（AGENTS.md 提到的 `cmd/putty` 二进制仓库里没提交，不影响主流程）。

## 优点

1. 用 [PuTTY](https://putty.org/) 一键 SSH 登录任一支持 SSH 的机器
2. 用 [WinSCP](https://winscp.net/) 或 [FileZilla](https://filezilla-project.org/) 上传/下载文件
3. 交互式 REPL（`go-prompt` 驱动），Tab 补全、命令提示
4. 集中管理 PuTTY / WinSCP / FileZilla 的连接信息
5. **免密登录**：连接信息里存的就是明文密码（加密保存），启动时自动喂给客户端
6. 用 **国密 SM3 + SM4** 加密保存用户名/密码，密钥派生自主密码

## 系统要求

### 通用

- Go 1.20+（编译用，运行时不需要）
- 图形桌面（PuTTY / FileZilla 都是 GUI 应用）

### Windows

- Windows 7 / Server 2012 或更新
- 已安装 PuTTY（`putty.exe`）和/或 WinSCP（`WinSCP.exe`）/ FileZilla

### Linux

- 任意主流发行版（Ubuntu / Debian / Fedora / Arch / Manjaro 等）
- X11 或 Wayland
- 已安装 [原生 Linux 版 PuTTY](https://putty.org/)（apt 包名 `putty`，提供 `/usr/bin/putty`，GTK3 GUI）和/或 FileZilla（apt 包名 `filezilla`）
- 推荐在 GNOME / KDE / XFCE 等主流桌面上使用

## 编译

### Windows 上编译

```cmd
cd cmd/mgr
go build
REM 产物：.\mgr.exe
```

### Linux / macOS 上编译

```bash
cd cmd/mgr
go build -o ssh-mgr .
# 产物：./ssh-mgr
```

### 交叉编译

Linux 上给 Windows 出包（`-s -w` 砍符号表，从 ~4.3 MB 缩到 ~2.9 MB）：

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o mgr.exe ./cmd/mgr
```

## 安装

### Windows

直接下 [v1.0.0 release](https://github.com/Breeze0806/ssh-mgr/releases/tag/v1.0.0) 的 64 位 zip，解压到任意目录。

### Linux

Debian / Ubuntu：

```bash
sudo apt update
sudo apt install putty filezilla       # 运行时依赖

git clone https://github.com/Breeze0806/ssh-mgr.git
cd ssh-mgr/cmd/mgr
go build -o ~/.local/bin/ssh-mgr .
```

Fedora：

```bash
sudo dnf install putty filezilla
```

Arch / Manjaro：

```bash
sudo pacman -S putty filezilla
```

## 配置

`-c` 指定配置文件，默认 `config.json`（在二进制同目录）：

```
Usage of mgr:
  -c string
        config file (default "config.json")
```

### Windows 配置示例

```json
{
    "ssh": "C:\\Program Files\\PuTTY\\putty.exe",
    "sftp": "C:\\Program Files (x86)\\WinSCP\\WinSCP.exe",
    "source": "C:\\Linux\\ssh",
    "isEncrypted": true,
    "password": "C:\\Linux\\passwd"
}
```

### Linux 配置示例

```json
{
    "ssh": "putty",
    "sftp": "filezilla",
    "source": "/home/yourname/.ssh-mgr/connections",
    "isEncrypted": true,
    "password": "/home/yourname/.ssh-mgr/passwd"
}
```

### 字段说明

| 字段 | 必填 | 说明 |
|------|------|------|
| `ssh` | 是 | SSH 客户端可执行文件。Windows 写绝对路径；Linux 写 `$PATH` 里的名字（`putty`）即可 |
| `sftp` | 是 | SFTP 客户端可执行文件。Windows 常用 WinSCP，Linux 常用 FileZilla |
| `source` | 是 | 连接信息存储目录。工具按 `<group>/<name>.json` 组织子目录和文件 |
| `isEncrypted` | 是 | 是否用 SM3+SM4 加密 `source` 里的 user/password 字段。强烈建议 `true` |
| `password` | 是 | 主密码的 SM3 哈希文件路径。第一次启动会让你输两次设置主密码，工具会写哈希到这里 |

## 首次使用

`isEncrypted: true` 且 `password` 文件不存在 → 让你设置主密码（输两次）：

```
please input password: ********
please confirm password: ********
```

主密码会用来解锁所有加密的连接信息。**忘记主密码 = 数据全丢**，工具不存明文，没救。

之后每次启动会先让你输入主密码解锁：

```
please input password: ********
>
```

## REPL 命令

| 命令 | 用法 | 说明 |
|------|------|------|
| `ssh` | `ssh <group> <name>` | 启动 PuTTY 连过去 |
| `sftp` | `sftp <group> <name>` | 启动 WinSCP / FileZilla |
| `add` | `add <group> <name>` | 新增一条连接（交互式问 address/user/password） |
| `show` | `show <group> <name>` | 打印这条连接的明文凭据（`user:password@host`） |
| `showAddr` | `showAddr <addr>` | 按 `ip:port` 查找匹配的所有连接 |
| `exit` | `exit` | 退出 REPL |

### `add` 命令示例

```bash
> add prod web1
please input ssh address:1.2.3.4:2222
please input ssh user:root
please input ssh password:********
add success！
```

`address` 写法：

- `1.2.3.4` → 默认端口 22（等价于 `1.2.3.4:22`）
- `1.2.3.4:2222` → 显式指定端口
- `[2001:db8::1]:22` → IPv6（要带方括号）

### Tab 补全

命令首字母 + Tab 出候选：

```bash
> ss
      ssh  ssh group name
> ssh t
       test   test
       test1  test1
```

### `show` 输出格式

`user:password@host:port`，**明文**——小心屏幕被偷窥。

## Linux 平台特别说明

### 原生 Linux PuTTY 不是 Windows 版

`apt install putty` 装的 `putty` 是 GTK3 GUI 版的 PuTTY（上游 [sgtatham/putty](https://putty.org/) 的 Unix 移植）。它接受跟 Windows PuTTY 几乎一样的命令行参数（`-pw`、`-ssh`），所以工具可以无缝启动。

### Wayland 下 PuTTY 字体问题

GNOME Wayland 下，原生 Linux PuTTY 的 GDK 后端有 bug：字体枚举失败（`PuTTY: unable to load font "client:Ubuntu Mono 16"`）、glibc 符号冲突。**当前代码已自动处理**：在启动 PuTTY 子进程时设置 `GDK_BACKEND=x11` 强制走 X11 后端，**用户不需要任何手动配置**。

如果还报错，手动修法：

```bash
# 编辑 ~/.putty/sessions/Default Settings
# 把
FontName=client:Ubuntu Mono 16
# 改成（去掉 client: 前缀）
FontName=Ubuntu Mono 16
```

### Linux 下退出不卡

按 `exit` 或 `Ctrl-C` 退出后，**直接回到 shell**，不会再问"press return"（Windows 下保留这个行为以防控制台窗口闪关）。`stty sane` 这种手动修 shell 的需求也用不上了。

### （规划中）"新窗口"模式

当前实现是 PuTTY 跑在 ssh-mgr 的同一终端里。go-prompt v0.2.6 有一个老 bug：进入 raw mode 时改了 `icnl` / `igncr` 等 stty 属性，退出时没完全恢复——会导致**退出 ssh-mgr 后同一终端的 shell 看不见输入字符**（命令照常执行，但屏幕空）。

短期 workaround：

```bash
stty sane
```

**正在计划改成"新窗口"模式**：

```bash
gnome-terminal -- putty user@host -pw pass -ssh
```

这样 PuTTY 在独立的 GNOME Terminal 窗口里跑，ssh-mgr 的主终端完全不受影响。

## 故障排除

### 退出后同一终端 shell 看不见输入字符

go-prompt v0.2.6 已知 bug（[c-bata/go-prompt#228](https://github.com/c-bata/go-prompt/issues/228)、[#266](https://github.com/c-bata/go-prompt/issues/266)）。短期 `stty sane`；新窗口模式上线后不再需要。

### FileZilla 弹出后连不上 / 认证失败

最常见原因是**密码解密后多了 `\x00`**：历史版本用零填充（zero-padding）加密的连接信息，gmsm v1.4.1 的 `Sm4Ecb` 静默吞 PKCS#7 unpadding 错误（`out, _ = pkcs7UnPadding(out)`），解密结果是 `"root\x00\x00...\x00"`（16 字节），URL 编码后变 `root%00`，SSH 服务器直接拒。

**当前代码已自动处理**（`bytes.TrimRight(data, "\x00")`），不用手动修。

### PuTTY 报 "unable to load font"

参见上方 "Linux 平台特别说明 → Wayland 下 PuTTY 字体问题"。

### 输完密码后 REPL 里看不见字符

确保：

- 用最新版的 ssh-mgr（已替换 gopass → x/term）
- 终端是 xterm-compatible（GNOME Terminal / Konsole / xterm 都行）
- tmux/screen 里跑的话设 `default-terminal "xterm-256color"`

### `ssh start fail. err: exec: "..."` — 可执行文件找不到

`config.json` 的 `ssh` 字段写的是可执行文件**本体**，不是 shell 命令。Linux 配 `"ssh": "putty"` 即可（`$PATH` 里有 `/usr/bin/putty`），**不要**写 `"ssh": "env GDK_BACKEND=x1 putty"`——`exec.Command` 不走 shell 解析，整串会被当一个可执行文件名去找。环境变量注入在代码里做了（`cmd.Env = append(os.Environ(), "GDK_BACKEND=x11")`），config 里不要再加。

### 主密码忘了

没救。工具不存明文。建议：

- 用密码管理器存主密码
- 定期备份 `source` 目录和 `passwd` 文件

## 数据迁移

### Windows → Linux

把 Windows 上的 `source` 目录（默认 `C:\Linux\ssh`）和 `passwd` 文件整个 scp 到 Linux 任意目录，然后在 Linux 的 `config.json` 里把 `source` / `password` 指向这个目录。**主密码两边一致，不需要重新录入**——加密格式跨平台是一样的。

```bash
# Windows 上
scp -r C:\Linux\ssh C:\Linux\passwd user@linuxbox:~/.ssh-mgr/

# Linux 上编辑 config.json
{
    "ssh": "putty",
    "sftp": "filezilla",
    "source": "/home/yourname/.ssh-mgr/ssh",
    "isEncrypted": true,
    "password": "/home/yourname/.ssh-mgr/passwd"
}
```

### 备份

只备份两个东西：

1. `source` 目录（所有连接信息，加密的）
2. `passwd` 文件（主密码的 SM3 哈希）

```bash
tar czf ssh-mgr-backup.tgz -C ~/.ssh-mgr source passwd
```

## 安全说明

- **加密算法**：SM3（哈希）+ SM4（对称加密，ECB 模式）
- **KDF 弱点**：主密码 → SM4 key 的派生是 `md5(master_pwd)[:16]`（单轮 MD5，无 salt、无拉伸）。拿到 `passwd` 文件的攻击者可以离线爆破。建议主密码用 16+ 字符的强密码。改进方向：换 Argon2id / scrypt + 随机 salt
- **进程参数泄露**：启动 PuTTY 时密码通过 `-pw` 传命令行参数，在 Linux `ps` / `top` 或 Windows 任务管理器里能看到。PuTTY 协议本身的限制，不修
- **配置文件权限**：当前用 `os.ModePerm` (0777)，所有用户可读。建议 `chmod 700 ~/.ssh-mgr` 保护
- **主密码忘了 = 数据全丢**，没救

## 注意点

虽然连接信息用国密加密保存，但在使用时（启动 PuTTY / FileZilla），密码会通过命令行参数传递，**Windows 任务管理器 / Linux `ps` 都能看到**。建议只在单人使用的机器上跑。

## 开发者

- [AGENTS.md](AGENTS.md) — 开发指南（架构、依赖、调试注意事项、go-prompt 终端状态坑）
- 跨平台编译：`GOOS=windows go build`
- 代码组织：
  - `cmd/mgr/` — 主 REPL 二进制
  - `dao/` — 数据模型 + 加密 mapper（SM3 + SM4）
  - `services/` — ssh / sftp / ini / show / pass 业务逻辑
  - `api/cmdline/` — REPL glue（executor + completer）
