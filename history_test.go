package prompt

import (
	"reflect"
	"testing"
)

func TestHistoryClear(t *testing.T) {
	h := NewHistory()
	h.Add("foo")
	h.Clear()
	expected := &History{
		histories: []string{"foo"},
		tmp:       []string{"foo", ""},
		selected:  1,
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}
}

func TestHistoryAdd(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	expected := &History{
		histories: []string{"echo 1"},
		tmp:       []string{"echo 1", ""},
		selected:  1,
	}
	if !reflect.DeepEqual(h, expected) {
		t.Errorf("Should be %v, but got %v", expected, h)
	}
}

func TestHistoryOlder(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")

	// Prepare buffer
	buf := NewBuffer()
	buf.InsertTextMoveCursor("echo 2", DefColCount, DefRowCount, false)

	// [1 time] Call Older function
	buf1, changed := h.Older(buf, DefColCount, DefRowCount)
	if !changed {
		t.Error("Should be changed history but not changed.")
	}
	if buf1.Text() != "echo 1" {
		t.Errorf("Should be %#v, but got %#v", "echo 1", buf1.Text())
	}

	// [2 times] Call Older function
	buf = NewBuffer()
	buf.InsertTextMoveCursor("echo 1", DefColCount, DefRowCount, false)
	buf2, changed := h.Older(buf, DefColCount, DefRowCount)
	if changed {
		t.Error("Should be not changed history but changed.")
	}
	if !reflect.DeepEqual("echo 1", buf2.Text()) {
		t.Errorf("Should be %#v, but got %#v", "echo 1", buf2.Text())
	}
}
