package vikunja

import (
	"fmt"
	"os"
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

func TestGetTasks(t *testing.T) {
	configs := config.GlobalConfigs
	v := Vikunja{
		Address: configs.VikunjaConfigs.Address,
		Token:   configs.VikunjaConfigs.Token,
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
