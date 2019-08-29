// +build linux

package parallel

import (
	"runtime"
	"sync/atomic"
)

type NoteAction func(id int, note interface{}) error

type SpinLock struct{ lock uintptr }

func (l *SpinLock) Lock() {
	for !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
	}
}
func (l *SpinLock) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}

type NoteTable struct {
	mu    SpinLock
	notes map[int]interface{}
}

func NewNoteTable() *NoteTable {
	return &NoteTable{notes: make(map[int]interface{})}
}

func (t *NoteTable) Get(id int) (note interface{}, ok bool) {
	t.mu.Lock()
	note, ok = t.notes[id]
	t.mu.Unlock()
	return
}

func (t *NoteTable) Add(id int, note interface{}) bool {
	t.mu.Lock()
	t.notes[id] = note
	n := len(t.notes)
	t.mu.Unlock()
	return n == 1
}

func (t *NoteTable) Del(id int) bool {
	t.mu.Lock()
	delete(t.notes, id)
	n := len(t.notes)
	t.mu.Unlock()
	return n == 0
}

func (t *NoteTable) Each(action NoteAction) error {
	t.mu.Lock()
	if len(t.notes) == 0 {
		t.mu.Unlock()
		return nil
	}
	notes := t.notes
	t.notes = nil
	t.mu.Unlock()
	for id, note := range notes {
		if err := action(id, note); err != nil {
			return err
		}
	}
	return nil
}
