package jambel

import (
	"errors"
	"reflect"
	"testing"
)

type mockConn struct {
	// sent keeps track of sent bytes
	sent [][]byte
	// fail indicates whether Send should return an error
	fail bool
}

func (m *mockConn) Write(cmd []byte) (int, error) {
	if m.fail {
		return 0, errors.New("send error")
	}
	m.sent = append(m.sent, cmd)
	return len(cmd), nil
}

func (m *mockConn) Close() {}

func TestJambel_Reset(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	err := j.Reset()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(mc.sent) != 1 || string(mc.sent[0]) != "reset\n" {
		t.Errorf("expected 'reset\n', got %v", mc.sent)
	}
}

func TestJambel_On(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	_ = j.On(Green)
	want := []byte("set=3,on\n")
	if !reflect.DeepEqual(mc.sent[0], want) {
		t.Errorf("expected %q, got %q", want, mc.sent[0])
	}
}

func TestJambel_Off(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	_ = j.Off(Red)
	want := []byte("set=1,off\n")
	if !reflect.DeepEqual(mc.sent[0], want) {
		t.Errorf("expected %q, got %q", want, mc.sent[0])
	}
}

func TestJambel_Blink(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	_ = j.Blink(Yellow)
	want := []byte("set=2,blink\n")
	if !reflect.DeepEqual(mc.sent[0], want) {
		t.Errorf("expected %q, got %q", want, mc.sent[0])
	}
}

func TestJambel_BlinkInverse(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	_ = j.BlinkInverse(Green)
	want := []byte("set=3,blink_inverse\n")
	if !reflect.DeepEqual(mc.sent[0], want) {
		t.Errorf("expected %q, got %q", want, mc.sent[0])
	}
}

func TestJambel_Flash(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	_ = j.Flash(Red)
	want := []byte("set=1,flash\n")
	if !reflect.DeepEqual(mc.sent[0], want) {
		t.Errorf("expected %q, got %q", want, mc.sent[0])
	}
}

func TestJambel_SetAll(t *testing.T) {
	mc := &mockConn{}
	j := &Jambel{conn: mc}
	_ = j.SetAll(On, Blink, Off)
	want := []byte("set_all=0,2,1,0\n")
	if !reflect.DeepEqual(mc.sent[0], want) {
		t.Errorf("expected %q, got %q", want, mc.sent[0])
	}
}

func TestJambel_SendError(t *testing.T) {
	mc := &mockConn{fail: true}
	j := &Jambel{conn: mc}
	if err := j.On(Green); err == nil {
		t.Error("expected error, got nil")
	}
}
