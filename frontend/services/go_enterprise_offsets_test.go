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
	wantMinors := []string{"1.2", "1.1", "1.0"}
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

func TestCompareGoEnterpriseOffsets_newVersionsAndModule(t *testing.T) {
	current := mustParseOffsetsFile(t, `{"timestamp":"2026-07-29T00:16:13.51429777Z","mods":[{"module":"example.com/mod","packages":[{"package":"example.com/mod/pkg","structs":[{"struct":"Foo","fields":[{"field":"Bar","offsets":[{"offset":8,"versions":["1.0.0","1.1.0"]}]}]}]}]}]}`)
	proposed := mustParseOffsetsFile(t, `{"timestamp":"2026-07-30T00:00:00Z","mods":[{"module":"example.com/mod","packages":[{"package":"example.com/mod/pkg","structs":[{"struct":"Foo","fields":[{"field":"Bar","offsets":[{"offset":8,"versions":["1.0.0","1.1.0","1.2.0"]}]}]}]}]},{"module":"example.com/new","packages":[{"package":"example.com/new/pkg","structs":[{"struct":"Baz","fields":[{"field":"Q","offsets":[{"offset":0,"versions":["2.0.0"]}]}]}]}]}]}`)

	result := compareGoEnterpriseOffsets(current, proposed)
	if !result.HasUpdates {
		t.Fatalf("expected hasUpdates")
	}
	if len(result.Mods) != 2 {
		t.Fatalf("unexpected mods: %+v", result.Mods)
	}
	if result.Mods[0].Module != "example.com/mod" || result.Mods[0].IsNew {
		t.Fatalf("unexpected first module: %+v", result.Mods[0])
	}
	foundNew := false
	for _, minor := range result.Mods[0].MinorVersions {
		for _, ver := range minor.Versions {
			if ver.Version == "1.2.0" {
				foundNew = true
				if !ver.IsNew {
					t.Fatalf("1.2.0 should be marked new")
				}
			}
			if ver.Version == "1.0.0" && ver.IsNew {
				t.Fatalf("1.0.0 should not be new")
			}
		}
	}
	if !foundNew {
		t.Fatalf("expected 1.2.0 in first module")
	}
	if result.Mods[1].Module != "example.com/new" || !result.Mods[1].IsNew {
		t.Fatalf("unexpected second module: %+v", result.Mods[1])
	}
}

func TestCompareGoEnterpriseOffsets_noUpdates(t *testing.T) {
	current := mustParseOffsetsFile(t, `{"timestamp":"2026-07-29T00:16:13.51429777Z","mods":[{"module":"example.com/mod","packages":[{"package":"example.com/mod/pkg","structs":[{"struct":"Foo","fields":[{"field":"Bar","offsets":[{"offset":8,"versions":["1.0.0"]}]}]}]}]}]}`)
	proposed := mustParseOffsetsFile(t, `{"timestamp":"2026-07-30T00:00:00Z","mods":[{"module":"example.com/mod","packages":[{"package":"example.com/mod/pkg","structs":[{"struct":"Foo","fields":[{"field":"Bar","offsets":[{"offset":8,"versions":["1.0.0"]}]}]}]}]}]}`)

	result := compareGoEnterpriseOffsets(current, proposed)
	if result.HasUpdates {
		t.Fatalf("expected no updates")
	}
	if len(result.Mods) != 1 || result.Mods[0].IsNew {
		t.Fatalf("unexpected mods: %+v", result.Mods)
	}
	if result.Mods[0].MinorVersions[0].Versions[0].IsNew {
		t.Fatalf("existing version should not be new")
	}
}

func mustParseOffsetsFile(t *testing.T, content string) *versionedModules {
	t.Helper()
	parsed, err := parseGoEnterpriseOffsetsFile(content)
	if err != nil {
		t.Fatalf("parseGoEnterpriseOffsetsFile: %v", err)
	}
	return parsed
}
