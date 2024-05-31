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
	v, err := New(config.GlobalConfigs.Vikunja.Address, config.GlobalConfigs.Vikunja.Token)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get tasks", func(t *testing.T) {
		tasks, err := v.GetTasks(-1, 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, task := range tasks {
			if task.ID == 0 {
				t.Fatal("task with ID 0")
			}
		}
	})
}

func TestSetTaskDone(t *testing.T) {
	v, err := New(config.GlobalConfigs.Vikunja.Address, config.GlobalConfigs.Vikunja.Token)
	if err != nil {
		t.Fatal(err)
	}
	taskID := 0

	t.Run("set task done (need to manually set a valid task ID!!!)", func(t *testing.T) {
		err := v.SetTaskDone(taskID)
		if err != nil {
			t.Fatal(err)
			return
		}
	})
}
