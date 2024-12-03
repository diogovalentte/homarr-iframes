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

func TestGetVersion(t *testing.T) {
	v, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get version", func(t *testing.T) {
		version, err := v.getVikunjaVersion()
		if err != nil {
			t.Fatal(err)
		}

		if version == "" {
			t.Fatal("version is empty")
		}
	})

	tests := map[string]string{
		"1.25.3":  "1.25.3",
		"1.25.4":  "1.25.3",
		"1.25.0":  "1.24.3",
		"2.25.3":  "1.55.4",
		"24.25.3": "5.0.8",
	}

	t.Run("extract version", func(t *testing.T) {
		for version1, version2 := range tests {
			v, err := IsVersionGreaterOrEqualTo(version1, version2)
			if err != nil {
				t.Fatal(err)
			}
			if !v {
				t.Fatalf("%s is not greater or equal to %s", version1, version2)
			}
		}
	})
}

func TestGetTasks(t *testing.T) {
	v, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get tasks", func(t *testing.T) {
		tasks, err := v.GetTasks(-1, 0, []*int{})
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
	v, err := New()
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
