package servicehost

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

func TestServiceSuppliesItsBoundaries(t *testing.T) {
	o := Options{Role: "custom", Port: "8080", Database: true, Concurrency: 4}
	if o.validate() == nil {
		t.Fatal("missing schema accepted")
	}
	o.Readiness = func(context.Context, *pgxpool.Pool) error { return nil }
	if o.validate() == nil {
		t.Fatal("missing policy accepted")
	}
	o.Authorize = func(string, string, string) bool { return false }
	if err := o.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestMemoryBudgetFromConfiguration(t *testing.T) {
	for value, want := range map[string]int64{"96MiB": 96 << 20, "1GiB": 1 << 30} {
		got, e := memoryLimit(value)
		if e != nil || got != want {
			t.Fatal(got, e)
		}
	}
	for _, value := range []string{"-1MiB", "999999999999999GiB", "0MiB", "garbage"} {
		if _, e := memoryLimit(value); e == nil {
			t.Fatal("accepted", value)
		}
	}
}
