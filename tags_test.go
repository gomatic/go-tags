package tags_test

import (
	"errors"
	"testing"

	typev1 "github.com/skykernel/api/src/proto/skykernel/type/v1"

	"github.com/skykernel/go-tags"
)

func TestToMap(t *testing.T) {
	got := tags.ToMap([]*typev1.Tag{{Key: "a", Value: "1"}, {Key: "b", Value: "2"}})
	if len(got) != 2 || got["a"] != "1" || got["b"] != "2" {
		t.Fatalf("ToMap = %v, want {a:1 b:2}", got)
	}
}

func TestToMapDuplicateLaterWins(t *testing.T) {
	got := tags.ToMap([]*typev1.Tag{{Key: "k", Value: "first"}, {Key: "k", Value: "last"}})
	if got["k"] != "last" {
		t.Fatalf("ToMap duplicate = %q, want %q", got["k"], "last")
	}
}

func TestToMapEmptyIsNonNil(t *testing.T) {
	got := tags.ToMap(nil)
	if got == nil {
		t.Fatal("ToMap(nil) = nil, want non-nil empty map")
	}
	if len(got) != 0 {
		t.Fatalf("ToMap(nil) = %v, want empty", got)
	}
}

func TestToProtoSorted(t *testing.T) {
	got := tags.ToProto(map[string]string{"b": "2", "a": "1", "c": "3"})
	want := []string{"a", "b", "c"}
	if len(got) != 3 {
		t.Fatalf("ToProto len = %d, want 3", len(got))
	}
	for i, k := range want {
		if got[i].GetKey() != k {
			t.Fatalf("ToProto[%d].Key = %q, want %q (not sorted)", i, got[i].GetKey(), k)
		}
	}
	if got[0].GetValue() != "1" {
		t.Fatalf("ToProto[0].Value = %q, want %q", got[0].GetValue(), "1")
	}
}

func TestToProtoEmptyIsNonNil(t *testing.T) {
	got := tags.ToProto(nil)
	if got == nil {
		t.Fatal("ToProto(nil) = nil, want non-nil empty slice")
	}
	if len(got) != 0 {
		t.Fatalf("ToProto(nil) = %v, want empty", got)
	}
}

func TestMerge(t *testing.T) {
	base := map[string]string{"a": "1", "b": "2"}
	got := tags.Merge(base, []*typev1.Tag{{Key: "b", Value: "two"}, {Key: "c", Value: "3"}})
	if got["a"] != "1" || got["b"] != "two" || got["c"] != "3" {
		t.Fatalf("Merge = %v, want {a:1 b:two c:3} (updates win)", got)
	}
	if base["b"] != "2" {
		t.Fatal("Merge mutated base")
	}
}

func TestMergeNilBase(t *testing.T) {
	got := tags.Merge(nil, []*typev1.Tag{{Key: "a", Value: "1"}})
	if got == nil || got["a"] != "1" {
		t.Fatalf("Merge(nil,...) = %v, want non-nil {a:1}", got)
	}
}

func TestMergeNilBaseNoUpdates(t *testing.T) {
	got := tags.Merge(nil, nil)
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

func TestErrorString(t *testing.T) {
	if tags.ErrInvalidPair.Error() != string(tags.ErrInvalidPair) {
		t.Fatal("Error() does not return the underlying string")
	}
}
