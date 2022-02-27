# ssh-mgr
ssh-mgr是一个 管理putty , winscp 或者filezilla的ssh连接的工具

## cmdline

### compile & install

下载代码并编译：

```bash
cd cmd/mgr
go build
```

### configure

```json
{
    "ssh": "F:\\Program Files\\PuTTY\\putty.exe",
    "sftp": "C:\\Program Files (x86)\\FileZilla FTP Client\\filezilla.exe",
    "source":"D:\\Linux",
    "isEncrypted":false
}
```

- ssh 用于配置putty程序的路径
- sftp 用于配置 winscp 或者filezilla程序的路径
- source 用于配置存储ssh连接信息的路径
- isEncrypted 是否加密ssh连接信息的路径
- password 用于存储密码信息，用于加密登录

### usage

```
Usage of mgr:
  -c string
        config file (default "config.json")
```

### start

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

