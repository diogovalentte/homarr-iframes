package vikunja

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

func setup() error {
	envFilePath := "../../../.env.test"
	err := config.SetConfigs(envFilePath)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetToken(t *testing.T) {
	configs := config.GlobalConfigs
	v := Vikunja{
		Address:  configs.VikunjaConfigs.Address,
		Username: configs.VikunjaConfigs.Username,
		Password: configs.VikunjaConfigs.Password,
	}

	t.Run("valid login", func(t *testing.T) {
		_, err := v.GetToken()
		if err != nil {
			t.Fatal(err)
			return
		}
	})

	t.Run("invalid password/username login", func(t *testing.T) {
		_, err := v.GetToken()
		if err == nil {
			t.Fatal(fmt.Errorf("Expected error while logging with invalid credentials"))
			return
		} else {
			if !strings.HasPrefix(err.Error(), "Login Error: unexpected status code: 412") {
				t.Fatal(fmt.Errorf("Expected status code while logging with invalid credentials: %s", err.Error()))
			}
		}
	})
}

func TestGetTasks(t *testing.T) {
	configs := config.GlobalConfigs
	v := Vikunja{
		Address:  configs.VikunjaConfigs.Address,
		Username: configs.VikunjaConfigs.Username,
		Password: configs.VikunjaConfigs.Password,
	}

	t.Run("get tasks", func(t *testing.T) {
		tasks, err := v.GetTasks(-1)
		if err != nil {
			t.Fatal(err)
			return
		}

		for _, task := range tasks {
			if task.ID == 0 {
				t.Fatal("task with ID 0")
				return
			}
		}
	})
}
