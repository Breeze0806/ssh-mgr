# ssh连接管理器
ssh连接管理器（ssh-mgr）是一个管理[putty](https://putty.org/) , [winscp](https://winscp.net/eng/index.php) 或者[filezilla](https://www.filezilla.cn/)的连接信息的工具

## 优点

1. 使用[putty](https://putty.org/) 进行ssh登录任一支持ssh的机器
2. 使用[winscp](https://winscp.net/eng/index.php) 或者[filezilla](https://www.filezilla.cn/)进行上传或下载文件
3. 使用命令行智能交互界面
4. 能够管理[putty](https://putty.org/) , [winscp](https://winscp.net/eng/index.php) 或者[filezilla](https://www.filezilla.cn/)的连接信息
5. 使用[putty](https://putty.org/) , [winscp](https://winscp.net/eng/index.php) 或者[filezilla](https://www.filezilla.cn/)进行免密登录
6. 使用**国密**加密保存用户名和密码的文件，密匙为登录密码

## 编译和安装

下载代码并编译：

```bash
cd cmd/mgr
go build
```

或者下载[ssh-mgr的windows64位发布版本](https://github.com/Breeze0806/ssh-mgr/releases/tag/v1.0.0)，注意其最小版本为win7或者win2012

## 使用方式

### 配置

```json
{
    "ssh": "C:\\Program Files\\PuTTY\\putty.exe",
    "sftp": "C:\\Program Files (x86)\\WinSCP\\WinSCP.exe",
    "source":"C:\\Linux\\ssh",
    "isEncrypted":true,
    "password":"C:\\Linux\\passwd"
}
```

- ssh 用于配置putty程序的路径
- sftp 用于配置 winscp 或者filezilla程序的路径
- source 用于配置存储ssh连接信息的路径
- isEncrypted 是否加密ssh连接信息的路径
- password 用于存储密码信息，用于加密登录

### 使用方法

```
Usage of mgr:
  -c string
        config file (default "config.json")
```

### 快速开始

- 在isEncrypted为ture时，开始需要输入密码，如果之前没有输入密码则会输入两次确认密码
- 使用下面的命令就可以进行ssh连接或者sftp连接，group是将ssh连接分组起的名称，而name是sh连接的别名

```bash
ssh group name      #启动putty进行ssh连接
sftp group name     #启动 winscp 或者filezilla进行sftp连接
add group name      #新增一个ssh连接信息
showAddr address    #显示ip:port相关的ssh连接信息
show group name     #显示对应分组1的ssh连接信息
exit                #退出程序
```

- 在打印对应的命令时会有对应的提醒，如下所示:

```bash
> ss
      ssh  ssh group name
```

- 另外，它提示出对应的提醒，可以使用tab键后上下选择

```
> ssh t
       test   test
       test1  test1
```

- 新增ssh连接信息，输入add group name后

```bash
#没有端口会默认为22，如果输入1.1.1.1:1234
please input ssh address:1.1.1.1 
please input ssh user:root
please input ssh password:*******
```

### 注意点

本工具虽然加密相关文件，但是在使用时，windows的任务管理器仍然会泄露shh的用户名和密码，为此在使用时确保只有一人在使用电脑。
