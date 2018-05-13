package tuid

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

var (
	id      = TUID(48080669904283764)
	idHex   = "00aad11d52352874"
	idBytes = []byte{0x00, 0xaa, 0xd1, 0x1d, 0x52, 0x35, 0x28, 0x74}
	idTime  = time.Date(2018, 5, 13, 16, 15, 25, 0, time.UTC)

	model    = &struct{ ID TUID }{id}
	jsonWant = []byte(`{"ID":"00aad11d52352874"}`)
)

// func TestBlah(t *testing.T) {
// 	n := New()
// 	t.Log("TUID:", n)
// 	t.Log("TINT:", int64(n))
// 	t.Log("Time:", n.Time())
// 	t.Logf("Byte: %#v", n.Bytes())
// 	t.Error()
// }

func TestNew(t *testing.T) {
	New()
}

func TestFromBytes(t *testing.T) {
	v, err := FromBytes(idBytes)
	if err != nil {
		t.Error(err)
	}
	if v != id {
		t.Error("FromBytes(...) not equals id")
	}
}

func TestFromString(t *testing.T) {
	v, err := FromString(idHex)
	if err != nil {
		t.Error(err)
	}
	if v != id {
		t.Error("FromString(...) not equals id")
	}
}

func TestBytes(t *testing.T) {
	if !bytes.Equal(idBytes, id.Bytes()) {
		t.Error("id.Bytes() is wrong")
	}
}

func TestConvert(t *testing.T) {
	a := New()
	if b, err := FromString(a.String()); err != nil {
		t.Error(err)
	} else if a != b {
		t.Error("a != b")
	}
}

func TestJsonMarshal(t *testing.T) {
	has, err := json.Marshal(model)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(has, jsonWant) {
		t.Error("JSON model doesn't match")
	}
}

func TestJsonUnmarshal(t *testing.T) {
	var has struct{ ID TUID }
	if err := json.Unmarshal(jsonWant, &has); err != nil {
		t.Error(err)
	}
	if has.ID != model.ID {
		t.Error("JSON model doesn't match")
	}
}

func TestScan(t *testing.T) {
	cases := []struct {
		In   interface{}
		Want TUID
	}{
		{[]byte{0x00, 0xaa, 0xd5, 0x7f, 0x5d, 0x8b, 0x28, 0xd4}, TUID(48085489047775444)},
		{[]byte("00aad5cdf0b7c886"), TUID(48085826524399750)},
		{"00aad5fa5796b942", TUID(48086017228847426)},
		{int64(8634732875096), TUID(8634732875096)},
	}

	for _, c := range cases {
		var v TUID
		if err := v.Scan(c.In); err != nil {
			t.Error(err)
		}
		if v != c.Want {
			t.Errorf("Invalid value %v for input %v", v, c.In)
		}
	}
}

func TestValue(t *testing.T) {
	v, _ := id.Value()
	if _, ok := v.(int64); !ok {
		t.Error("Value() doen't return an int64")
	}
}

func TestTime(t *testing.T) {
	if !id.Time().Equal(idTime) {
		t.Error("Time doesn't match")
	}
}
