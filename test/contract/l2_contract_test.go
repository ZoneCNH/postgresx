package contract_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
	"github.com/ZoneCNH/postgresx/test/postgresxtest"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const l2ReleaseLevel = "L2-T2"

var (
	requiredProfiles    = []string{"unit", "contract", "integration"}
	requiredP0Contracts = []string{
		"sql.exec",
		"sql.query_row",
		"sql.query_many",
		"sql.not_found",
		"sql.syntax_error",
		"sql.unique_violation",
		"sql.foreign_key_violation",
		"sql.context_timeout",
		"tx.commit",
		"tx.rollback",
		"tx.rollback_on_error",
		"pool.exhaustion",
	}
	requiredHardFailures = []string{
		"secret_leak",
		"layer_violation",
		"missing_required_contract",
		"missing_required_evidence",
		"race_detected",
		"goroutine_leak",
		"release_level_overclaimed",
	}
	requiredEvidence = []string{
		".agent/evidence/raw/unit-test.json",
		".agent/evidence/raw/contract-test.json",
		".agent/evidence/raw/integration-test.json",
		".agent/evidence/normalized/contract-check.json",
		".agent/evidence/normalized/integration-check.json",
		".agent/evidence/normalized/layer-guard.json",
		".agent/evidence/normalized/secret-scan.json",
		".agent/evidence/decision/test-plan.json",
		".agent/evidence/decision/release-readiness.json",
		".agent/evidence/trace/traceability-matrix.json",
		".agent/evidence/retrospective.json",
		".agent/evidence/manifest.json",
	}
)

func TestP0SQLContract(t *testing.T) {
	ctx := t.Context()
	queryer := &postgresxtest.QueryAdapter{
		ExecTag:   postgresxtest.CommandTag{Rows: 2},
		QueryRows: &postgresxtest.Rows{Rows: [][]any{{1, "one"}, {2, "two"}}},
		Row:       &postgresxtest.Row{Values: []any{3, "three"}},
	}

	tag, err := queryer.Exec(ctx, "insert into l2_contract(id) values($1)", 1)
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if got := tag.RowsAffected(); got != 2 {
		t.Fatalf("RowsAffected() = %d, want 2", got)
	}

	rows, err := queryer.Query(ctx, "select id, name from l2_contract where id > $1", 0)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	defer rows.Close()
	var seen []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			t.Fatalf("Rows.Scan() error = %v", err)
		}
		seen = append(seen, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows.Err() = %v", err)
	}
	if !reflect.DeepEqual(seen, []string{"one", "two"}) {
		t.Fatalf("rows = %v, want [one two]", seen)
	}

	var id int
	var name string
	if err := queryer.QueryRow(ctx, "select id, name from l2_contract where id=$1", 3).Scan(&id, &name); err != nil {
		t.Fatalf("QueryRow().Scan() error = %v", err)
	}
	if id != 3 || name != "three" {
		t.Fatalf("row = (%d, %q), want (3, three)", id, name)
	}

	wantCalls := []postgresxtest.QueryCall{
		{Operation: "exec", SQL: "insert into l2_contract(id) values($1)", Args: []any{1}},
		{Operation: "query", SQL: "select id, name from l2_contract where id > $1", Args: []any{0}},
		{Operation: "query_row", SQL: "select id, name from l2_contract where id=$1", Args: []any{3}},
	}
	if !reflect.DeepEqual(queryer.Calls, wantCalls) {
		t.Fatalf("calls = %#v, want %#v", queryer.Calls, wantCalls)
	}
}

func TestP0SQLContractErrorPropagation(t *testing.T) {
	ctx := t.Context()
	execErr := errors.New("exec failed")
	queryErr := errors.New("query failed")
	rowErr := errors.New("row failed")
	rowsErr := errors.New("rows failed")

	queryer := &postgresxtest.QueryAdapter{ExecErr: execErr}
	if _, err := queryer.Exec(ctx, "delete from l2_contract where id=$1", 10); !errors.Is(err, execErr) {
		t.Fatalf("Exec() error = %v, want %v", err, execErr)
	}
	if len(queryer.Calls) != 1 || queryer.Calls[0].Operation != "exec" {
		t.Fatalf("Exec() calls = %#v, want one exec call", queryer.Calls)
	}

	queryer = &postgresxtest.QueryAdapter{QueryErr: queryErr}
	if _, err := queryer.Query(ctx, "select id from l2_contract"); !errors.Is(err, queryErr) {
		t.Fatalf("Query() error = %v, want %v", err, queryErr)
	}

	queryer = &postgresxtest.QueryAdapter{Row: &postgresxtest.Row{Err: rowErr}}
	if err := queryer.QueryRow(ctx, "select id from l2_contract where id=$1", 11).Scan(new(int)); !errors.Is(err, rowErr) {
		t.Fatalf("QueryRow().Scan() error = %v, want %v", err, rowErr)
	}

	queryer = &postgresxtest.QueryAdapter{QueryRows: &postgresxtest.Rows{Rows: [][]any{{12}}, ErrValue: rowsErr}}
	rows, err := queryer.Query(ctx, "select id from l2_contract")
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if !rows.Next() {
		t.Fatal("Rows.Next() = false, want one row before Err()")
	}
	var id int
	if err := rows.Scan(&id); err != nil {
		t.Fatalf("Rows.Scan() error = %v", err)
	}
	if err := rows.Err(); !errors.Is(err, rowsErr) {
		t.Fatalf("Rows.Err() = %v, want %v", err, rowsErr)
	}
}

func TestP0TxContract(t *testing.T) {
	var _ postgresx.Tx = (*postgresxtest.QueryAdapter)(nil)

	ctx := t.Context()
	tx := &postgresxtest.QueryAdapter{ExecTag: postgresxtest.CommandTag{Rows: 1}}
	txFn := postgresx.TxFunc(func(ctx context.Context, tx postgresx.Tx) error {
		tag, err := tx.Exec(ctx, "update l2_contract set seen=true where id=$1", 7)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			t.Fatalf("RowsAffected() = %d, want 1", tag.RowsAffected())
		}
		return nil
	})
	if err := txFn(ctx, tx); err != nil {
		t.Fatalf("TxFunc() error = %v", err)
	}
	if len(tx.Calls) != 1 || tx.Calls[0].Operation != "exec" {
		t.Fatalf("tx calls = %#v, want one exec", tx.Calls)
	}
}

func TestP0TxContractPropagatesExecutorErrors(t *testing.T) {
	ctx := t.Context()
	execErr := errors.New("tx exec failed")
	tx := &postgresxtest.QueryAdapter{ExecErr: execErr}
	txFn := postgresx.TxFunc(func(ctx context.Context, tx postgresx.Tx) error {
		_, err := tx.Exec(ctx, "update l2_contract set seen=true where id=$1", 17)
		return err
	})

	if err := txFn(ctx, tx); !errors.Is(err, execErr) {
		t.Fatalf("TxFunc() error = %v, want %v", err, execErr)
	}
	if len(tx.Calls) != 1 || tx.Calls[0].Operation != "exec" {
		t.Fatalf("tx calls = %#v, want one exec", tx.Calls)
	}
}

func TestP0LivePostgresTxRollbackContract(t *testing.T) {
	ctx := t.Context()
	fixture := postgresxtest.Start(ctx, t, "postgresx-l2-rollback-contract")
	client := fixture.Client()
	sentinel := errors.New("rollback contract sentinel")

	if _, err := client.Exec(ctx, "create temporary table l2_contract_rollback(id integer) on commit preserve rows"); err != nil {
		t.Fatalf("create temporary table: %v", err)
	}
	err := client.WithTx(ctx, func(ctx context.Context, tx postgresx.Tx) error {
		if _, err := tx.Exec(ctx, "insert into l2_contract_rollback(id) values($1)", 1); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("WithTx() error = %v, want sentinel rollback error", err)
	}

	var count int
	if err := client.QueryRow(ctx, "select count(*) from l2_contract_rollback").Scan(&count); err != nil {
		t.Fatalf("count rollback table: %v", err)
	}
	if count != 0 {
		t.Fatalf("rollback table count = %d, want 0", count)
	}
}

func TestP0PoolContract(t *testing.T) {
	cfg := postgresx.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Database = "postgres"
	cfg.User = "postgres"
	cfg.Password = foundationx.NewSecretString("contract-secret")
	cfg.MaxOpenConns = 3
	cfg.MinIdleConns = 1
	cfg.ApplicationName = "postgresx-l2-contract"

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	dsn := cfg.RedactedDSN()
	if strings.Contains(dsn, "contract-secret") {
		t.Fatalf("RedactedDSN() = %q, leaked password", dsn)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse RedactedDSN(): %v", err)
	}
	if masked, ok := parsed.User.Password(); !ok || masked != "***" {
		t.Fatalf("RedactedDSN() password = %q, %v; want masked password", masked, ok)
	}
	sanitized := cfg.Sanitize()
	if sanitized.Password != "***" {
		t.Fatalf("Sanitize().Password = %q, want mask", sanitized.Password)
	}

	statsType := reflect.TypeOf(postgresx.PoolStats{})
	for _, field := range []string{"TotalConns", "IdleConns", "AcquiredConns", "ConstructingConns", "MaxConns"} {
		if _, ok := statsType.FieldByName(field); !ok {
			t.Fatalf("PoolStats missing field %s", field)
		}
	}
}

func TestP0LivePostgresContract(t *testing.T) {
	fixture := postgresxtest.Start(t.Context(), t, "postgresx-l2-contract")
	client := fixture.Client()

	if err := client.Ping(t.Context()); err != nil {
		t.Fatalf("Ping() error = %v", err)
	}
	if stats := client.Stats(); stats.MaxConns <= 0 {
		t.Fatalf("Stats().MaxConns = %d, want positive", stats.MaxConns)
	}
	if err := client.WithTx(t.Context(), func(ctx context.Context, tx postgresx.Tx) error {
		var got int
		if err := tx.QueryRow(ctx, "select 1").Scan(&got); err != nil {
			return err
		}
		if got != 1 {
			t.Fatalf("select 1 = %d, want 1", got)
		}
		return nil
	}); err != nil {
		t.Fatalf("WithTx() error = %v", err)
	}
}

func TestP0ErrorMappingContract(t *testing.T) {
	tests := []struct {
		contract  string
		err       error
		kind      foundationx.ErrorKind
		retryable bool
	}{
		{
			contract:  "sql.not_found",
			err:       pgx.ErrNoRows,
			kind:      foundationx.ErrorKindNotFound,
			retryable: false,
		},
		{
			contract:  "sql.syntax_error",
			err:       &pgconn.PgError{Code: "42601"},
			kind:      foundationx.ErrorKindValidation,
			retryable: false,
		},
		{
			contract:  "sql.unique_violation",
			err:       &pgconn.PgError{Code: "23505"},
			kind:      foundationx.ErrorKindAlreadyExist,
			retryable: false,
		},
		{
			contract:  "sql.foreign_key_violation",
			err:       &pgconn.PgError{Code: "23503"},
			kind:      foundationx.ErrorKindConflict,
			retryable: false,
		},
		{
			contract:  "sql.context_timeout",
			err:       context.DeadlineExceeded,
			kind:      foundationx.ErrorKindTimeout,
			retryable: true,
		},
		{
			contract:  "pool.exhaustion",
			err:       &pgconn.PgError{Code: "53300"},
			kind:      foundationx.ErrorKindUnavailable,
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.contract, func(t *testing.T) {
			err := postgresx.MapError(tt.contract, tt.err)
			if !foundationx.IsKind(err, tt.kind) {
				t.Fatalf("MapError() = %v, want kind %s", err, tt.kind)
			}
			if got := postgresx.IsRetryable(err); got != tt.retryable {
				t.Fatalf("IsRetryable() = %v, want %v", got, tt.retryable)
			}
			if !errors.Is(err, tt.err) {
				t.Fatalf("MapError() does not unwrap contract error %v", tt.err)
			}
		})
	}
}

func TestL2StandardMetadataMatchesXlibStandard(t *testing.T) {
	repoRoot := findRepoRoot(t)
	manifest := readJSONFile[l2CapabilityManifest](t, filepath.Join(repoRoot, ".agent/l2-capabilities.yaml"))
	if manifest.Package != "postgresx" || manifest.Layer != "L2" {
		t.Fatalf("capability manifest package/layer = %s/%s, want postgresx/L2", manifest.Package, manifest.Layer)
	}
	if manifest.StandardSource != "github.com/ZoneCNH/xlib-standard" {
		t.Fatalf("standard_source = %q, want github.com/ZoneCNH/xlib-standard", manifest.StandardSource)
	}
	if manifest.ReleaseLevelTarget != l2ReleaseLevel {
		t.Fatalf("release_level_target = %q, want %q", manifest.ReleaseLevelTarget, l2ReleaseLevel)
	}
	if !reflect.DeepEqual(manifest.ReleaseContract.RequiredProfiles, requiredProfiles) {
		t.Fatalf("manifest required_profiles = %v, want %v", manifest.ReleaseContract.RequiredProfiles, requiredProfiles)
	}
	if manifest.ReleaseContract.ReleaseAllowed {
		t.Fatal("L2-T2 must not set release_allowed=true")
	}
	if manifest.ReleaseContract.FactoryGradeAllowed {
		t.Fatal("L2-T2 must not set factory_grade_allowed=true")
	}
	if manifest.ReleaseContract.MinScore != 75 {
		t.Fatalf("L2-T2 min_score = %d, want 75", manifest.ReleaseContract.MinScore)
	}
	if manifest.Provider.Image != "postgres:16-alpine" {
		t.Fatalf("provider image = %q, want postgres:16-alpine", manifest.Provider.Image)
	}
	assertSameSet(t, "core capabilities", manifest.Capabilities.Core, []string{"common", "sql", "transaction", "pool"})
	assertSameSet(t, "optional capabilities", manifest.Capabilities.Optional, []string{"migration", "advisory_lock", "batch_insert", "copy"})
	assertSameSet(t, "p0 contracts", manifest.P0Contracts, requiredP0Contracts)
	assertSameSet(t, "forbidden runtime dependencies", manifest.ForbiddenRuntimeDependencies, []string{
		"github.com/ZoneCNH/xlib-standard",
		"github.com/ZoneCNH/testkitx",
		"github.com/ZoneCNH/xlibgate",
	})

	registry := readJSONFile[l2ContractPackRegistry](t, filepath.Join(repoRoot, ".agent/registry/l2-contract-packs.yaml"))
	pack, ok := findContractPack(registry.Packs, "postgresx-p0-sql-tx-pool")
	if !ok {
		t.Fatal("contract pack postgresx-p0-sql-tx-pool not found")
	}
	if pack.Package != "postgresx" || pack.Layer != "L2" || pack.ReleaseLevel != l2ReleaseLevel {
		t.Fatalf("contract pack package/layer/release = %s/%s/%s, want postgresx/L2/%s", pack.Package, pack.Layer, pack.ReleaseLevel, l2ReleaseLevel)
	}
	if !reflect.DeepEqual(pack.RequiredProfiles, requiredProfiles) {
		t.Fatalf("pack required_profiles = %v, want %v", pack.RequiredProfiles, requiredProfiles)
	}
	assertSameSet(t, "pack required contracts", contractNames(pack.RequiredContracts), requiredP0Contracts)

	gate := readJSONFile[l2Gate](t, filepath.Join(repoRoot, ".agent/gates/l2gate.yaml"))
	if gate.Package != "postgresx" || gate.Layer != "L2" {
		t.Fatalf("gate package/layer = %s/%s, want postgresx/L2", gate.Package, gate.Layer)
	}
	if gate.ReleaseLevelTarget != l2ReleaseLevel || gate.ReleaseLevelActual != l2ReleaseLevel {
		t.Fatalf("gate release target/actual = %s/%s, want %s/%s", gate.ReleaseLevelTarget, gate.ReleaseLevelActual, l2ReleaseLevel, l2ReleaseLevel)
	}
	if gate.MinScore != 75 || gate.Score != 75 {
		t.Fatalf("gate score/min_score = %d/%d, want 75/75", gate.Score, gate.MinScore)
	}
	if gate.ReleaseAllowed || gate.FactoryGradeAllowed {
		t.Fatalf("L2-T2 gate release_allowed/factory_grade_allowed = %v/%v, want false/false", gate.ReleaseAllowed, gate.FactoryGradeAllowed)
	}
	if !reflect.DeepEqual(gate.RequiredProfiles, requiredProfiles) {
		t.Fatalf("gate required_profiles = %v, want %v", gate.RequiredProfiles, requiredProfiles)
	}
	assertSameSet(t, "gate hard failures", gate.HardFailures, requiredHardFailures)
	assertSameSet(t, "gate required evidence", gate.RequiredEvidence, requiredEvidence)
}

func TestFormalPackagesDoNotDependOnL2TestTooling(t *testing.T) {
	repoRoot := findRepoRoot(t)
	for _, dir := range []string{"pkg", "internal"} {
		root := filepath.Join(repoRoot, dir)
		if _, err := os.Stat(root); err != nil {
			continue
		}
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			text := string(contents)
			for _, forbidden := range []string{"testkitx", "xlibgate", "xlib-standard"} {
				if strings.Contains(text, forbidden) {
					t.Fatalf("formal package file %s contains forbidden L2 tooling dependency marker %q", path, forbidden)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}

	cmd := exec.Command("go", "list", "-deps", "./pkg/postgresx")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps ./pkg/postgresx failed: %v\n%s", err, out)
	}
	deps := string(out)
	for _, forbidden := range []string{
		"github.com/ZoneCNH/xlib-standard",
		"github.com/ZoneCNH/testkitx",
		"github.com/ZoneCNH/xlibgate",
		"github.com/ZoneCNH/x.go",
		"github.com/bytechainx/x.go",
	} {
		if strings.Contains(deps, forbidden) {
			t.Fatalf("formal package dependencies contain forbidden runtime dependency %q", forbidden)
		}
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

type l2CapabilityManifest struct {
	SchemaVersion                string   `json:"schema_version"`
	StandardSource               string   `json:"standard_source"`
	Package                      string   `json:"package"`
	Layer                        string   `json:"layer"`
	ReleaseLevelTarget           string   `json:"release_level_target"`
	P0Contracts                  []string `json:"p0_contracts"`
	ForbiddenRuntimeDependencies []string `json:"forbidden_runtime_dependencies"`
	ReleaseContract              struct {
		RequiredProfiles    []string `json:"required_profiles"`
		ReleaseAllowed      bool     `json:"release_allowed"`
		FactoryGradeAllowed bool     `json:"factory_grade_allowed"`
		MinScore            int      `json:"min_score"`
	} `json:"release_contract"`
	Capabilities struct {
		Core     []string `json:"core"`
		Optional []string `json:"optional"`
	} `json:"capabilities"`
	Provider struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"provider"`
}

type l2ContractPackRegistry struct {
	SchemaVersion  string           `json:"schema_version"`
	StandardSource string           `json:"standard_source"`
	Packs          []l2ContractPack `json:"packs"`
}

type l2ContractPack struct {
	Name              string               `json:"name"`
	Package           string               `json:"package"`
	Layer             string               `json:"layer"`
	ReleaseLevel      string               `json:"release_level"`
	RequiredProfiles  []string             `json:"required_profiles"`
	RequiredContracts []l2RequiredContract `json:"required_contracts"`
}

type l2RequiredContract struct {
	Name string `json:"name"`
}

type l2Gate struct {
	SchemaVersion       string   `json:"schema_version"`
	StandardSource      string   `json:"standard_source"`
	Package             string   `json:"package"`
	Layer               string   `json:"layer"`
	ReleaseLevelTarget  string   `json:"release_level_target"`
	ReleaseLevelActual  string   `json:"release_level_actual"`
	MinScore            int      `json:"min_score"`
	Score               int      `json:"score"`
	RequiredProfiles    []string `json:"required_profiles"`
	ReleaseAllowed      bool     `json:"release_allowed"`
	FactoryGradeAllowed bool     `json:"factory_grade_allowed"`
	HardFailures        []string `json:"hard_failures"`
	RequiredEvidence    []string `json:"required_evidence"`
}

func readJSONFile[T any](t *testing.T, path string) T {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()

	var value T
	if err := json.NewDecoder(file).Decode(&value); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return value
}

func findContractPack(packs []l2ContractPack, name string) (l2ContractPack, bool) {
	for _, pack := range packs {
		if pack.Name == name {
			return pack, true
		}
	}
	return l2ContractPack{}, false
}

func contractNames(contracts []l2RequiredContract) []string {
	names := make([]string, 0, len(contracts))
	for _, contract := range contracts {
		names = append(names, contract.Name)
	}
	return names
}

func assertSameSet(t *testing.T, label string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
	seen := make(map[string]int, len(got))
	for _, value := range got {
		seen[value]++
	}
	for _, value := range want {
		if seen[value] == 0 {
			t.Fatalf("%s = %v, missing %q from %v", label, got, value, want)
		}
		seen[value]--
	}
	for value, count := range seen {
		if count != 0 {
			t.Fatalf("%s = %v, unexpected %q not in %v", label, got, value, want)
		}
	}
}
