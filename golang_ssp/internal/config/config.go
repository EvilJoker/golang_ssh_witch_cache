package config

import (
	"bufio"
	"fmt"
	"golang_ssp/golang_ssp/pkg/logger"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SSHConfig struct {
	Host          string
	Hostname      string
	User          string
	Port          string
	Password      string // Not recommended to store passwords in plain text
	LoginTimes    string
	LastLoginTime string // 2022-01-01T15:04:05
}

const TIMEFORMAT = "2006-01-02T15:04:05"

func (s *SSHConfig) Equals(s2 *SSHConfig) bool {
	return s.Host == s2.Host && s.Hostname == s2.Hostname && s.User == s2.User && s.Port == s2.Port && s.Password == s2.Password && s.LoginTimes == s2.LoginTimes && s.LastLoginTime == s2.LastLoginTime
}

func (s *SSHConfig) Compare(s2 *SSHConfig) bool {
	// high than s2 ,return true

	// 谁最新谁优先
	if s.LastLoginTime != "" && s2.LastLoginTime != "" && s.LastLoginTime != s2.LastLoginTime {

		aTime, err1 := time.Parse(TIMEFORMAT, s.LastLoginTime)
		bTime, err2 := time.Parse(TIMEFORMAT, s2.LastLoginTime)
		if err1 != nil && err2 == nil {
			return false
		}
		if err2 != nil && err1 == nil {
			return true
		}

		return aTime.After(bTime)
	}

	// 谁登陆次数多谁优先
	if s.LoginTimes != s2.LastLoginTime {
		ai, _ := strconv.Atoi(s.LoginTimes)
		bi, _ := strconv.Atoi(s2.LoginTimes)

		return bi < ai

	}

	return strings.Compare(s2.Host, s.Host) < 0
}

func (s *SSHConfig) Update(s2 *SSHConfig) {
	s.Hostname = s2.Hostname
	s.Host = s2.Host
	s.User = s2.User
	s.Password = s2.Password
	s.Port = s2.Port
	s.LastLoginTime = s2.LastLoginTime
	s.LoginTimes = s2.LoginTimes
}
func (s *SSHConfig) String() string {
	if s.LoginTimes == "" {
		s.LoginTimes = "0"
	}

	if s.LastLoginTime == "" {
		s.LastLoginTime = "1977-01-01T15:04:05"
	}

	if s.Port == "" {
		s.Port = "22"
	}

	return fmt.Sprintf("Host %s\n  HostName %s\n  User %s\n  Port %s\n  #Password %s\n  #LoginTimes %s\n  #LastLoginTime %s\n", s.Host, s.Hostname, s.User, s.Port, s.Password, s.LoginTimes, s.LastLoginTime)
}

func (s *SSHConfig) Increase() {
	// 增加登录次数
	if s.LoginTimes == "" {
		s.LoginTimes = "0"
	}
	times, _ := strconv.Atoi(s.LoginTimes)

	s.LoginTimes = strconv.Itoa(times + 1)

	// 更新上次登录时间
	s.LastLoginTime = time.Now().Format("2006-01-02T15:04:05")

}

// ReadConfig 读取SSH配置文件并解析成SSHConfig结构体的切片。
// 它打开SSH配置文件，逐行读取并解析配置，然后关闭文件。
// 返回值是解析后的SSH配置切片和可能出现的错误。
func ReadConfig(configPath string) (*[]SSHConfig, error) {

	configPath = absPath(configPath)

	if _, err := os.Stat(configPath); err != nil {

		if os.IsNotExist(err) {
			logger.Logger.Printf("SSH config file not found, try to create %s\n", configPath)
			if _, err = os.Create(configPath); err != nil {
				logger.Logger.Printf("Create SSH config file failed, %s\n", err)
				return nil, err
			}
		}
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var configs []SSHConfig
	var currentConfig SSHConfig = SSHConfig{Password: "", LoginTimes: "0", LastLoginTime: ""}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") &&
			!strings.Contains(line, "Password") &&
			!strings.Contains(line, "LoginTimes") &&
			!strings.Contains(line, "LastLoginTime") {
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "#"))

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		switch key {
		case "Host":
			if currentConfig.Host != "" {
				configs = append(configs, currentConfig)
			}
			currentConfig = SSHConfig{Host: value, Port: "22", Password: "", LoginTimes: "0", LastLoginTime: ""}
		case "HostName":
			currentConfig.Hostname = value
		case "User":
			currentConfig.User = value
		case "Port":
			currentConfig.Port = value
		case "Password":
			currentConfig.Password = value
		case "LastLoginTime":
			currentConfig.LastLoginTime = value
		case "LoginTimes":
			currentConfig.LoginTimes = value
		}
	}

	if currentConfig.Host != "" {
		configs = append(configs, currentConfig)
	}
	SortConfigs(&configs)

	return &configs, scanner.Err()
}

func absPath(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Logger.Fatalf("获取用户主目录失败: %v\n", err)
		}
		// 替换 `~` 为用户主目录
		path = filepath.Join(homeDir, path[1:])

	}
	abs, _ := filepath.Abs(path)
	return abs
}

func SortConfigs(configs *[]SSHConfig) {
	sort.Slice(*configs, func(i, j int) bool {
		return (*configs)[i].Compare(&(*configs)[j])
	})
}

func WriteConfig(configPath string, configs []SSHConfig) error {
	configPath = absPath(configPath)
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, config := range configs {
		_, err := writer.WriteString(config.String())
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func ListConfigs(configs []SSHConfig) {
	if len(configs) == 0 {
		logger.Logger.Println("No configurations found")
		return
	}
	for i, config := range configs {
		fmt.Printf("%-4d%-25s%-15s%-15s%-15s%-15s%-15s\n", i, config.Host, config.Hostname, config.User, config.Port, config.LoginTimes, config.LastLoginTime)
		if i > 20 {
			break
		}
	}
}

func GetSSHConfig(c *[]SSHConfig, t *SSHConfig) (*SSHConfig, error) {
	for _, config := range *c {
		if t.Host != "" && config.Host == t.Host {
			return &config, nil
		}
	}
	for _, config := range *c {
		if t.Hostname != "" && config.Hostname == t.Hostname {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("no config found for host %v", t)
}
