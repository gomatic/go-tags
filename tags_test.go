package tags_test

import (
	"errors"
	"testing"

	"github.com/gomatic/go-tags"
)

// pair is a minimal Pair implementation standing in for any generated message
// exposing GetKey/GetValue.
type pair struct{ key, value string }

func (p pair) GetKey() string   { return p.key }
func (p pair) GetValue() string { return p.value }

func TestToMap(t *testing.T) {
	got := tags.ToMap([]pair{{"a", "1"}, {"b", "2"}})
	if len(got) != 2 || got["a"] != "1" || got["b"] != "2" {
		t.Fatalf("ToMap = %v, want {a:1 b:2}", got)
	}
}

func TestToMapDuplicateLaterWins(t *testing.T) {
	got := tags.ToMap([]pair{{"k", "first"}, {"k", "last"}})
	if got["k"] != "last" {
		t.Fatalf("ToMap duplicate = %q, want %q", got["k"], "last")
	}
}

func TestToMapEmptyIsNonNil(t *testing.T) {
	got := tags.ToMap[pair](nil)
	if got == nil {
		t.Fatal("ToMap(nil) = nil, want non-nil empty map")
	}
	if len(got) != 0 {
		t.Fatalf("ToMap(nil) = %v, want empty", got)
	}
}

func TestFromMapSorted(t *testing.T) {
	got := tags.FromMap(map[string]string{"b": "2", "a": "1", "c": "3"}, func(key, value string) pair {
		return pair{key, value}
	})
	want := []string{"a", "b", "c"}
	if len(got) != 3 {
		t.Fatalf("FromMap len = %d, want 3", len(got))
	}
	for i, k := range want {
		if got[i].GetKey() != k {
			t.Fatalf("FromMap[%d].Key = %q, want %q (not sorted)", i, got[i].GetKey(), k)
		}
	}
	if got[0].GetValue() != "1" {
		t.Fatalf("FromMap[0].Value = %q, want %q", got[0].GetValue(), "1")
	}
}

func TestFromMapEmptyIsNonNil(t *testing.T) {
	got := tags.FromMap(nil, func(key, value string) pair { return pair{key, value} })
	if got == nil {
		t.Fatal("FromMap(nil) = nil, want non-nil empty slice")
	}
	if len(got) != 0 {
		t.Fatalf("FromMap(nil) = %v, want empty", got)
	}
}

func TestMerge(t *testing.T) {
	base := map[string]string{"a": "1", "b": "2"}
	got := tags.Merge(base, []pair{{"b", "two"}, {"c", "3"}})
	if got["a"] != "1" || got["b"] != "two" || got["c"] != "3" {
		t.Fatalf("Merge = %v, want {a:1 b:two c:3} (updates win)", got)
	}
	if base["b"] != "2" {
		t.Fatal("Merge mutated base")
	}
}

func TestMergeNilBase(t *testing.T) {
	got := tags.Merge(nil, []pair{{"a", "1"}})
	if got == nil || got["a"] != "1" {
		t.Fatalf("Merge(nil,...) = %v, want non-nil {a:1}", got)
	}
}

func TestMergeNilBaseNoUpdates(t *testing.T) {
	got := tags.Merge[pair](nil, nil)
	if got == nil {
		t.Fatal("Merge(nil,nil) = nil, want non-nil empty map (contract is always non-nil)")
	}
	if len(got) != 0 {
		t.Fatalf("Merge(nil,nil) = %v, want empty", got)
	}
}

func TestParse(t *testing.T) {
	got, err := tags.Parse([]string{"a=1", "b=2"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got["a"] != "1" || got["b"] != "2" {
		t.Fatalf("Parse = %v, want {a:1 b:2}", got)
	}
}

func TestParseValueWithEquals(t *testing.T) {
	got, err := tags.Parse([]string{"url=https://x?a=b"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got["url"] != "https://x?a=b" {
		t.Fatalf("Parse value = %q, want only the first = split", got["url"])
	}
}

func TestParseEmptyValue(t *testing.T) {
	got, err := tags.Parse([]string{"k="})
	if err != nil {
		t.Fatalf("Parse(%q) returned error: %v", "k=", err)
	}
	if v, ok := got["k"]; !ok || v != "" {
		t.Fatalf("Parse(%q) = %v, want key %q present with empty value", "k=", got, "k")
	}
}

func TestParseDuplicateLaterWins(t *testing.T) {
	got, err := tags.Parse([]string{"k=first", "k=last"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got["k"] != "last" {
		t.Fatalf("Parse duplicate = %q, want %q", got["k"], "last")
	}
}

func TestParseEmptyIsNonNil(t *testing.T) {
	got, err := tags.Parse(nil)
	if err != nil {
		t.Fatalf("Parse(nil) returned error: %v", err)
	}
	if got == nil {
		t.Fatal("Parse(nil) = nil, want non-nil empty map")
	}
}

func TestParseMissingEquals(t *testing.T) {
	got, err := tags.Parse([]string{"novalue"})
	if !errors.Is(err, tags.ErrInvalidPair) {
		t.Fatalf("Parse error = %v, want ErrInvalidPair", err)
	}
	if got != nil {
		t.Fatalf("Parse = %v, want nil on error", got)
	}
}

func TestParseEmptyKey(t *testing.T) {
	got, err := tags.Parse([]string{"=value"})
	if !errors.Is(err, tags.ErrInvalidPair) {
		t.Fatalf("Parse empty-key error = %v, want ErrInvalidPair", err)
	}
	if got != nil {
		t.Fatalf("Parse = %v, want nil on error", got)
	}
}
