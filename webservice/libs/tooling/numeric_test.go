package tooling_test

import (
	"testing"
	"time"

	n "webservice/libs/tooling"
)

func TestRoundFloatFast(t *testing.T) {
	a := time.Now().UnixMilli()
	for i := 0.; i < 1_000_000; i += 1.23456789 {
		n.RoundFloatFast(i, 2)
	}
	t.Logf("RoundFloatFast millis=%d", time.Now().UnixMilli()-a)
}

func TestRoundFloat(t *testing.T) {
	f, _ := n.RoundFloat(0.000000000000000000000000000000011, 32)
	expected := 0.00000000000000000000000000000001
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.RoundFloat(92234, 14)
	expected = 92234.
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.RoundFloat(2.289, 2)
	expected = 2.29
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.RoundFloat(111222333444555666777888999000111222333.123456789, 8)
	expected = 111222333444555666777888999000111222333.12345679
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	a := time.Now().UnixMilli()
	for i := 0.; i < 1_000_000; i += 1.23456789 {
		n.RoundFloat(i, 2)
	}
	t.Logf("RoundFloat millis=%d", time.Now().UnixMilli()-a)
}

func TestTruncFloat(t *testing.T) {
	f, _ := n.TruncFloat(0.000000000000000000000000000000019, 32)
	expected := 0.00000000000000000000000000000001
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.TruncFloat(2.289, 2)
	expected = 2.28
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.TruncFloat(111222333444555666777888999000111222333.123456789, 8)
	expected = 111222333444555666777888999000111222333.12345678
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.TruncFloat(1, 0)
	expected = 1.
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	f, _ = n.TruncFloat(1, 2)
	expected = 1.
	if f != expected {
		t.Errorf("Expected=%+v Actual=%+v", expected, f)
	}

	a := time.Now().UnixMilli()
	for i := 0.; i < 1_000_000; i += 1.23456789 {
		n.TruncFloat(i, 2)
	}
	t.Logf("TruncFloat millis=%d", time.Now().UnixMilli()-a)
}
