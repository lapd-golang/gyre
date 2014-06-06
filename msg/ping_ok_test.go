package msg

import (
	zmq "github.com/pebbe/zmq4"

	"testing"
)

// Yay! Test function.
func TestPingOk(t *testing.T) {

	// Output
	output, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()

	address := "Shout"
	output.SetIdentity(address)
	err = output.Bind("inproc://selftest-pingok")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Unbind("inproc://selftest-pingok")

	// Input
	input, err := zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()

	err = input.Connect("inproc://selftest-pingok")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Disconnect("inproc://selftest-pingok")

	// Create a Ping_Ok message and send it through the wire
	pingok := NewPingOk()
	pingok.SetSequence(123)

	err = pingok.Send(output)
	if err != nil {
		t.Fatal(err)
	}
	transit, err := Recv(input)
	if err != nil {
		t.Fatal(err)
	}

	tr := transit.(*PingOk)
	if tr.Sequence() != 123 {
		t.Fatalf("expected %d, got %d", 123, tr.Sequence())
	}

	err = tr.Send(input)
	if err != nil {
		t.Fatal(err)
	}
	transit, err = Recv(output)
	if err != nil {
		t.Fatal(err)
	}
	if address != tr.Address() {
		t.Fatalf("expected %v, got %v", address, tr.Address())
	}
}
