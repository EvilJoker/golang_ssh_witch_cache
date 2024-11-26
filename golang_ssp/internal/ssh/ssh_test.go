package ssh

import (
	"golang_ssp/golang_ssp/internal/config"
	"log"
	"testing"

	ssh3 "github.com/gliderlabs/ssh"
)

func TestCheckConnection(t *testing.T) {
	// 启动 sshd.Server
	go func() {
		ssh3.Handle(func(s ssh3.Session) {
			s.Write([]byte("Welcome to local SSH server!\n"))
			s.Exit(0)
		})

		log.Println("Starting SSH server on 127.0.0.1:2222...")
		err := ssh3.ListenAndServe("127.0.0.1:2222", nil, ssh3.PasswordAuth(func(ctx ssh3.Context, password string) bool {
			return ctx.User() == "test" && password == "1234" // 用户名和密码验证
		}))
		if err != nil {
			log.Fatal(err)
		}

	}()
	// time.Sleep(2 * time.Second)
	// 测试
	// 登录正常
	cfg := &config.SSHConfig{
		Host:          "test",
		Hostname:      "127.0.0.1",
		User:          "test",
		Port:          "2222",
		Password:      "1234",
		LoginTimes:    "0",
		LastLoginTime: "2023-07-01T12:00:00",
	}

	if !checkConnection(cfg) {
		t.Fatalf("test Connection right failed")
	}
	// 登录异常

	cfg.Password = "12345"

	if checkConnection(cfg) {
		t.Fatalf("test Connection wrong failed")
	}
}
