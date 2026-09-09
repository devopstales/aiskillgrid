// Package service is the shared facade over memory, codeindex, webcache, and search.
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/codeindex"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/config"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/embedder"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/files"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/graph"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/hybrid"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/memory"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/project"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/search"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/store"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/webcache"
)

// Service is the shared facade over memory, codeindex, webcache, and search.
type Service struct {
	dataDir string
}

// New creates a service using dataDir for per-project SQLite stores.
func New(dataDir string) *Service {
	return &Service{dataDir: dataDir}
}

// DefaultDataDir returns the mnemonic data directory from env or ~/.skillgrid/mnemonic.
func DefaultDataDir() (string, error) {
	if v := os.Getenv("SKILLGRID_MNEMONIC_DATA_DIR"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".skillgrid", "mnemonic"), nil
}

type projectHandle struct {
	store     *store.Store
	projectID string
	root      string // workspace directory for ContentPlane (.skillgrid/files)
	memory    *memory.Service
	web       *webcache.Service
	content   *files.ContentPlane
}

func (s *Service) openProject(projectID, configRoot string) (*projectHandle, func(), error) {
	if s == nil {
		return nil, nil, fmt.Errorf("service not initialized")
	}
	st, err := store.Open(s.dataDir, projectID)
	if err != nil {
		return nil, nil, err
	}
	root := configRoot
	if abs, absErr := filepath.Abs(configRoot); absErr == nil {
		root = abs
	}
	cfg := config.Load(root)
	h := &projectHandle{
		store:     st,
		projectID: projectID,
		root:      root,
		memory:    memory.New(st, projectID),
		web:       webcache.New(st, projectID, cfg.WebCache),
		content:   files.NewContentPlane(root),
	}
	return h, func() { st.Close() }, nil
}

func (s *Service) openProjectForDirectory(directory string) (*projectHandle, func(), error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve directory: %w", err)
	}
	res, resErr := project.ResolveDetailed(absDir)
	if resErr != nil {
		// Hard abort for writes: never open/create under an ambiguous
		// directory-hash fallback (or other resolve failures such as binding
		// write errors). Recover via MNEMONIC_PROJECT or explicit project=.
		return nil, nil, fmt.Errorf("resolve project: %w", resErr)
	}
	// Best-effort: fold any pre-identity directory-hash store for this path
	// into the canonical identity bucket so prior memories are reachable and
	// future alias-named writes route here. Idempotent and read-mostly.
	if res.SeedID != "" && res.SeedID != res.ID && res.Source == project.SourceIdentity {
		if _, _, err := s.MergeProjects(context.Background(), res.SeedID, res.ID); err == nil {
			// recorded
		}
	}
	return s.openProject(res.ID, absDir)
}

func (s *Service) openProjectFromCWD() (*projectHandle, func(), error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	return s.openProjectForDirectory(cwd)
}

// ListProjects returns the project IDs with a store in dataDir, sorted.
func (s *Service) ListProjects() ([]string, error) {
	entries, err := os.ReadDir(s.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sqlite") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".sqlite")
		if strings.TrimSpace(id) != "" {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids, nil
}

// ObservationsRecent returns stored observations, newest first.
func (s *Service) ObservationsRecent(ctx context.Context, projectID string, limit int) ([]memory.Observation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.memory.Recent(ctx, limit)
}

// ResolveProject returns the project ID for directory.
func (s *Service) ResolveProject(directory string) (string, error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	return project.Resolve(absDir)
}

// SessionStart creates a workspace session for directory. title is the
// session name shown in the web dashboard (mem-sessions).
func (s *Service) SessionStart(ctx context.Context, directory, title string) (sessionID, projectID string, err error) {
	h, cleanup, err := s.openProjectForDirectory(directory)
	if err != nil {
		return "", "", err
	}
	defer cleanup()
	sessionID, err = h.memory.SessionStart(ctx, directory, title)
	if err != nil {
		return "", "", err
	}
	return sessionID, h.projectID, nil
}

// SessionSetTitle renames a session so the dashboard session list shows it.
func (s *Service) SessionSetTitle(ctx context.Context, projectID, sessionID, title string) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.SessionSetTitle(ctx, sessionID, title)
}

// SessionEnd ends a session with optional summary.
func (s *Service) SessionEnd(ctx context.Context, projectID, sessionID, summary string) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.SessionEnd(ctx, sessionID, summary)
}

// SessionSummary stores an end-of-session summary.
func (s *Service) SessionSummary(ctx context.Context, projectID, sessionID, summary string) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.SessionSummary(ctx, sessionID, summary)
}

// SessionStartByClientID registers a session under the caller's ID (idempotent).
func (s *Service) SessionStartByClientID(ctx context.Context, sessionID, directory, title string) (id, projectID string, existed bool, err error) {
	h, cleanup, err := s.openProjectForDirectory(directory)
	if err != nil {
		return "", "", false, err
	}
	defer cleanup()
	id, projID, existed, err := h.memory.SessionStartByClientID(ctx, sessionID, directory, title)
	if err != nil {
		return "", "", false, err
	}
	return id, projID, existed, nil
}

// PromptInput is a captured user prompt.
type PromptInput = memory.PromptInput

// SavePrompt stores a captured user prompt.
func (s *Service) SavePrompt(ctx context.Context, projectID string, in PromptInput) (int64, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return 0, err
	}
	defer cleanup()
	return h.memory.SavePrompt(ctx, in)
}

// PassiveInput is free text for server-side learnings extraction.
type PassiveInput = memory.PassiveInput

// PassiveResult reports what the passive extractor found.
type PassiveResult = memory.CapturePassiveResult

// CapturePassive extracts learnings from raw text and persists them.
func (s *Service) CapturePassive(ctx context.Context, projectID string, in PassiveInput) (PassiveResult, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return PassiveResult{}, err
	}
	defer cleanup()
	return h.memory.CapturePassive(ctx, in)
}

// CompactionContext assembles session-scoped context for the compaction prompt.
type CompactionContext = memory.CompactionContext

// ContextForCompaction returns the compaction context.
func (s *Service) ContextForCompaction(ctx context.Context, projectID, sessionID string, limit int) (CompactionContext, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return CompactionContext{}, err
	}
	defer cleanup()
	return h.memory.CompactionContext(ctx, sessionID, limit)
}

// MigrateProjects rolls data recorded under oldProject into the newProject
// store. In Mnemonic each project has its own SQLite file, so the data lives
// in <oldProject>.sqlite while new writes go to <newProject>.sqlite. We open
// the old store, ATTACH the new one, copy every row tagged oldProject across
// (observations, sessions, web_cache, prompts, files), re-tag it as
// newProject, and record the migration for idempotency.
//
// Idempotent: if no rows tagged oldProject remain, nothing is copied and we
// just ensure the record exists. If the old store file is missing, this is a
// clean no-op (0 moved, nil error).
//
// Called from the agent plugin on first run when the project id it computed
// differs from one previously recorded in the data dir.
func (s *Service) MigrateProjects(ctx context.Context, oldProject, newProject string) (int, error) {
	oldProject = strings.TrimSpace(oldProject)
	newProject = strings.TrimSpace(newProject)
	if oldProject == "" || newProject == "" || oldProject == newProject {
		return 0, nil
	}

	oldPath := filepath.Join(s.dataDir, oldProject+".sqlite")
	if _, err := os.Stat(oldPath); err != nil {
		if os.IsNotExist(err) {
			return 0, nil // nothing to migrate
		}
		return 0, err
	}
	// Ensure the destination store exists (and is migrated) before we attach.
	newStore, err := store.Open(s.dataDir, newProject)
	if err != nil {
		return 0, fmt.Errorf("open destination store: %w", err)
	}
	defer newStore.Close()
	oldStore, err := store.Open(s.dataDir, oldProject)
	if err != nil {
		return 0, fmt.Errorf("open source store: %w", err)
	}
	defer oldStore.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	attachSQL := "ATTACH DATABASE ? AS newdb"
	res, err := oldStore.DB.ExecContext(ctx, attachSQL, newStore.Path())
	if err != nil {
		return 0, fmt.Errorf("attach new store: %w", err)
	}
	defer func() {
		_, _ = oldStore.DB.Exec("DETACH DATABASE newdb")
	}()
	_ = res

	t := oldStore.DB

	// Bulk copy across attached DBs must disable FK checks: sessions must land
	// before observations, and OR IGNORE / schema drift can otherwise leave
	// dangling session_id references mid-copy.
	if _, err := t.ExecContext(ctx, `PRAGMA foreign_keys=OFF`); err != nil {
		return 0, fmt.Errorf("disable foreign_keys: %w", err)
	}
	defer func() { _, _ = t.Exec(`PRAGMA foreign_keys=ON`) }()

	var total int
	copied := []string{}
	// sessions first: observations.session_id REFERENCES sessions(id).
	// sessions/prompts have no deleted_at column — do not filter on it.
	for _, table := range []string{"sessions", "prompts", "observations", "web_cache", "files"} {
		hasDeletedAt := table == "observations" || table == "web_cache" || table == "files"
		countSQL := `SELECT COUNT(*) FROM ` + table + ` WHERE project = ?`
		if hasDeletedAt {
			countSQL += ` AND deleted_at IS NULL`
		}
		var n int
		err := t.QueryRowContext(ctx, countSQL, oldProject).Scan(&n)
		if err != nil || n == 0 {
			continue
		}
		insertSQL := `INSERT OR IGNORE INTO newdb.` + table + ` SELECT * FROM ` + table + ` WHERE project = ?`
		if hasDeletedAt {
			insertSQL += ` AND (deleted_at IS NULL OR deleted_at = '')`
		}
		res, err := t.ExecContext(ctx, insertSQL, oldProject)
		if err != nil {
			if strings.Contains(err.Error(), "no such table") || strings.Contains(err.Error(), "duplicate column") || strings.Contains(err.Error(), "datatype mismatch") {
				continue
			}
			return total, fmt.Errorf("copy %s: %w", table, err)
		}
		affected, _ := res.RowsAffected()
		total += int(affected)
		copied = append(copied, fmt.Sprintf("%s=%d", table, affected))
		// Re-tag copied rows in the destination.
		if _, err := t.ExecContext(ctx, `
			UPDATE newdb.`+table+` SET project = ? WHERE project = ?`,
			newProject, oldProject); err != nil {
			if !strings.Contains(err.Error(), "no such table") {
				return total, fmt.Errorf("retag %s: %w", table, err)
			}
		}
	}

	// Record the migration in the destination for idempotency.
	if _, err := newStore.DB.ExecContext(ctx, `
		INSERT INTO project_migrations (old_project, new_project, migrated_at)
		VALUES (?, ?, ?)`,
		oldProject, newProject, now); err != nil && !strings.Contains(err.Error(), "no such table") {
		// ignore — best effort
	}

	return total, nil
}

// LastObservationAt returns the newest observation time for the project, if any.
func (s *Service) LastObservationAt(ctx context.Context, projectID string) (time.Time, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return time.Time{}, err
	}
	defer cleanup()
	return h.memory.LastObservationAt(ctx)
}

// RecordRelation stores a semantic link between two observations in the
// given project (mem_judge / mem_compare).
func (s *Service) RecordRelation(ctx context.Context, projectID string, srcID, dstID int64, relation, reason string, confidence *float64) (memory.Relation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return memory.Relation{}, err
	}
	defer cleanup()
	return h.memory.RecordRelation(ctx, srcID, dstID, relation, reason, confidence)
}

// RemoveRelation clears a live link between two observations (mem_judge
// not_conflict verdict).
func (s *Service) RemoveRelation(ctx context.Context, projectID string, srcID, dstID int64, relation string) (bool, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return false, err
	}
	defer cleanup()
	return h.memory.RemoveRelation(ctx, srcID, dstID, relation)
}

// RelationsOf returns every live relation touching observation id.
func (s *Service) RelationsOf(ctx context.Context, projectID string, id int64) ([]memory.Relation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.memory.RelationsOf(ctx, id)
}

// RelationsBetween returns the live links between two specific observations.
func (s *Service) RelationsBetween(ctx context.Context, projectID string, srcID, dstID int64) ([]memory.Relation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.memory.RelationsBetween(ctx, srcID, dstID)
}

// ProjectDrift mirrors the memory package's drift report, for API consumers.
type ProjectDrift = memory.ProjectDrift

// CheckProjectDrift returns a drift report when projectName is a known alias
// of a canonical project. It probes the store that most likely holds the
// alias row (the canonical project's store, then the current CWD store, then
// every known store) and does not modify state.
func (s *Service) CheckProjectDrift(ctx context.Context, projectName string) (*ProjectDrift, error) {
	name := strings.TrimSpace(projectName)
	if name == "" {
		return nil, nil
	}
	probe := func(storeID string) (*ProjectDrift, bool) {
		if storeID == "" {
			return nil, false
		}
		h, cleanup, err := s.openProject(storeID, ".")
		if err != nil {
			return nil, false
		}
		d, err := h.memory.CheckProjectDrift(ctx, name)
		cleanup()
		if err != nil {
			return nil, false
		}
		if d != nil {
			return d, true
		}
		return nil, false
	}
	// 1) The canonical project (if the alias points at one of our stores).
	if canonical := s.canonicalForAlias(ctx, name); canonical != "" {
		if d, ok := probe(canonical); ok {
			return d, nil
		}
	}
	// 2) The explicitly-named store (it may hold its own alias row).
	if n := storeIDFor(name); n != "" {
		if d, ok := probe(n); ok {
			return d, nil
		}
	}
	// 3) Current CWD store.
	if cwd, err := os.Getwd(); err == nil {
		if pid, err := s.ResolveProject(cwd); err == nil {
			if d, ok := probe(pid); ok {
				return d, nil
			}
		}
	}
	return nil, nil
}

// storeIDFor normalizes a project name into its store ID (best-effort — no IO).
func storeIDFor(name string) string {
	n := strings.ToLower(strings.TrimSpace(name))
	if n == "" {
		return ""
	}
	return n
}

// canonicalForAlias scans existing stores for a project_aliases row naming
// projectName as an alias, and returns the canonical name. Best-effort.
func (s *Service) canonicalForAlias(ctx context.Context, alias string) string {
	projects, err := s.ListProjects()
	if err != nil {
		return ""
	}
	for _, p := range projects {
		h, cleanup, err := s.openProject(p, ".")
		if err != nil {
			continue
		}
		var canonical string
		err = h.store.DB.QueryRowContext(ctx, `
			SELECT canonical FROM project_aliases WHERE alias = ? LIMIT 1`, alias,
		).Scan(&canonical)
		cleanup()
		if err == nil && canonical != "" {
			return canonical
		}
	}
	return ""
}

// MergeProjects consolidates data recorded under the source project name into
// the canonical project, then records a project alias. Returns the number of
// rows moved and the canonical name.
//
// Behaviour:
//   - If source == canonical or either is blank, no-op (0, canonical).
//   - If the source store file does not exist, no-op (0, canonical) — but the
//     alias is still recorded so future writes to the source name land in the
//     canonical store.
//   - Otherwise: copies observations/sessions/web_cache/prompts/files rows
//     tagged source into the canonical store (INSERT OR IGNORE so re-runs do
//     not fail on PK collisions), re-tags them as canonical, and records the
//     migration + alias (idempotent via unique primary key).
func (s *Service) MergeProjects(ctx context.Context, source, canonical string) (int, string, error) {
	source = strings.TrimSpace(source)
	canonical = strings.TrimSpace(canonical)
	if source == "" || canonical == "" || source == canonical {
		return 0, canonical, nil
	}

	// Record the alias first so that future resolution always lands in the
	// canonical store regardless of whether the copy below succeeds.
	if err := s.recordProjectAlias(ctx, canonical, source); err != nil {
		return 0, canonical, fmt.Errorf("record alias: %w", err)
	}

	// No data to move if the source store file is missing.
	srcPath := filepath.Join(s.dataDir, source+".sqlite")
	if _, err := os.Stat(srcPath); err != nil {
		if os.IsNotExist(err) {
			return 0, canonical, nil
		}
		return 0, canonical, err
	}

	moved, err := s.MigrateProjects(ctx, source, canonical)
	if err != nil {
		return 0, canonical, err
	}
	return moved, canonical, nil
}

// recordProjectAlias marks source as an alias of canonical in the canonical
// store. Idempotent: re-merge keeps the first merge time.
func (s *Service) recordProjectAlias(ctx context.Context, canonical, source string) error {
	newStore, err := store.Open(s.dataDir, canonical)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "no such table") {
			return nil // no canonical store yet — alias is implicit
		}
		return err
	}
	defer newStore.Close()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err = newStore.DB.ExecContext(ctx, `
		INSERT INTO project_aliases (alias, canonical, merged_at)
		VALUES (?, ?, ?)
		ON CONFLICT(alias) DO NOTHING`,
		source, canonical, now,
	)
	return err
}

// SessionStartedAt returns the started_at of a session (zero if missing).
func (s *Service) SessionStartedAt(ctx context.Context, projectID, sessionID string) (time.Time, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return time.Time{}, err
	}
	defer cleanup()
	return h.memory.SessionStartedAt(ctx, sessionID)
}

// ListProjectsForMigration returns all project ids with a store file.
func (s *Service) ListProjectsForMigration() ([]string, error) {
	return s.ListProjects()
}

// SaveObservationInput holds fields for saving an observation (HTTP + MCP).
type SaveObservationInput struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Scope     string `json:"scope"`
	TopicKey  string `json:"topic_key"`
	SessionID string `json:"session_id"`
	// CapturePrompt, when true, best-effort links the session's most recent
	// user prompt to the observation (mem_save capture_prompt).
	CapturePrompt bool `json:"capture_prompt"`
	// ProjectName, when non-empty, names the logical project for the save. It
	// is validated against project_aliases and a drift warning surfaced by the
	// caller when the name has been retired by mem_merge_projects.
	ProjectName string `json:"project_name"`
	// ToolName is optional provenance for which tool produced the save.
	ToolName string `json:"tool_name"`
}

// SaveObservation stores an observation with scope normalization matching MCP mem_save.
func (s *Service) SaveObservation(ctx context.Context, projectID string, in SaveObservationInput) (int64, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return 0, err
	}
	defer cleanup()
	scope := in.Scope
	if scope == "" {
		scope = "project"
	}
	if scope == "personal" {
		scope = "user"
	}
	return h.memory.Save(ctx, memory.SaveInput{
		Title:         in.Title,
		Type:          in.Type,
		Content:       in.Content,
		Scope:         scope,
		TopicKey:      in.TopicKey,
		SessionID:     in.SessionID,
		CapturePrompt: in.CapturePrompt,
		ProjectName:   in.ProjectName,
		ToolName:      in.ToolName,
	})
}

// SearchObservations runs FTS over observations (scope = any).
func (s *Service) SearchObservations(ctx context.Context, projectID, query, matchMode string, limit int) ([]memory.Observation, error) {
	return s.SearchObservationsScoped(ctx, projectID, query, matchMode, "", limit)
}

// SearchObservationsAll runs the same FTS query across every store in dataDir
// and returns the union, ordered by global rank (each store returns its own
// bm25-ranked list; the merged result interleaves by cross-store rank, so a
// #1 hit in one store is never buried under #5 hits from another). Used by
// mem_search all_projects=true so an agent at a parent directory can still
// find memories stored under a child project's store.
func (s *Service) SearchObservationsAll(ctx context.Context, query, matchMode, scope string, limit int) ([]memory.Observation, error) {
	if limit <= 0 {
		limit = 20
	}
	projects, err := s.ListProjects()
	if err != nil {
		return nil, err
	}
	type ranked struct {
		obs  memory.Observation
		rank int // 0-based rank within a single store
	}
	seen := map[string]bool{}
	var collected []ranked
	for _, pid := range projects {
		res, err := s.SearchObservationsScoped(ctx, pid, query, matchMode, scope, limit)
		if err != nil {
			continue
		}
		for i, o := range res {
			key := strconv.FormatInt(o.ID, 10) + "/" + o.Project
			if seen[key] {
				continue
			}
			seen[key] = true
			collected = append(collected, ranked{obs: o, rank: i})
		}
	}
	sort.SliceStable(collected, func(i, j int) bool {
		if collected[i].rank != collected[j].rank {
			return collected[i].rank < collected[j].rank
		}
		return collected[i].obs.UpdatedAt > collected[j].obs.UpdatedAt
	})
	out := make([]memory.Observation, 0, len(collected))
	for _, r := range collected {
		out = append(out, r.obs)
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// SearchObservationsScoped runs FTS over observations, restricting to a
// visibility scope when scope is non-empty (project|user|global).
func (s *Service) SearchObservationsScoped(ctx context.Context, projectID, query, matchMode, scope string, limit int) ([]memory.Observation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.memory.SearchWithScope(ctx, query, matchMode, scope, limit)
}

// GetObservation returns a single observation by ID.
func (s *Service) GetObservation(ctx context.Context, projectID string, id int64) (memory.Observation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return memory.Observation{}, err
	}
	defer cleanup()
	return h.memory.Get(ctx, id)
}

// RecentContext returns recent session summaries.
func (s *Service) RecentContext(ctx context.Context, projectID string, limit int) ([]memory.Session, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.memory.RecentContext(ctx, limit)
}

// SearchAllProjects runs the FTS query across every store in the data dir and
// returns the unified, rank-merged result set. This is the backing for
// mem_search(all_projects=true) and rescues memories stranded under a
// directory-hash store for the same logical project.
func (s *Service) SearchAllProjects(ctx context.Context, query, matchMode, scope string, limit int) ([]memory.Observation, error) {
	return s.SearchObservationsAll(ctx, query, matchMode, scope, limit)
}

// BlendedSearch runs the FTS leg and, when a non-empty query vector is
// supplied and embedding recall is enabled, the vector leg, merging the two
// ranked lists with reciprocal rank fusion (P4). Passing an empty vector
// returns the plain FTS result unchanged — FTS5 is the floor.
func (s *Service) BlendedSearch(ctx context.Context, projectID, query, matchMode, scope string, queryVec []float32, limit int) ([]memory.Observation, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	vec := memory.Vector{Data: queryVec}
	return h.memory.BlendedSearch(ctx, query, matchMode, scope, vec, limit)
}

// SetObservationEmbedding stores a precomputed embedding vector for an
// observation (P4). blob is the little-endian float32 encoding.
func (s *Service) SetObservationEmbedding(ctx context.Context, projectID string, id int64, blob []byte, model string) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.SetEmbedding(ctx, id, blob, model)
}

// SearchAllProjectsInProjects runs the FTS query over a fixed set of project
// store ids (the caller's explicit list) and merges results by rank. This is
// used when the caller names specific stores to span rather than "all".
func (s *Service) SearchAllProjectsInProjects(ctx context.Context, projectIDs []string, query, matchMode, scope string, limit int) ([]memory.Observation, error) {
	if limit <= 0 {
		limit = 20
	}
	var collected []memory.Observation
	for _, pid := range projectIDs {
		res, err := s.SearchObservationsScoped(ctx, pid, query, matchMode, scope, limit)
		if err != nil {
			continue
		}
		collected = append(collected, res...)
	}
	// de-dup by (project, id)
	seen := map[string]bool{}
	out := make([]memory.Observation, 0, len(collected))
	for _, o := range collected {
		k := o.Project + "|" + strconv.FormatInt(o.ID, 10)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, o)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].UpdatedAt > out[j].UpdatedAt })
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// PinObservation marks an observation as pinned (mem_pin).
func (s *Service) PinObservation(ctx context.Context, projectID string, id int64) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.Pin(ctx, id)
}

// UnpinObservation clears the pinned flag (mem_unpin).
func (s *Service) UnpinObservation(ctx context.Context, projectID string, id int64) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.Unpin(ctx, id)
}

// TTLRetire soft-deletes expired observations for the project and returns how
// many were retired. Backs mem_review(action=retire_expired) and any
// maintenance path that wants a single sweep.
func (s *Service) TTLRetire(ctx context.Context, projectID string) (int, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return 0, err
	}
	defer cleanup()
	return h.memory.TTLRetire(ctx)
}

// TTLPending returns the count of observations that are past their expires_at
// timestamp and have not been soft-deleted (diagnostic; feeds mem_doctor).
func (s *Service) TTLPending(ctx context.Context, projectID string) (int, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return 0, err
	}
	defer cleanup()
	return h.memory.TTLSoftExpiry(ctx)
}

// Unify consolidates one or more source project stores into a single canonical
// project. It is the admin-facing wrapper around MergeProjects for cases where
// the caller wants to fold several names at once (e.g. three directory-hash
// variants of the same repo). For each source it records the alias and
// copies + re-tags rows, so a single mem_search(all_projects=false, project=
// canonical) then returns the combined history.
func (s *Service) Unify(ctx context.Context, canonical string, sources ...string) (int, error) {
	canonical = strings.TrimSpace(canonical)
	if canonical == "" {
		return 0, fmt.Errorf("canonical project is required")
	}
	total := 0
	for _, src := range sources {
		src = strings.TrimSpace(src)
		if src == "" || src == canonical {
			continue
		}
		moved, _, err := s.MergeProjects(ctx, src, canonical)
		if err != nil {
			return total, err
		}
		total += moved
	}
	return total, nil
}

// UpdateObservation modifies an existing observation by ID. Only non-empty
// fields in in are applied. Bumps updated_at; FTS trigger keeps the index
// in sync.
func (s *Service) UpdateObservation(ctx context.Context, projectID string, id int64, in memory.UpdateInput) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.Update(ctx, id, in)
}

// DeleteObservation removes an observation. Soft-delete by default; hard
// when hardDelete is true.
func (s *Service) DeleteObservation(ctx context.Context, projectID string, id int64, hardDelete bool) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.Delete(ctx, id, hardDelete)
}

// Timeline returns the progressive-disclosure window around an observation.
func (s *Service) ObservationTimeline(ctx context.Context, projectID string, anchorID int64, window time.Duration, limit int) (memory.Timeline, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return memory.Timeline{}, err
	}
	defer cleanup()
	return h.memory.Timeline(ctx, anchorID, window, limit)
}

// ListReviews returns observations due for review, oldest review_after first.
func (s *Service) ListReviews(ctx context.Context, projectID string, limit int) ([]memory.ReviewDue, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.memory.ListReviews(ctx, limit)
}

// MarkReviewed advances an observation's review cycle.
func (s *Service) MarkReviewReviewed(ctx context.Context, projectID string, id int64) (string, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return "", err
	}
	defer cleanup()
	return h.memory.MarkReviewed(ctx, id)
}

// SetReviewAfter sets the review_after for an observation.
func (s *Service) SetObservationReviewAfter(ctx context.Context, projectID string, id int64, reviewAfter string) error {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return err
	}
	defer cleanup()
	return h.memory.SetReviewAfter(ctx, id, reviewAfter)
}

// ProjectInfo describes the resolved project for cwd: id, source, all known
// projects, plus — when cwd is a parent of several git repositories — the
// candidate list so the caller (typically the agent) can pick one and retry
// the write with the chosen name.
type ProjectInfo struct {
	Project           string   `json:"project"`
	Source            string   `json:"source"`
	Directory         string   `json:"directory"`
	Projects          []string `json:"projects"`
	Ambiguous         bool     `json:"ambiguous,omitempty"`
	AvailableProjects []string `json:"available_projects,omitempty"`
	Warning           string   `json:"warning,omitempty"`
	SeedID            string   `json:"seed_id,omitempty"`
}

// CurrentProject returns the resolved project for cwd along with all
// available projects.
func (s *Service) CurrentProject(directory string) (ProjectInfo, error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return ProjectInfo{}, err
	}
	res, resErr := project.ResolveDetailed(absDir)
	projects, err := s.ListProjects()
	if err != nil {
		return ProjectInfo{}, err
	}
	out := ProjectInfo{
		Project:   res.ID,
		Source:    string(res.Source),
		Directory: absDir,
		Projects:  projects,
	}
	if resErr != nil {
		var amb *project.AmbiguousProjectError
		if errors.As(resErr, &amb) {
			out.Ambiguous = true
			out.AvailableProjects = res.Available
		}
	}
	if res.Warning != "" {
		out.Warning = res.Warning
	}
	if res.SeedID != "" {
		out.SeedID = res.SeedID
	}
	return out, nil
}

// MemoryDoctor describes mnemonic store health for a project.
type MemoryDoctor struct {
	SchemaVersion   int            `json:"schema_version"`
	WALMode         string         `json:"wal_mode,omitempty"`
	Observations    int            `json:"observations"`
	ObservationsFTS int            `json:"observations_fts"`
	Files           int            `json:"files"`
	Chunks          int            `json:"chunks"`
	ChunksFTS       int            `json:"chunks_fts"`
	WebCache        int            `json:"web_cache"`
	WebCacheFTS     int            `json:"web_cache_fts"`
	Prompts         int            `json:"prompts"`
	ByType          map[string]int `json:"by_type"`
	DiskSizeBytes   int64          `json:"disk_size_bytes"`
	FTSIntegrityOK  bool           `json:"fts_integrity_ok"`
	FTSDrift        int            `json:"fts_drift"`
}

// MemoryDoctor runs read-only diagnostics: schema count, FTS row counts
// and drift, WAL state, and on-disk size.
func (s *Service) MemoryDoctor(ctx context.Context, projectID string) (MemoryDoctor, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return MemoryDoctor{}, err
	}
	defer cleanup()

	out := MemoryDoctor{}
	if err := h.store.DB.QueryRowContext(ctx, `SELECT schema_version FROM index_meta WHERE key='schema_version'`).Scan(&out.SchemaVersion); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return out, fmt.Errorf("schema_version: %w", err)
		}
	}
	if err := h.store.DB.QueryRowContext(ctx, `PRAGMA journal_mode`).Scan(&out.WALMode); err != nil {
		return out, fmt.Errorf("journal_mode: %w", err)
	}
	byType := map[string]int{}
	h.rowCount(ctx, "observations", &out.Observations)
	h.rowCount(ctx, "files", &out.Files)
	h.rowCount(ctx, "chunks", &out.Chunks)
	h.rowCount(ctx, "web_cache", &out.WebCache)
	h.rowCount(ctx, "prompts", &out.Prompts)

	ftsObs, ftsChunks, ftsWeb := int64(0), int64(0), int64(0)
	_ = h.store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM observations_fts`).Scan(&ftsObs)
	_ = h.store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM chunks_fts`).Scan(&ftsChunks)
	_ = h.store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM web_cache_fts`).Scan(&ftsWeb)
	out.ObservationsFTS = int(ftsObs)
	out.ChunksFTS = int(ftsChunks)
	out.WebCacheFTS = int(ftsWeb)

	rows, err := h.store.DB.QueryContext(ctx, `
		SELECT type, COUNT(*) FROM observations
		WHERE project = ? AND deleted_at IS NULL GROUP BY type`, projectID)
	if err == nil {
		for rows.Next() {
			var t string
			var c int
			if rows.Scan(&t, &c) == nil {
				byType[t] = c
			}
		}
		rows.Close()
	}
	out.ByType = byType

	if info, err := os.Stat(h.store.Path()); err == nil {
		out.DiskSizeBytes = info.Size()
	}
	out.FTSDrift = int(ftsObs) - out.Observations
	out.FTSIntegrityOK = out.FTSDrift >= 0
	return out, nil
}

// rowCount helper (keeps MemoryDoctor readable).
func (h *projectHandle) rowCount(ctx context.Context, table string, out *int) {
	_ = h.store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(out)
}

// SymbolHitDTO is a public identifier-FTS hit.
type SymbolHitDTO = search.SymbolHit

// SymbolSearch runs identifier-aware FTS over indexed symbols.
func (s *Service) SymbolSearch(ctx context.Context, projectID, query string, limit int) ([]search.SymbolHit, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return search.SymbolFTS(h.store.DB, query, limit)
}

// OrientResult is the Tier-1 orientation answer for one symbol: its metadata,
// the file TOC (all symbols in the file), a signature, and linked rationale.
type OrientResult struct {
	Found     bool             `json:"found"`
	Symbol    map[string]any   `json:"symbol,omitempty"`
	FileTOC   []map[string]any `json:"file_toc,omitempty"`
	Signature string           `json:"signature,omitempty"`
	List      []map[string]any `json:"list,omitempty"`
	Rationale []map[string]any `json:"rationale,omitempty"`
	Reason    string           `json:"reason,omitempty"` // not-found note
}

// OrientSymbol returns Tier-1 orientation for a resolved symbol: signature,
// file TOC, map, list, and metadata, plus linked rationale. An unknown symbol
// returns Found=false with a not-found reason (no fabricated symbol).
func (s *Service) OrientSymbol(ctx context.Context, projectID, symbol string) (*OrientResult, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return orientSymbol(h.store.DB, symbol)
}

// CodeGrep runs a structural by-example grep over root, index-free.
func (s *Service) CodeGrep(ctx context.Context, root, pattern string) (*search.GrepResult, error) {
	if pattern == "" {
		return nil, fmt.Errorf("code_grep: pattern is required")
	}
	return search.GrepByExample(root, pattern)
}

// orientSymbol resolves symbol (exact or FTS) and returns its orientation.
func orientSymbol(db *sql.DB, symbol string) (*OrientResult, error) {
	// Resolve the symbol: exact name match first (deterministic), then
	// identifier-FTS. Multiple exact matches are the first (lowest id) —
	// orientation is a single-symbol answer, not a candidate list (that is
	// step 03's code_impact job).
	row := db.QueryRow(`
		SELECT s.id, s.name, s.qualified_name, s.kind, s.language, s.signature,
		       s.start_line, s.end_line, f.path, s.file_id
		FROM symbols s INNER JOIN files f ON f.id = s.file_id
		WHERE s.name = ?
		ORDER BY s.id LIMIT 1`, symbol)
	var id, fileID int64
	var name, qualified, kind, lang, sig string
	var startLine, endLine int
	var path string
	err := row.Scan(&id, &name, &qualified, &kind, &lang, &sig, &startLine, &endLine, &path, &fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Fallback: identifier-FTS.
			hits, e2 := search.SymbolFTS(db, symbol, 1)
			if e2 != nil {
				return nil, e2
			}
			if len(hits) == 0 {
				return &OrientResult{Found: false, Reason: "symbol not found: " + symbol}, nil
			}
			hit := hits[0]
			row2 := db.QueryRow(`
				SELECT s.id, s.name, s.qualified_name, s.kind, s.language, s.signature,
				       s.start_line, s.end_line, f.path, s.file_id
				FROM symbols s INNER JOIN files f ON f.id = s.file_id
				WHERE s.id = ?`, hit.ID)
			if err := row2.Scan(&id, &name, &qualified, &kind, &lang, &sig, &startLine, &endLine, &path, &fileID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	out := &OrientResult{
		Found: true,
		Symbol: map[string]any{
			"id":             id,
			"name":           name,
			"qualified_name": qualified,
			"kind":           kind,
			"language":       lang,
			"path":           path,
			"start_line":     startLine,
			"end_line":       endLine,
		},
		Signature: sig,
	}

	// File TOC + list: every symbol in the file, ordered by start line.
	tocRows, err := db.Query(`
		SELECT id, name, qualified_name, kind, start_line, end_line
		FROM symbols WHERE file_id = ? ORDER BY start_line, id`, fileID)
	if err != nil {
		return out, nil
	}
	for tocRows.Next() {
		var tID int64
		var tName, tQualified, tKind string
		var tStart, tEnd int
		if err := tocRows.Scan(&tID, &tName, &tQualified, &tKind, &tStart, &tEnd); err != nil {
			tocRows.Close()
			return out, nil
		}
		entry := map[string]any{
			"id":         tID,
			"name":       tName,
			"kind":       tKind,
			"start_line": tStart,
			"end_line":   tEnd,
		}
		out.FileTOC = append(out.FileTOC, entry)
	}
	tocRows.Close()
	out.List = out.FileTOC

	// Rationale linked to this symbol.
	rationaleRows, err := db.Query(`SELECT text, kind, line FROM rationale WHERE symbol_id = ? ORDER BY line`, id)
	if err == nil {
		for rationaleRows.Next() {
			var text, kind string
			var line int
			if rationaleRows.Scan(&text, &kind, &line) == nil {
				out.Rationale = append(out.Rationale, map[string]any{
					"text": text,
					"kind": kind,
					"line": line,
				})
			}
		}
		rationaleRows.Close()
	}
	return out, nil
}

// CodeStatus returns index stats and whether the index is stale.
func (s *Service) CodeStatus(ctx context.Context, projectID string) (codeindex.Status, bool, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return codeindex.Status{}, false, err
	}
	defer cleanup()
	status, err := codeindex.GetStatus(h.store)
	if err != nil {
		return codeindex.Status{}, false, err
	}
	stale := status.FileCount == 0 || status.LastIndexed == ""
	return status, stale, nil
}

// CodeSearch runs BM25 FTS over indexed code chunks.
func (s *Service) CodeSearch(ctx context.Context, projectID, query string, limit int) ([]search.CodeHit, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return search.CodeSearch(h.store.DB, query, limit)
}

// CodeFiles returns all indexed file paths, sorted.
func (s *Service) CodeFiles(ctx context.Context, projectID string) ([]string, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	rows, err := h.store.DB.QueryContext(ctx, `SELECT path FROM files ORDER BY path`)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}
	defer rows.Close()
	var paths []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return paths, nil
}

// ReadIndexedCode returns indexed source for path and optional line range.
func (s *Service) ReadIndexedCode(ctx context.Context, projectID, path string, startLine, endLine int) (map[string]any, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return readIndexedCode(h.store.DB, path, startLine, endLine)
}

// ImpactOptions tunes a blast-radius traversal (narrowing + confidence).
type ImpactOptions struct {
	File          string
	UID           string
	Kind          string
	MinConfidence string
	MaxDepth      int
}

// ImpactResultDTO is the code_impact answer. Either Target+tiers is set (a
// single resolved symbol) or Candidates is set (a ranked ambiguous list —
// never a silent pick).
type ImpactResultDTO struct {
	Found      bool                 `json:"found"`
	Ambiguous  bool                 `json:"ambiguous,omitempty"`
	Candidates []graph.Symbol       `json:"candidates,omitempty"`
	Target     *graph.Symbol        `json:"target,omitempty"`
	WillBreak  []graph.ImpactEdge   `json:"will_break,omitempty"`
	Likely     []graph.ImpactEdge   `json:"likely_affected,omitempty"`
	Excluded   int                  `json:"excluded_low_confidence,omitempty"`
	Reason     string               `json:"reason,omitempty"`
}

func (r *ImpactResultDTO) Summary() string {
	if r == nil {
		return ""
	}
	if r.Ambiguous {
		return fmt.Sprintf("ambiguous: %d candidates (narrow with file/uid/kind)", len(r.Candidates))
	}
	return fmt.Sprintf("will_break: %d, likely_affected: %d", len(r.WillBreak), len(r.Likely))
}

// Impact is the risk-tiered blast-radius DTO for the graph package result.
type Impact = graph.ImpactResult

// CodeImpact returns the risk-tiered blast radius for symbol, honoring the
// narrowing + confidence options. A name matching several symbols returns a
// ranked candidate list (never a silent pick); an unknown symbol returns an
// empty (not-found) result.
func (s *Service) CodeImpact(ctx context.Context, projectID, symbol string, opts ImpactOptions) (*ImpactResultDTO, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	res, err := graph.Resolve(ctx, h.store.DB, symbol, graph.ResolveFilter{
		File: opts.File, UID: opts.UID, Kind: opts.Kind,
	})
	if err != nil {
		return nil, err
	}
	if res.NotFound {
		return &ImpactResultDTO{Found: false, Reason: "symbol not found: " + symbol}, nil
	}
	if res.Ambiguous {
		ranked, err := graph.RankCandidates(ctx, h.store.DB, res.Matches)
		if err != nil {
			return nil, err
		}
		return &ImpactResultDTO{Found: true, Ambiguous: true, Candidates: ranked}, nil
	}
	impact, err := graph.Impact(ctx, h.store.DB, res.Target, graph.ImpactOptions{
		MinConfidence: opts.MinConfidence,
		MaxDepth:      opts.MaxDepth,
	})
	if err != nil {
		return nil, err
	}
	target := res.Target
	return &ImpactResultDTO{
		Found:     true,
		Target:    &target,
		WillBreak: impact.WillBreak,
		Likely:    impact.Likely,
		Excluded:  impact.Excluded,
	}, nil
}

// NeighborsDTO is the graph neighbor answer (every edge confidence-labeled).
type NeighborsDTO struct {
	Symbol *graph.Symbol `json:"symbol,omitempty"`
	Edges  []graph.Edge  `json:"edges"`
	Reason string        `json:"reason,omitempty"`
}

// GrabNeighbors returns a symbol's confidence-labeled edges for view. An
// unknown symbol returns not-found (no invented edges).
func (s *Service) GrabNeighbors(ctx context.Context, projectID, view, symbol string) (*NeighborsDTO, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	res, err := graph.Resolve(ctx, h.store.DB, symbol, graph.ResolveFilter{})
	if err != nil {
		return nil, err
	}
	if res.NotFound {
		return &NeighborsDTO{Edges: []graph.Edge{}, Reason: "symbol not found: " + symbol}, nil
	}
	if res.Ambiguous {
		// Report ambiguity as a not-found-with-candidates rather than a
		// silent pick; graph views are single-symbol answers.
		return &NeighborsDTO{Edges: []graph.Edge{}, Reason: fmt.Sprintf("ambiguous symbol %q: %d candidates (narrow with file/uid/kind)", symbol, len(res.Matches))}, nil
	}
	v := graph.View(view)
	edges, err := graph.Neighbors(ctx, h.store.DB, res.Target, v)
	if err != nil {
		return nil, err
	}
	sym := res.Target
	return &NeighborsDTO{Symbol: &sym, Edges: edges}, nil
}

// PathDTO is the code_path answer.
type PathDTO struct {
	Found      bool             `json:"found"`
	Path       []graph.Edge     `json:"path,omitempty"`
	GraphStops *graph.GraphStops `json:"graph_stops,omitempty"`
	Reason     string           `json:"reason,omitempty"`
}

// CodePath returns the shortest edge path between two symbols or a
// where-the-graph-stops answer. Unknown endpoints yield a not-found reason.
func (s *Service) CodePath(ctx context.Context, projectID, from, to string) (*PathDTO, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	resFrom, err := graph.Resolve(ctx, h.store.DB, from, graph.ResolveFilter{})
	if err != nil {
		return nil, err
	}
	resTo, err := graph.Resolve(ctx, h.store.DB, to, graph.ResolveFilter{})
	if err != nil {
		return nil, err
	}
	if resFrom.NotFound || resTo.NotFound {
		missing := from
		if !resFrom.NotFound {
			missing = to
		}
		return &PathDTO{Found: false, Reason: "symbol not found: " + missing}, nil
	}
	res, err := graph.Path(ctx, h.store.DB, resFrom.Target, resTo.Target)
	if err != nil {
		return nil, err
	}
	return &PathDTO{Found: res.Found, Path: res.Path, GraphStops: res.GraphStops}, nil
}

// ExplainDTO is the code_explain answer.
type ExplainDTO struct {
	Found       bool          `json:"found"`
	Symbol      *graph.Symbol `json:"symbol,omitempty"`
	Degree      int           `json:"degree"`
	Connections []graph.Conn  `json:"connections"`
	Reason      string        `json:"reason,omitempty"`
}

// CodeExplain returns a symbol's node, degree, and connections ranked by the
// neighbor's degree. Unknown symbol returns not-found.
func (s *Service) CodeExplain(ctx context.Context, projectID, symbol string) (*ExplainDTO, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	res, err := graph.Resolve(ctx, h.store.DB, symbol, graph.ResolveFilter{})
	if err != nil {
		return nil, err
	}
	if res.NotFound {
		return &ExplainDTO{Found: false, Connections: []graph.Conn{}, Reason: "symbol not found: " + symbol}, nil
	}
	out, err := graph.Explain(ctx, h.store.DB, res.Target)
	if err != nil {
		return nil, err
	}
	sym := res.Target
	return &ExplainDTO{Found: true, Symbol: &sym, Degree: out.Degree, Connections: out.Connections}, nil
}

// FairCoverageDTO is the per-language fair-coverage field for code_status.
type FairCoverageDTO = graph.CoverageLang

// CodeFairCoverage returns measured per-language fair coverage from edges.
func (s *Service) CodeFairCoverage(ctx context.Context, projectID string) (map[string]FairCoverageDTO, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return graph.FairCoverage(ctx, h.store.DB)
}

// ExploreSourceSpan is one symbol's verbatim source slice.
type ExploreSourceSpan struct {
	Symbol    string
	StartLine int
	EndLine   int
	Content   string
}

// ExploreFlowEdge is one call-flow hop between returned symbols.
type ExploreFlowEdge struct {
	From       string
	To         string
	Kind       string
	Confidence string
	Line       int
}

// ExploreResult is the composite code_explore answer.
type ExploreResult struct {
	Symbol string
	Source map[string][]ExploreSourceSpan
	Flow   []ExploreFlowEdge
	Impact *ImpactResultDTO
}

// CodeExplore assembles the composite answer: the symbol's source (grouped by
// file), its call-flow neighbors (including INFERRED dynamic-dispatch hops),
// and a blast-radius summary.
func (s *Service) CodeExplore(ctx context.Context, projectID, symbol string) (*ExploreResult, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	res, err := graph.Resolve(ctx, h.store.DB, symbol, graph.ResolveFilter{})
	if err != nil {
		return nil, err
	}
	out := &ExploreResult{Symbol: symbol, Source: map[string][]ExploreSourceSpan{}}
	if res.NotFound {
		return out, nil
	}
	target := res.Target
	// Pull verbatim source for the target symbol.
	spans, err := s.readSymbolSpans(ctx, h.store.DB, []graph.Symbol{target})
	if err == nil {
		for path, ss := range spans {
			out.Source[path] = ss
		}
	}
	// Gather the call-flow: forward (callees) and reverse (callers) edges,
	// including INFERRED dynamic-dispatch hops.
	callees, _ := graph.Neighbors(ctx, h.store.DB, target, graph.ViewCallees)
	callers, _ := graph.Neighbors(ctx, h.store.DB, target, graph.ViewCallers)
	related := []graph.Symbol{}
	seen := map[int64]bool{target.ID: true}
	for _, e := range append(append([]graph.Edge{}, callees...), callers...) {
		other := e.To
		if e.From.ID == target.ID && e.To.ID != target.ID {
			other = e.To
		}
		if e.From.ID != target.ID && e.To.ID != target.ID {
			other = e.To
		}
		if other.ID == 0 {
			continue
		}
		out.Flow = append(out.Flow, ExploreFlowEdge{
			From:       nameOf(target, other),
			To:         nameOf(other, target),
			Kind:       e.Kind,
			Confidence: e.Confidence,
			Line:       e.Line,
		})
		if !seen[other.ID] {
			seen[other.ID] = true
			related = append(related, other)
		}
	}
	if len(related) > 0 {
		if spans, err := s.readSymbolSpans(ctx, h.store.DB, related); err == nil {
			for path, ss := range spans {
				out.Source[path] = append(out.Source[path], ss...)
			}
		}
	}
	out.Impact, err = s.CodeImpact(ctx, projectID, symbol, ImpactOptions{})
	if err != nil {
		out.Impact = &ImpactResultDTO{}
	}
	return out, nil
}

// nameOf returns the from-side name for a flow hop relative to target.
func nameOf(a, b graph.Symbol) string {
	if a.ID != 0 {
		return a.Name
	}
	if b.ID != 0 {
		return b.Name
	}
	return ""
}

// readSymbolSpans returns each symbol's verbatim source grouped by file.
func (s *Service) readSymbolSpans(ctx context.Context, db *sql.DB, syms []graph.Symbol) (map[string][]ExploreSourceSpan, error) {
	out := map[string][]ExploreSourceSpan{}
	for _, sym := range syms {
		if sym.ID == 0 {
			continue
		}
		res, err := readIndexedCode(db, sym.Path, sym.StartLine, sym.EndLine)
		if err != nil {
			continue
		}
		text, _ := res["text"].(string)
		out[sym.Path] = append(out[sym.Path], ExploreSourceSpan{
			Symbol:    sym.Name,
			StartLine: sym.StartLine,
			EndLine:   sym.EndLine,
			Content:   text,
		})
	}
	return out, nil
}

// RunCodeIndex runs incremental code indexing for directory.
func (s *Service) RunCodeIndex(ctx context.Context, directory string) (codeindex.Stats, error) {
	h, cleanup, err := s.openProjectForDirectory(directory)
	if err != nil {
		return codeindex.Stats{}, err
	}
	defer cleanup()
	cfg := config.Load(directory)
	idxCfg := codeindex.Config{
		Include:      cfg.Include,
		Exclude:      cfg.Exclude,
		ChunkLines:   cfg.ChunkLines,
		ChunkOverlap: cfg.ChunkOverlap,
		MaxFileSize:  cfg.MaxFileSize,
	}
	return codeindex.New(h.store).Run(ctx, directory, idxCfg)
}

// CodeHybridResult is the hybrid code search answer (per-signal provenance on
// every hit).
type CodeHybridResult = hybrid.Result

// CodeHybridSearch runs the offline hybrid code search (FTS + deterministic
// signals + optional embeddings) for projectID. The embedder is resolved from
// config; a down/missing embedder degrades to FTS + signals (the floor).
func (s *Service) CodeHybridSearch(ctx context.Context, projectID, query string, limit int) (*hybrid.Result, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	emb := resolveEmbedder(h.root)
	return hybrid.Search(ctx, h.store.DB, query, hybrid.Options{
		Limit:    limit,
		Embedder: emb,
	})
}

// CodeSemanticSearch runs the vector leg only: symbol-level hits return the
// named symbol + file + line; chunk-level hits return the line range. With no
// active embedder it returns an empty result (the hybrid tool keeps the
// FTS+signals floor).
func (s *Service) CodeSemanticSearch(ctx context.Context, projectID, query string, limit int) (*hybrid.Result, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	emb := resolveEmbedder(h.root)
	return hybrid.Search(ctx, h.store.DB, query, hybrid.Options{
		Limit:    limit,
		Semantic: true,
		Embedder: emb,
	})
}

// CodeEmbeddingStatus reports the active provider/model, vector dimension,
// embedded counts, and model-swap guard state (read-only).
func (s *Service) CodeEmbeddingStatus(ctx context.Context, projectID string) (map[string]any, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	cfg := config.Load(h.root)
	emb := resolveEmbedder(h.root)
	status := map[string]any{
		"provider": cfg.Embedder.Provider,
	}
	if emb == nil {
		status["active"] = false
		status["model"] = ""
		status["dimension"] = 0
		status["reason"] = "embedder off"
	} else {
		status["active"] = true
		status["model"] = emb.Model()
		status["dimension"] = emb.Dimension()
	}
	var symCount int
	_ = h.store.DB.QueryRow(`SELECT COUNT(*) FROM embeddings`).Scan(&symCount)
	status["embedded_symbols"] = symCount
	var model sql.NullString
	_ = h.store.DB.QueryRow(`SELECT value FROM embed_meta WHERE key = 'embedding_model'`).Scan(&model)
	status["indexed_model"] = model.String
	return status, nil
}

// toAsym converts config params to embedder params.
func toAsym(p config.EmbedderParams) embedder.AsymParams {
	return embedder.AsymParams{
		Instructions: p.Instructions,
		InputType:    p.InputType,
		MaxTokens:    p.MaxTokens,
	}
}

// resolveEmbedder builds the process embedder from config. An empty/nil
// embedder means "off" — the FTS+signals floor applies.
func resolveEmbedder(configRoot string) embedder.Embedder {
	cfg := config.Load(configRoot)
	switch cfg.Embedder.Provider {
	case "external":
		return embedder.NewExternal(embedder.ExternalConfig{
			BaseURL:   cfg.Embedder.BaseURL,
			Model:     cfg.Embedder.Model,
			APIKey:    cfg.Embedder.APIKey,
			Dimension: cfg.Embedder.Dimension,
			Indexing:  toAsym(cfg.Embedder.Indexing),
			Query:     toAsym(cfg.Embedder.Query),
		})
	case "off", "":
		return nil
	default: // "onnx" is the default
		return embedder.NewOnnx(embedder.OnnxConfig{
			Model:     cfg.Embedder.Model,
			Dimension: cfg.Embedder.Dimension,
			Indexing:  toAsym(cfg.Embedder.Indexing),
			Query:     toAsym(cfg.Embedder.Query),
		})
	}
}

// WebLookup checks the web cache for a matching entry.
func (s *Service) WebLookup(ctx context.Context, projectID string, in webcache.LookupInput) (webcache.LookupResult, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return webcache.LookupResult{}, err
	}
	defer cleanup()
	return h.web.Lookup(ctx, in)
}

// WebSave persists a web research snapshot.
func (s *Service) WebSave(ctx context.Context, projectID string, in webcache.SaveWebInput) (int64, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return 0, err
	}
	defer cleanup()
	return h.web.Save(ctx, in)
}

// WebSearch runs FTS over cached web snapshots.
func (s *Service) WebSearch(ctx context.Context, projectID, query, source string, freshOnly bool, limit int) ([]webcache.WebHit, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return h.web.Search(ctx, query, source, freshOnly, limit)
}

// WebGet returns a full cached snapshot by ID.
func (s *Service) WebGet(ctx context.Context, projectID string, id int64) (webcache.WebEntry, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return webcache.WebEntry{}, err
	}
	defer cleanup()
	return h.web.Get(ctx, id)
}

// MemoryStatus returns memory store health stats.
func (s *Service) MemoryStatus(ctx context.Context, projectID string) (memory.Status, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return memory.Status{}, err
	}
	defer cleanup()
	return h.memory.Status(ctx)
}

// WebCacheStatus returns web cache health stats.
func (s *Service) WebCacheStatus(ctx context.Context, projectID string) (webcache.Status, error) {
	h, cleanup, err := s.openProject(projectID, ".")
	if err != nil {
		return webcache.Status{}, err
	}
	defer cleanup()
	return h.web.CacheStatus(ctx)
}

// OpenForCWD opens project-scoped services from the current working directory.
func (s *Service) OpenForCWD() (*projectHandle, func(), error) {
	return s.openProjectFromCWD()
}

// OpenForDirectory opens project-scoped services for directory.
func (s *Service) OpenForDirectory(directory string) (*projectHandle, func(), error) {
	return s.openProjectForDirectory(directory)
}

// ProjectID returns the project ID for an open handle.
func (h *projectHandle) ProjectID() string { return h.projectID }

// Memory returns the memory service for an open handle.
func (h *projectHandle) Memory() *memory.Service { return h.memory }

// Web returns the webcache service for an open handle.
func (h *projectHandle) Web() *webcache.Service { return h.web }

// Store returns the underlying store for an open handle.
func (h *projectHandle) Store() *store.Store { return h.store }

func readIndexedCode(db *sql.DB, path string, startLine, endLine int) (map[string]any, error) {
	var fileID int64
	err := db.QueryRow(`SELECT id FROM files WHERE path = ?`, path).Scan(&fileID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not indexed: %s", path)
	}
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if startLine > 0 {
		if endLine <= 0 {
			endLine = startLine
		}
		rows, err = db.Query(`
			SELECT start_line, end_line, text FROM chunks
			WHERE file_id = ? AND start_line <= ? AND end_line >= ?
			ORDER BY start_line`,
			fileID, endLine, startLine,
		)
	} else {
		rows, err = db.Query(`
			SELECT start_line, end_line, text FROM chunks
			WHERE file_id = ?
			ORDER BY start_line`,
			fileID,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parts []string
	firstLine := 0
	lastLine := 0
	for rows.Next() {
		var chunkStart, chunkEnd int
		var text string
		if err := rows.Scan(&chunkStart, &chunkEnd, &text); err != nil {
			return nil, err
		}
		if firstLine == 0 || chunkStart < firstLine {
			firstLine = chunkStart
		}
		if chunkEnd > lastLine {
			lastLine = chunkEnd
		}
		parts = append(parts, text)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("no indexed chunks for %s", path)
	}
	return map[string]any{
		"path":       path,
		"start_line": firstLine,
		"end_line":   lastLine,
		"text":       strings.Join(parts, "\n"),
	}, nil
}
