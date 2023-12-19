package tasks

import "testing"

func TestListPublishedTasks(t *testing.T) {
	tasks, err := ListPublishedTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) == 0 {
		t.Fatal("no tasks returned")
	}
}
