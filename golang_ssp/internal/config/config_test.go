package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {

	// 创建临时目录
	outputDir := "tmp"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(outputDir)

	// 创建临时配置文件
	configPath := outputDir + "/test_config.yaml"
	configContents := `Host test1
  HostName 192.168.1.1
  User ubuntu
  #Password abcdefg
  #LoginTimes 20
  #LastLoginTime 2023-07-01T10:00:00
Host test2
  HostName 192.168.1.2
  User ubuntu
  #Password abcdefg
  #LoginTimes 20
  #LastLoginTime 2023-07-01T12:00:00
Host test3
  HostName 192.168.1.3
  User ubuntu
  #Password abcdefg
  #LoginTimes 21
  #LastLoginTime 2023-07-01T12:00:00
Host test4
  HostName 192.168.1.4
  User ubuntu
  #Password abcdefg
`
	configFile, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer configFile.Close()

	configFile.WriteString(configContents)

	// 执行
	configs, err := ReadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	// 验证结果
	expectedConfigs := []SSHConfig{
		{Host: "test3", Hostname: "192.168.1.3", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "21", LastLoginTime: "2023-07-01T12:00:00"},
		{Host: "test2", Hostname: "192.168.1.2", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "20", LastLoginTime: "2023-07-01T12:00:00"},
		{Host: "test1", Hostname: "192.168.1.1", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "20", LastLoginTime: "2023-07-01T10:00:00"},
		{Host: "test4", Hostname: "192.168.1.4", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "0", LastLoginTime: ""},
	}

	for i, expectedConfig := range expectedConfigs {
		config := (*configs)[i]
		if config.Host != expectedConfig.Host ||
			config.Hostname != expectedConfig.Hostname ||
			config.User != expectedConfig.User ||
			config.Port != expectedConfig.Port ||
			config.Password != expectedConfig.Password ||
			config.LoginTimes != expectedConfig.LoginTimes ||
			config.LastLoginTime != expectedConfig.LastLoginTime {
			t.Errorf("Expected config %v, got %v", expectedConfig, config)
		}
	}

}

func TestWriteConfig(t *testing.T) {
	// 创建临时目录
	outputDir := "tmp"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(outputDir)

	configPath := outputDir + "/test_config.yaml"
	defer os.Remove(configPath)

	config := SSHConfig{
		Host:          "test",
		Hostname:      "test.com",
		User:          "testuser",
		Port:          "22",
		Password:      "testpassword",
		LoginTimes:    "0",
		LastLoginTime: "2023-07-01T12:00:00",
	}

	err := WriteConfig(configPath, []SSHConfig{config})
	if err != nil {
		t.Errorf("Error writing config: %v", err)
	}
	cfgs, err := ReadConfig(configPath)
	if err != nil || cfgs == nil || len(*cfgs) == 0 || !(*cfgs)[0].Equals(&config) {
		t.Error("Error reading config")
	}
}

func TestListConfigs(t *testing.T) {
	configs := []SSHConfig{
		{Host: "test3", Hostname: "192.168.1.3", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "21", LastLoginTime: "2023-07-01T12:00:00"},
		{Host: "test2", Hostname: "192.168.1.2", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "20", LastLoginTime: "2023-07-01T12:00:00"},
		{Host: "test1", Hostname: "192.168.1.1", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "20", LastLoginTime: "2023-07-01T10:00:00"},
		{Host: "test4", Hostname: "192.168.1.4", User: "ubuntu", Port: "22", Password: "abcdefg", LoginTimes: "0", LastLoginTime: ""},
	}

	ListConfigs(configs)
}

func TestGetSSHConfig(t *testing.T) {
	configs := []SSHConfig{
		{Host: "test1", Hostname: "192.168.1.1", User: "ubuntu", Password: "abcdefg", LoginTimes: "20", LastLoginTime: "2023-07-01T10:00:00"},
		{Host: "test2", Hostname: "192.168.1.2", User: "ubuntu", Password: "abc"},
	}
	cfg := SSHConfig{Host: "test1"}
	config, _ := GetSSHConfig(&configs, &cfg)

	if config.Host != "test1" {
		t.Errorf("Expected host to be 'test1', got '%s'", config.Host)
	}
}
