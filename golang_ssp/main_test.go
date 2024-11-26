package main

import (
	"golang_ssp/golang_ssp/internal/config"
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		args          []string
		expectedCfg   *config.SSHConfig
		expectedModel string
	}{
		{
			args:          []string{"-list"},
			expectedCfg:   &config.SSHConfig{},
			expectedModel: "list",
		},
		{
			args:          []string{"-host", "192.168.1.1"},
			expectedCfg:   &config.SSHConfig{Host: "192.168.1.1"},
			expectedModel: "login",
		},
		{
			args:          []string{"user@192.168.1.1"},
			expectedCfg:   &config.SSHConfig{Hostname: "192.168.1.1", User: "user"},
			expectedModel: "login",
		},
		{
			args:          []string{"192.168.1.1"},
			expectedCfg:   &config.SSHConfig{Host: "192.168.1.1", Hostname: "192.168.1.1"},
			expectedModel: "login",
		},
		{
			args:          []string{"1"},
			expectedCfg:   &config.SSHConfig{},
			expectedModel: "index",
		},
	}

	for _, tc := range testCases {
		// reset
		*listOpt = false
		*hostnameOpt = ""
		*hostOpt = ""

		// 模拟命令行参数
		oldArgs := os.Args
		os.Args = append([]string{"ssp"}, tc.args...)

		// 调用被测试函数
		model, data := ParseArgs()

		// 恢复命令行参数
		os.Args = oldArgs

		// 验证结果
		if data["config"] != nil && tc.expectedCfg != nil {

			if !data["config"].(*config.SSHConfig).Equals(tc.expectedCfg) {
				t.Errorf("Expected config: %v, got: %v", tc.expectedCfg, data)
			}
		} else if (data["config"] == nil) != (tc.expectedCfg == nil) {
			t.Errorf("Expected cfg to be %v, got %v", tc.expectedCfg, data)
		}

		if model != tc.expectedModel {
			t.Errorf("Expected model: %s, got: %s", tc.expectedModel, model)
		}
	}
}

func TestReadInput(*testing.T) {
	sshConfig := config.SSHConfig{}

	ReadInput(&sshConfig)
}
