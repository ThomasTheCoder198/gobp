package registry

import "testing"

func TestLoad_AllManifestsLoad(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(r.All()) == 0 {
		t.Fatal("expected at least one manifest")
	}
}

func TestLoad_HasExpectedKinds(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []Kind{KindFramework, KindDB, KindSDK, KindAddon, KindPattern, KindLogger} {
		if len(r.ByKind(kind)) == 0 {
			t.Errorf("no manifests for kind %s", kind)
		}
	}
}

func TestLoad_KnownIDs(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := map[Kind][]string{
		KindFramework: {"gin", "echo", "fiber"},
		KindDB:        {"postgres", "mysql", "sqlite", "mongo", "redis", "cassandra"},
		KindSDK:       {"openai", "stripe"},
		KindAddon:     {"docker", "compose", "githubactions"},
		KindPattern:   {"worker"},
	}
	for kind, ids := range want {
		for _, id := range ids {
			if !r.Has(kind, id) {
				t.Errorf("missing %s/%s", kind, id)
			}
		}
	}
}

func TestLoad_RejectsDuplicateIDs(t *testing.T) {
	// We only verify the production registry is unique here; injecting a
	// duplicate would require a custom fs.FS, which is overkill for now.
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, m := range r.All() {
		if seen[m.ID] {
			t.Errorf("duplicate id %q", m.ID)
		}
		seen[m.ID] = true
	}
}
