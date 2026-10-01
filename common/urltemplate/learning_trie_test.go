package urltemplate

import (
	"testing"
)

func TestBuildPathTrie(t *testing.T) {
	trie := BuildPathTrie(map[string]string{
		"/users/1":       "10",
		"/users/2":       "5",
		"/users/1/orders": "3",
		"/health":        "100",
	})

	users := trie.Root.Children["users"]
	if users == nil {
		t.Fatal("expected users node")
	}
	if users.Leaf {
		t.Fatal("users should not be a leaf")
	}
	if users.UniqueLeaves != 3 {
		t.Fatalf("users uniqueLeaves=%d, want 3", users.UniqueLeaves)
	}
	if users.TotalObservations != 18 {
		t.Fatalf("users totalObservations=%d, want 18", users.TotalObservations)
	}

	one := users.Children["1"]
	if one == nil || !one.Leaf {
		t.Fatal("expected /users/1 to be a leaf")
	}
	if one.ObservationCount != 10 {
		t.Fatalf("/users/1 obs=%d, want 10", one.ObservationCount)
	}
	if one.UniqueLeaves != 2 { // itself + /users/1/orders
		t.Fatalf("/users/1 uniqueLeaves=%d, want 2", one.UniqueLeaves)
	}
	if one.TotalObservations != 13 {
		t.Fatalf("/users/1 totalObservations=%d, want 13", one.TotalObservations)
	}

	orders := one.Children["orders"]
	if orders == nil || !orders.Leaf || orders.ObservationCount != 3 {
		t.Fatalf("expected /users/1/orders leaf with obs=3, got %+v", orders)
	}

	health := trie.Root.Children["health"]
	if health == nil || !health.Leaf || health.ObservationCount != 100 {
		t.Fatalf("expected /health leaf with obs=100, got %+v", health)
	}

	if trie.Root.UniqueLeaves != 4 {
		t.Fatalf("root uniqueLeaves=%d, want 4", trie.Root.UniqueLeaves)
	}
	if trie.Root.TotalObservations != 118 {
		t.Fatalf("root totalObservations=%d, want 118", trie.Root.TotalObservations)
	}
}
