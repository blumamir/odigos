package services

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseGoEnterpriseOffsetsContent_jsonString(t *testing.T) {
	inner := `{"timestamp":"2026-07-29T00:16:13.51429777Z","mods":[{"module":"example.com/mod","packages":[{"package":"example.com/mod/pkg","structs":[{"struct":"Foo","fields":[{"field":"Bar","offsets":[{"offset":8,"versions":["1.0.0","1.1.0"]},{"offset":16,"versions":["1.1.0","1.2.0"]}]}]}]}]}]}`
	wrapped, err := json.Marshal(inner)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	parsed, err := parseGoEnterpriseOffsetsContent(string(wrapped))
	if err != nil {
		t.Fatalf("parseGoEnterpriseOffsetsContent: %v", err)
	}
	model := goEnterpriseOffsetsToModel(parsed)
	if len(model.Mods) != 1 || model.Mods[0].Module != "example.com/mod" {
		t.Fatalf("unexpected module: %+v", model.Mods)
	}
	if model.Mods[0].MinVersion != "1.0.0" {
		t.Fatalf("unexpected minVersion: %q", model.Mods[0].MinVersion)
	}
	if model.Mods[0].MaxVersion != "1.2.0" {
		t.Fatalf("unexpected maxVersion: %q", model.Mods[0].MaxVersion)
	}
	wantMinors := []string{"1.0", "1.1", "1.2"}
	if len(model.Mods[0].MinorVersions) != len(wantMinors) {
		t.Fatalf("unexpected minorVersions: %+v", model.Mods[0].MinorVersions)
	}
	for i, want := range wantMinors {
		if model.Mods[0].MinorVersions[i].MinorVersion != want {
			t.Fatalf("unexpected minorVersions[%d]: got %q want %q", i, model.Mods[0].MinorVersions[i].MinorVersion, want)
		}
	}
	if parsed.Timestamp.UTC() != time.Date(2026, 7, 29, 0, 16, 13, 514297770, time.UTC) {
		t.Fatalf("unexpected timestamp: %v", parsed.Timestamp)
	}
}

func TestParseGoEnterpriseOffsetsContent_jsonStringWithSignature(t *testing.T) {
	inner := `{"timestamp":"2026-07-29T00:16:13.51429777Z","mods":[]}---SIGNATURE---abc`
	wrapped, err := json.Marshal(inner)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	parsed, err := parseGoEnterpriseOffsetsContent(string(wrapped))
	if err != nil {
		t.Fatalf("parseGoEnterpriseOffsetsContent: %v", err)
	}
	if len(parsed.Mods) != 0 {
		t.Fatalf("expected empty mods")
	}
}

func TestParseGoEnterpriseOffsetsContent_empty(t *testing.T) {
	parsed, err := parseGoEnterpriseOffsetsContent("  ")
	if err != nil {
		t.Fatalf("parseGoEnterpriseOffsetsContent: %v", err)
	}
	if len(parsed.Mods) != 0 {
		t.Fatalf("expected empty mods")
	}
}
