package memory

import (
	"context"
	"testing"
)

// seedLayeredStore creates a session + L2 scenario + L3 persona-delta personas
// records (the step-02 bootstrap sources) and an L1 atom observation (the RRF
// fallback source), each provenance-linked. It returns the session id.
func seedLayeredStore(t *testing.T, fx *fixture) {
	t.Helper()
	ctx := context.Background()
	// L0 source: the fixture's session1 (a resolvable live session row).
	l0 := session1
	// L1 atom observation (the RRF fallback source — reachable via FTS).
	if _, err := fx.svc.Save(ctx, SaveInput{
		SessionID: l0, Type: "architecture",
		Title:   "Persona delta marker PERSONA_DELTA_XYZ",
		Content: "Persona profile increment PERSONA_DELTA_XYZ",
	}); err != nil {
		t.Fatalf("save persona atom: %v", err)
	}
	// L2 scenario + L3 persona-delta personas records (step-02 bootstrap).
	if _, err := fx.st.DB.Exec(`
		INSERT INTO personas (project, kind, title, content, content_hash, created_at)
		VALUES (?, 'persona_delta', 'Persona delta marker PERSONA_DELTA_XYZ', 'Persona profile increment PERSONA_DELTA_XYZ', 'h-delta', '2026-01-02T00:00:00Z')`,
		fx.svc.ProjectID()); err != nil {
		t.Fatalf("insert persona_delta: %v", err)
	}
	if _, err := fx.st.DB.Exec(`
		INSERT INTO personas (project, kind, title, content, content_hash, created_at)
		VALUES (?, 'scenario', 'Scenario marker SCENARIO_ABC', 'Working context scenario SCENARIO_ABC', 'h-scen', '2026-01-02T00:00:00Z')`,
		fx.svc.ProjectID()); err != nil {
		t.Fatalf("insert scenario: %v", err)
	}
	deltaID := lastPersonaID(t, fx, "persona_delta")
	scenID := lastPersonaID(t, fx, "scenario")
	if _, err := fx.st.DB.Exec(`
		INSERT INTO observation_layers (project, layer, target_kind, target_id, source_session, source_topic, content_hash, created_at)
		VALUES (?, 'L3', 'persona', ?, ?, 'delta-topic', 'h-delta', '2026-01-02T00:00:00Z')`,
		fx.svc.ProjectID(), deltaID, l0); err != nil {
		t.Fatalf("link L3: %v", err)
	}
	if _, err := fx.st.DB.Exec(`
		INSERT INTO observation_layers (project, layer, target_kind, target_id, source_session, source_topic, content_hash, created_at)
		VALUES (?, 'L2', 'persona', ?, ?, 'scen-topic', 'h-scen', '2026-01-02T00:00:00Z')`,
		fx.svc.ProjectID(), scenID, l0); err != nil {
		t.Fatalf("link L2: %v", err)
	}
}

func lastPersonaID(t *testing.T, fx *fixture, kind string) int64 {
	t.Helper()
	var id int64
	if err := fx.st.DB.QueryRow(`
		SELECT id FROM personas WHERE project = ? AND kind = ? ORDER BY id DESC LIMIT 1`,
		fx.svc.ProjectID(), kind).Scan(&id); err != nil {
		t.Fatalf("last persona %s: %v", kind, err)
	}
	return id
}

// TestRetrieve is 03.2 [RED] — the "Mnemonic tool surface" threat. It proves:
// (a) layered retrieval returns L2/L3 first (bootstrap) and falls back to L1/L0
// via the existing RRF path for a specific fact; (b) every in-list hit carries
// its full-content fetch id (mem_get_observation is the only full-content
// path); and (c) bad retrieval args are rejected clearly.
//
// Scenarios: layered-retrieval-l2-l3-first-with-l1-l0-rrf-fallback,
// mem-get-observation-is-only-full-content-path.
func TestRetrieve(t *testing.T) {
	fx := newFixture(t, "retrieve-proj")
	ctx := context.Background()
	seedLayeredStore(t, fx)

	t.Run("bootstrap-returns-l2-l3-first", func(t *testing.T) {
		hits, err := fx.svc.Retrieve(ctx, RetrieveOpts{Mode: "bootstrap", Limit: 10})
		if err != nil {
			t.Fatalf("retrieve bootstrap: %v", err)
		}
		if len(hits) < 2 {
			t.Fatalf("expected L2+L3 bootstrap hits, got %d", len(hits))
		}
		// L3 (persona_delta) must come before L2 (scenario): the most stable
		// layer first.
		if hits[0].Layer != "L3" {
			t.Fatalf("expected L3 first, got %q (hit %+v)", hits[0].Layer, hits[0])
		}
		// Every in-list hit carries its full-content fetch id.
		for _, h := range hits {
			if h.GetObservationID == 0 {
				t.Fatalf("in-list hit missing full-content fetch id: %+v", h)
			}
		}
	})

	t.Run("fact-falls-back-to-l1-l0-rrf", func(t *testing.T) {
		// A specific-fact query that matches the L1 atom content falls back to
		// the existing RRF path (BlendedSearch / FTS) and returns L1 hits.
		hits, err := fx.svc.Retrieve(ctx, RetrieveOpts{Mode: "fact", Query: "persona delta marker", Limit: 10})
		if err != nil {
			t.Fatalf("retrieve fact: %v", err)
		}
		if len(hits) == 0 {
			t.Fatal("expected an L1/L0 RRF fallback hit for the specific fact")
		}
		for _, h := range hits {
			if h.Layer != "L1" {
				t.Fatalf("RRF fallback should return L1 hits, got layer %q", h.Layer)
			}
			if h.GetObservationID == 0 {
				t.Fatalf("RRF fallback hit missing full-content fetch id: %+v", h)
			}
		}
	})

	t.Run("bad-mode-rejected", func(t *testing.T) {
		_, err := fx.svc.Retrieve(ctx, RetrieveOpts{Mode: "nonsense", Limit: 5})
		if err == nil {
			t.Fatal("unknown retrieve mode should be rejected")
		}
	})
}
