# sshMutualTrust说明文档
## 一、简介
> 本项目用于方便主机(间)的互信操作。

## 二、项目结构
```shell
$ tree sshMutualTrustBin/
sshMutualTrustBin
├── configs
│   ├── manyNode.conf      ## 多个主机的连接信息
│   ├── settings.ini       ## 配置文件
│   └── singleNode.conf    ## 单个主机的连接信息（可以填写多个）
├── sshMutualTrust         ## 执行脚本（Linux AMD64） 
└── sshMutualTrustMac      ## 执行脚本（MacOS）

1 directory, 5 files
```

## 三、配置说明
### 3.1 `configs/settings.ini`
```vim
[default]
## ssh操作的超时时间
connection_timeout = 60
## 为false时相当于在~/.ssh/config配置'StrictHostKeyChecking no'参数（即首次连接无需输入yes）
strict_hostkey_checking = false

[logging]
## 日志级别
level = debug
development = true
## 是否开启celler
disable_caller = true
## 是否开启stacktrace
disable_stacktrace = true
## 日志格式
encoding = json
## 日志时间的key。如： {"time":"2019-10-23 09:46:34"}
encoder_config_time_key = time
## 日志level的key。 如： {"level": "info"}
encoder_config_level_key = level
encoder_config_name_key = log
## celler的key。 如： {"celler": "xxx"}
encoder_config_caller_key = caller
## msg的key。 如： {"msg": "xxx"}
encoder_config_msg_key = msg
## stacktrace名称
encoder_config_trace_key = stacktrace
## 日志文件的输出位置（stdout为console；其他为日志文件路径）
output_paths = stdout
## 错误日志的输出位置（stderr为console；其他为日志文件路径）
error_output_paths = stderr
## 自定义的key
initial_fields_key = service
## 自定义的value。与initial_fields_key打印如： {"service": "SSHMutualTrust"}
initial_fields_value = SSHMutualTrust
## 是否开启日志切割, 打开之后只会识别rotate_file, 如果需要打印console，请打开rotate_console
enable_rotate = true
## 屏幕打印日志
rotate_console = true
## 日志切割文件
rotate_file = ./logs/ssh.log
## 日志切割最大大小
rotate_max_size = 1024
## 日志最大保存份数
rotate_max_backups = 5
## 日志最大保存时间
rotate_max_age = 30
## 是否开启日志压缩
rotate_compress = false
```

### 3.2 `singleNode.conf`
> 此配置文件在操作`./sshMutualTrust single`时需要配置。

> 格式：`ip port username password`
```vim
127.0.0.1 22 dongxiaoyi xxxxxx
192.168.75.174 22 root redhat
```
说明：
- 配置中的主机的密钥将会分发给`manyNode.conf`中的主机。

### 3.3 `manyNode.conf`
> 此配置文件在操作`./sshMutualTrust single`和`./sshMutualTrust many`时都需要配置。
> 

> 格式：`ip port username password`
```vim
192.168.75.174 22 root redhat
192.168.75.175 22 root redhat
```
说明：
- 在操作`./sshMutualTrust single`时每个主机都会接受`singleNode.conf`中的主机的密钥。
- 在操作`./sshMutualTrust many`时每个主机之间形成互信。

## 四、互信操作
### 4.1 操作指令说明
```shell
$ ./sshMutualTrust --help                                 
Usage:
  sshMutualTrust [command]

Available Commands:
  help        Help about any command
  many        多主机之间互信（配置manyNode.conf）.
  single      单节点互信多节点（配置singleNode.conf和manyNode.conf）.

Flags:
  -h, --help   help for sshMutualTrust

Use "sshMutualTrust [command] --help" for more information about a command.
```

### 4.2 多主机间互信
(1) 配置`manyNode.conf`
example:
```vim
192.168.75.174 22 root redhat
192.168.75.175 22 root redhat
```
(2) 执行`./sshMutualTrust many`
example:
```shell
$ ./sshMutualTrust many
{"level":"info","time":"2020-03-13 20:29:40","msg":"开始配置主机互信!","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"连接主机[192.168.75.175:22]成功！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"检查主机[192.168.75.175:22]是否存在密钥，不存在将创建！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"连接主机[192.168.75.174:22]成功！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"检查主机[192.168.75.174:22]是否存在密钥，不存在将创建！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"主机[192.168.75.175:22]密钥检查完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"正在搜集主机[192.168.75.175:22]密钥！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"主机[192.168.75.174:22]密钥检查完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"正在搜集主机[192.168.75.174:22]密钥！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"主机[192.168.75.175:22]密钥搜集完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:29:50","msg":"主机[192.168.75.174:22]密钥搜集完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:30:00","msg":"主机[192.168.75.175:22]写入互信密钥完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:30:00","msg":"主机[192.168.75.175:22]互信中间缓存文件清理完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:30:00","msg":"主机[192.168.75.174:22]写入互信密钥完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:30:00","msg":"主机[192.168.75.174:22]互信中间缓存文件清理完毕！","service":"SSHMutualTrust"}
```

### 4.3 单节点互信多节点
> 多个节点将接受单个节点的ssh密钥

(1) 配置`singleNode.conf`
example:
```vim
127.0.0.1 22 dongxiaoyi xxxxxx
192.168.75.111 22 root redhat
```
(2) 配置`manyNode.conf`
example:
```vim
192.168.75.174 22 root redhat
192.168.75.175 22 root redhat
```
(3)执行`./sshMutualTrust single`
```shell
$ ./sshMutualTrust single
{"level":"info","time":"2020-03-13 20:31:31","msg":"开始配置主机互信!","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:31","msg":"连接主机[127.0.0.1:22]成功！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:31","msg":"检查主机[127.0.0.1:22]是否存在密钥，不存在将创建！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:31","msg":"主机[127.0.0.1:22]密钥检查完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:31","msg":"正在搜集主机[127.0.0.1:22]密钥！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:31","msg":"主机[127.0.0.1:22]密钥搜集完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:41","msg":"连接主机[192.168.75.111:22]成功！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:41","msg":"检查主机[192.168.75.111:22]是否存在密钥，不存在将创建！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:41","msg":"主机[192.168.75.111:22]密钥检查完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:41","msg":"正在搜集主机[192.168.75.111:22]密钥！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:41","msg":"主机[192.168.75.111:22]密钥搜集完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:41","msg":"开始配置主机互信!","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:51","msg":"主机[192.168.75.175:22]写入互信密钥完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:51","msg":"主机[192.168.75.175:22]互信中间缓存文件清理完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:51","msg":"主机[192.168.75.174:22]写入互信密钥完毕！","service":"SSHMutualTrust"}
{"level":"info","time":"2020-03-13 20:31:51","msg":"主机[192.168.75.174:22]互信中间缓存文件清理完毕！","service":"SSHMutualTrust"}
```
说明：
- 操作完成后`192.168.75.174`和`192.168.75.175`的`~/.ssh/authorized_keys`中会添加`127.0.0.1`、`192.168.75.111`的`~/.ssh/id_rsa.pub`

### 4.4 日志
> 操作会在项目生成日志，路径`tmp/ssh.log`
