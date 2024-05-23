package reexp

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(0)
}

func TestPGMarker(t *testing.T) {
	isPg = true
	if pgMarker("..WHERE id=?") != "..WHERE id=$1" {
		t.Fail()
	}
	if pgMarker("..WHERE id=? AND fk=? ORDER..") != "..WHERE id=$1 AND fk=$2 ORDER.." {
		t.Fail()
	}
}

func TestComponent(t *testing.T) {
	Comp := []any{Client{}, Fund{}, Portfolio{}, CFP{}}
	Val := []string{
		`{"id":109,"dob":"2000-01-01T00:00:00Z","name":"Mr. John Doe","ni":"AA123456A"}`,
		`{"id":1,"name":"General","sector":"Technology","type":"Sustainable"}`,
		`{"id":307,"amount":25000,"name":"My Investment","opened_at":"2024-05-21T00:00:00Z","state":3}`,
		`{"client":109,"fund":1,"portfolio":307}`,
	}

	var c Client
	var p Portfolio
	var f Fund
	var cfp CFP

	var bfr bytes.Buffer
	for i := range 4 {
		dec := json.NewDecoder(strings.NewReader(Val[i]))
		enc := json.NewEncoder(&bfr)

		var err error
		switch Comp[i].(type) {
		case Client:
			err = dec.Decode(&c)
			enc.Encode(&c)
		case Fund:
			err = dec.Decode(&f)
			enc.Encode(&f)
		case Portfolio:
			err = dec.Decode(&p)
			enc.Encode(&p)
		case CFP:
			err = dec.Decode(&cfp)
			enc.Encode(&cfp)
		}
		if err != nil {
			t.Fatalf("Can't decode %T: %v", Comp[i], err)
		}
		log.Print(bfr.String())
		if Val[i] != strings.Trim(bfr.String(), "\r\n") {
			t.Fatalf("json In:Out don't match for: %T", Comp[i])
		}
		bfr.Reset()
	}
}
