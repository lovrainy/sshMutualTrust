package main

import (
	"github.com/spf13/cobra"
	"github.com/zhuah/socker"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"sshMutualTrust/configs"
	"sshMutualTrust/utils"
	"strings"
	"time"
)

var (
	cmd    *exec.Cmd
	output []byte
	err    error
)

type Host struct {
	IP string
	Port string
	User string
	Password string
}

func main() {
	// 初始化配置
	configs.InitConfig()

	var rootCmd = &cobra.Command{Use: "sshMutualTrust"}
	var cmdManyTrust = &cobra.Command{
		Use:   "many",
		Short: "多主机之间互信（配置manyNode.conf）.",
		Long:  `多主机之间互信（配置manyNode.conf）.`,
		Run: func(cmd *cobra.Command, args []string) {
			configs.Logger.Info("开始配置主机互信!")

			// 多主机连接配置
			authKeys := make([]string, 0)
			key := make(chan string, 0)

			ipDict, err := NodeConfig("configs/manyNode.conf")
			if err != nil {
				configs.Logger.Error(err)
			}

			countSelect := 0
			for _, h := range ipDict {
				go SelectSecret(h, key)

			}
			loopSelect:
				for {
					select {
					case k := <- key:
						authKeys = append(authKeys, k)
						countSelect += 1
						if countSelect == len(ipDict) {
							break loopSelect
						}
					case <-time.After(time.Duration(configs.ConnTimeout) * time.Second):
						break loopSelect
					}
				}
			// 将所有密钥写入到主机临时文件
			down := make(chan bool, 0)
			countWrite := 0
			for _, h := range ipDict {
				go SecretWrite(h, authKeys, down)

			}
			loopWrite:
				for {
					select {
					case <- down:
						countWrite += 1
						if countWrite == len(ipDict) {
							break loopWrite
						}
					case <-time.After(time.Duration(configs.ConnTimeout) * time.Second):
						break loopWrite
					}
				}
		},
	}

	var cmdSingleTrust = &cobra.Command{
		Use:   "single",
		Short: "单节点互信多节点（配置singleNode.conf和manyNode.conf）.",
		Long:  `单节点互信多节点（配置singleNode.conf和manyNode.conf）.`,
		Run: func(cmd *cobra.Command, args []string) {
			configs.Logger.Info("开始配置主机互信!")

			// 读取需要互信给多节点的单个节点信息（可以是多个）
			singleKeys := make([]string, 0)
			key := make(chan string, 0)

			singleIpList, err := NodeConfig("configs/singleNode.conf")
			if err != nil {
				configs.Logger.Error(err)
			}

			countSingle := 0
			for _, h := range singleIpList {
				go SelectSecret(h, key)
			}
			loopSelect:
				for {
					select {
					case keySingle := <- key:
						singleKeys = append(singleKeys, keySingle)
						countSingle += 1
						if countSingle == len(singleIpList) {
							break loopSelect
						}
					case <-time.After(time.Duration(configs.ConnTimeout) * time.Second):
						break loopSelect
					}
				}

			// 获取多节点连接信息
			manyIpList, err := NodeConfig("configs/manyNode.conf")
			if err != nil {
				configs.Logger.Error(err)
			}
			// 连接多节点，将互信密钥写入
			down := make(chan bool, 0)
			countWrite := 0
			for _, h := range manyIpList {
				go SecretWrite(h, singleKeys, down)

			}
			loopWrite:
				for {
					select {
					case <- down:
						countWrite += 1
						if countWrite == len(manyIpList) {
							break loopWrite
						}
					case <-time.After(time.Duration(configs.ConnTimeout) * time.Second):
						break loopWrite
					}
				}
		},
	}

	rootCmd.AddCommand(cmdManyTrust)
	rootCmd.AddCommand(cmdSingleTrust)

	err := rootCmd.Execute()
	if err != nil {
		configs.Logger.Error(err)
	}

}


// 获取主机连接配置
func NodeConfig(conf string) ([]Host, error) {
	ipDict := make([]Host, 0)
	// 解析manyNode.conf配置文件
	filePath := path.Join(utils.AbsPath(), conf)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		configs.Logger.Error(err)
		return ipDict, nil
	}
	ipuserList := strings.Split(string(bytes),"\n")

	for _, ipUser := range ipuserList {
		if ipUser == "" {
			break
		}
		l := strings.Fields(strings.TrimSpace(ipUser))
		h := Host{l[0], l[1],l[2],l[3]}
		ipDict = append(ipDict, h)
	}

	return ipDict, nil
}

// 远程主机密钥搜集
func SelectSecret(h Host, key chan string) {
	var sshConfig = &socker.Auth{User: h.User, Password: h.Password}

	agent, err1 := socker.Dial(h.IP+":"+h.Port, sshConfig)

	if err1 != nil {
		configs.Logger.Error("dial agent failed:", err1)
	} else {
		configs.Logger.Infof("连接主机[%s]成功！", h.IP+":"+h.Port)
		configs.Logger.Infof("检查主机[%s]是否存在密钥，不存在将创建！", h.IP+":"+h.Port)
		agent.Rcmd("if [ ! -f ~/.ssh/id_rsa ]; then ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa;fi")
		agent.Rcmd("if [ ! -f ~/.ssh/authorized_keys ]; then touch ~/.ssh/authorized_keys;fi")
		configs.Logger.Infof("主机[%s]密钥检查完毕！", h.IP+":"+h.Port)
		configs.Logger.Infof("正在搜集主机[%s]密钥！", h.IP+":"+h.Port)

		r, _ := agent.Rcmd("cat ~/.ssh/id_rsa.pub")
		configs.Logger.Infof("主机[%s]密钥搜集完毕！", h.IP+":"+h.Port)
		agent.Close()
		key <- string(r)
	}
}

// 写入互信密钥到远程主机
func SecretWrite(h Host, authKeys []string, down chan bool) {
	var sshConfig = &socker.Auth{User: h.User, Password: h.Password}

	agent, err1 := socker.Dial(h.IP+":"+h.Port, sshConfig)

	if err1 != nil {
		configs.Logger.Error("dial agent failed:", err1)
		down <- false
	} else {
		agent.Rcmd("if [ ! -f ~/.ssh/auth.tmp ]; then touch ~/.ssh/auth.tmp;fi")
		agent.Rcmd("> ~/.ssh/auth.tmp")
		for _, k := range authKeys {
			str := fmt.Sprintf(`echo "%s" >> ~/.ssh/auth.tmp`, strings.Trim(k, "\n"))
			agent.Rcmd(str)
		}

		agent.Rcmd("awk '{print $0}' ~/.ssh/auth.tmp ~/.ssh/authorized_keys|sort | uniq > ~/.ssh/auth.tmp2")
		if !configs.StrictHostKeyChecking {
			agent.Rcmd("if [ ! -f ~/.ssh/config ]; then touch ~/.ssh/config;fi")
			configs.Logger.Warnf("主机[%s]更新.ssh/config", h.IP+":"+h.Port)
			agent.Rcmd("if ! grep 'StrictHostKeyChecking' ~/.ssh/config >/dev/null; then echo 'StrictHostKeyChecking no' > ~/.ssh/config; else sed -i 's#^StrictHostKeyChecking.*$#StrictHostKeyChecking no#g' ~/.ssh/config;fi")
			agent.Rcmd("chmod 0600 ~/.ssh/config")
		}
		agent.Rcmd("cat ~/.ssh/auth.tmp2 > ~/.ssh/authorized_keys")
		agent.Rcmd("chown " + h.User + " ~/.ssh/authorized_keys")
		agent.Rcmd("chmod 0600 ~/.ssh/authorized_keys")


		configs.Logger.Infof("主机[%s]写入互信密钥完毕！", h.IP+":"+h.Port)
		agent.Rcmd("rm -f ~/.ssh/auth.tmp")
		agent.Rcmd("rm -f ~/.ssh/auth.tmp2")
		configs.Logger.Infof("主机[%s]互信中间缓存文件清理完毕！", h.IP+":"+h.Port)

		agent.Close()
		down <- true
	}
}