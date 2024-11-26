package ssh

import (
	"bytes"
	"fmt"
	"golang_ssp/golang_ssp/internal/config"
	"os"
	"os/exec"
	"syscall"
)

func Login(cfg *config.SSHConfig, cfgs *[]config.SSHConfig, configPath string) {

	ret := checkConnection(cfg)
	if !ret {
		panic("Connection failed")
	} else {
		updateConfigs(cfg, cfgs, configPath)
	}

	binary, err := exec.LookPath("sshpass")
	if err != nil {
		fmt.Printf("Error looking up sshpass: %v\n", err)
		os.Exit(1)
	}

	args := []string{"sshpass", "-p", cfg.Password, "ssh", "-p", cfg.Port, fmt.Sprintf("%s@%s", cfg.User, cfg.Hostname)}
	err = syscall.Exec(binary, args, os.Environ())
	if err != nil {
		fmt.Printf("Error executing sshpass: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully logged in!")
}

func updateConfigs(cfg *config.SSHConfig, cfgs *[]config.SSHConfig, configPath string) {
	// 增加登录次数&时间
	cfg.Increase()

	isExist := false
	for i, c := range *cfgs {
		if c.Host == cfg.Host {
			(*cfgs)[i].Update(cfg)
			isExist = true
			break
		}
	}
	if !isExist {
		*cfgs = append(*cfgs, *cfg)
	}

	config.SortConfigs(cfgs)
	if err := config.WriteConfig(configPath, *cfgs); err != nil {
		fmt.Println("Error writing config:", err)
	}
}

func checkConnection(cfg *config.SSHConfig) bool {
	// 删除 known_hosts 记录
	if cfg.Port == "" {
		cfg.Port = "22"
	}

	kownhost := cfg.Hostname

	if cfg.Port != "22" {
		kownhost = "[" + cfg.Hostname + "]:" + cfg.Port
	}

	execCmd := exec.Command("ssh-keygen", "-R", kownhost)
	var testErr bytes.Buffer

	execCmd.Stderr = &testErr
	if err := execCmd.Run(); err != nil {
		fmt.Printf("ssh-keygen remove failed: %v, stderr: \n\n%s\n", err, testErr.String())
		return false
	}

	// 尝试连接测试
	testCmd := exec.Command("sshpass", "-p", cfg.Password, "ssh", "-p", cfg.Port, fmt.Sprintf("%s@%s", cfg.User, cfg.Hostname), "true")
	fmt.Println(testCmd.String())
	testCmd.Stderr = &testErr

	if err := testCmd.Run(); err != nil {
		fmt.Printf("Connection test failed: %v, stderr: \n\n%s\n", err, testErr.String())
		return false
	}

	fmt.Println("Connection test passed, proceeding with login...")
	return true
}
