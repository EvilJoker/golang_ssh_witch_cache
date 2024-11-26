package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang_ssp/golang_ssp/internal/config"
	"golang_ssp/golang_ssp/internal/ssh"
	"golang_ssp/golang_ssp/pkg/logger"
	"os"
	"strconv"
	"strings"
)

var cacheConfigPath = "~/.ssh/config_cache"

var (
	hostOpt = flag.String("host", "", "SSH host to connect")
	// ssp -hostOpt ip/hostname
	hostnameOpt = flag.String("hostname", "", "SSH hostname to connect")
	// ssp -list
	listOpt = flag.Bool("list", false, "List cached hosts")
)

func ParseArgs() (string, map[string]interface{}) {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of ssp (depends on sshpaas):\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  Description:\n")
		fmt.Fprintf(flag.CommandLine.Output(), `    ssp could simplify ssh login that auto compleled info by finding and caching ssh record,
    all record cache in ~/.ssh/config_cache`)
		fmt.Fprintf(flag.CommandLine.Output(), "\n\nssp [options] [host]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -host string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "     SSH host to connect (e.g., ssp -host node1)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -hostname string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "     SSH hostname to connect (e.g., ssp -hostname 127.0.0.1)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -list\n")
		fmt.Fprintf(flag.CommandLine.Output(), "     List cached hosts (e.g., ssp -list)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  host/hostname\n")
		fmt.Fprintf(flag.CommandLine.Output(), "     Same as -host -hostname, but needn`t '-' (e.g., ssp node1 or ssp 127.0.0.1 )\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  user@hostname\n")
		fmt.Fprintf(flag.CommandLine.Output(), "     Like ssh command (e.g., ssp root@127.0.0.1 )\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  index\n")
		fmt.Fprintf(flag.CommandLine.Output(), "     Use inde of '-list' reusult to login  (e.g., ssp 2, meaning use 2nd host in cache )\n")
	}

	flag.Parse()

	data := map[string]interface{}{}

	if *listOpt {
		data["config"] = &config.SSHConfig{}
		return "list", data
	}

	if *hostOpt != "" {
		data["config"] = &config.SSHConfig{Host: *hostOpt}
		return "login", data
	}

	if *hostnameOpt != "" {
		data["config"] = &config.SSHConfig{Hostname: *hostnameOpt}
		return "login", data
	}
	// 检查是否有非标志参数
	args := flag.Args()
	if len(args) > 0 {
		// 解析非标志参数
		if strings.Contains(args[0], "@") {
			parts := strings.Split(args[0], "@")
			if len(parts) == 2 {
				user := strings.TrimSpace(parts[0])
				hostname := strings.TrimSpace(parts[1])
				data["config"] = &config.SSHConfig{Hostname: hostname, User: user}
				return "login", data

			} else {
				fmt.Println("Invalid format for non-flag argument. Expected 'user@host'.")
				os.Exit(1)
			}
		} else if isInt(args[0]) {
			data["config"] = &config.SSHConfig{}
			data["index"] = args[0]
			return "index", data

		} else {
			data["config"] = &config.SSHConfig{Host: args[0], Hostname: args[0]}
			return "login", data
		}
	}
	fmt.Println("Invalid number of arguments.")
	panic("Invalid number of arguments.")
}

func isInt(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false

}

func ReadInput(cfg *config.SSHConfig) *config.SSHConfig {

	if cfg == nil {
		cfg = &config.SSHConfig{}
	}
	reader := bufio.NewReader(os.Stdin)

	if cfg.Host == "" {
		fmt.Print("Enter Host Like \"node1\": ")
		host, _ := reader.ReadString('\n')
		cfg.Host = strings.TrimSpace(host)

		if cfg.Host == "" {
			fmt.Println("Host cannot be empty.")
			os.Exit(1)
		}
	}

	if cfg.Hostname == "" {
		fmt.Print("Enter Hostname Like \"127.0.0.1\": ")
		hostname, _ := reader.ReadString('\n')
		cfg.Hostname = strings.TrimSpace(hostname)

		if cfg.Hostname == "" {
			fmt.Println("Hostname cannot be empty.")
			os.Exit(1)
		}
	}

	if cfg.User == "" {
		fmt.Print("Enter Username (default \"root\") ")
		user, _ := reader.ReadString('\n')
		cfg.User = strings.TrimSpace(user)

		if cfg.User == "" {
			cfg.User = "root"
		}
	}

	if cfg.Password == "" {
		fmt.Print("Enter Password : ")
		password, _ := reader.ReadString('\n')
		cfg.Password = strings.TrimSpace(password)

		if cfg.Password == "" {
			fmt.Println("Hostname cannot be empty.")
			os.Exit(1)
		}
	}

	if cfg.Port == "" {
		fmt.Print("Enter Port (default 22) ")
		port, _ := reader.ReadString('\n')
		cfg.Port = strings.TrimSpace(port)

		if cfg.Port == "" {
			cfg.Port = "22"
		}
	}

	return cfg

}

func main() {
	logger.Logger.Println("ssp start!")
	model, data := ParseArgs()
	inputCfg := data["config"].(*config.SSHConfig)

	cfgs, err := config.ReadConfig(cacheConfigPath)

	if err != nil {
		fmt.Printf("Error reading cache config: %v\n", err)
		os.Exit(1)
	}

	switch model {
	case "list":

		config.ListConfigs(*cfgs)

	case "login":
		if inputCfg == nil {
			fmt.Println("Invalid number of arguments.")
			panic("Invalid number of arguments.")
		}

		cfg, err := config.GetSSHConfig(cfgs, inputCfg)
		// 从命令行获取配置
		if err != nil {

			cfg = ReadInput(cfg)
			if cfg == nil {
				fmt.Printf("Error getting SSH config: %v\n", err)
				panic("Invalid number of arguments.")
			}
			fmt.Println(cfg)
		}

		ssh.Login(cfg, cfgs, cacheConfigPath)
	case "index":

		index, _ := strconv.Atoi(data["index"].(string))

		if index < 0 || index >= len(*cfgs) {
			fmt.Printf("Invalid index: %d\n", index)
		}

		cfg := (*cfgs)[index]

		ssh.Login(&cfg, cfgs, cacheConfigPath)

	}

}
