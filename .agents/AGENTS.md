# Repository Design Principles

*   **Generics & Type Safety:** Leverage Go 1.23+ generics (`any`, `comparable`, `cmp.Ordered`, constraints). Avoid `interface{}` / `any` boxing or runtime type assertions where type parameters can be used.
*   **Four Core Paradigms:**
    *   `slice`: Eager functional transformations over in-memory slices (`Map`, `Filter`, `FoldLeft`, `FoldRight`, `Reduce`, `Zip`, `UnZip`, `Partition`, `Sum`, `Concat`).
    *   `iterator`: Lazy evaluation based on Go 1.23 standard library `iter.Seq` and `iter.Seq2`. Short-circuiting operations (`AllMatch`, `AnyMatch`, `Limit`, `TakeWhile`) must terminate early without consuming unnecessary elements.
    *   `channel`: Concurrent streaming pipelines using Go channels. Every channel pipeline operation must accept `context.Context` as its first parameter and cleanly handle cancellation without goroutine leaks.
    *   `predicate`: Composable higher-order predicates and nil-safety helpers (`IsNil`, `NotNil`, `Not`, `True`, `False`).
*   **Mathematical & Functional Terminology:** Use standard algebraic and functional programming terminology accurately (e.g. `Monoid` for types supporting associative binary addition with an identity element, `FoldLeft`, `FoldRight`, `FlatMap`).
*   **Concurrency & Goroutine Lifecycle:**
    *   Always select on `<-ctx.Done()` when reading from or writing to channels.
    *   Ensure all spawned goroutines terminate and close destination channels cleanly upon completion or cancellation.
    *   Avoid spawning unbounded per-element goroutines in streaming operators.
*   **Zero-Allocation & Performance Awareness:** Minimize intermediate slice allocations in iterator pipelines; pre-allocate slice capacities where lengths are known in eager slice operations.

# Go Version & Toolchain Consistency

*   **Single Source of Truth (`go.mod`):** The Go version specified in `go.mod` is the canonical version for the entire repository.
*   **CI Workflow Consistency:** GitHub Actions workflows must use `go-version-file: 'go.mod'` with `actions/setup-go@v5` rather than hardcoded version strings.

# Testing Conventions

For all tests in this repository, strictly adhere to the following conventions:

*   **Table-Driven Tests:** Always use table-driven tests for unit testing. This is the idiomatic Go way and makes tests easy to read, extend, and maintain.
*   **Structure:**
    *   Define a `cases` slice of structs.
    *   Each test case struct should have at least a `name string` field.
    *   Iterate through the cases with a `for _, tc := range cases` loop.
    *   Inside the loop, rebind `tc := tc` and use `t.Run(tc.name, func(t *testing.T) { ... })`.
    *   Use `t.Parallel()` inside sub-tests where appropriate.
*   **Short-Circuiting Tests:** Matchers and limiting functions must be tested for early termination (e.g. ensuring supplier/predicate call counts do not exceed the minimum needed).
*   **Cancellation Tests:** Test channel pipelines with cancelled contexts to ensure immediate exit and clean termination without deadlocks or leaked goroutines.

# Examples (`example_*_test.go`)

*   **Runnable Examples:** Significant functions across packages should have runnable examples in `example_*_test.go` files.
*   **Verified Output:** Include `// Output:` comments at the end of example functions so they are validated by `go test ./...` and rendered on `pkg.go.dev`.
*   **Package Naming:** Place example tests in the `<pkg>_test` package to demonstrate public API usage from a consumer's perspective.

# Benchmarks (`*_bench_test.go`)

*   **Location & Naming:** Benchmarks belong in `<pkg>_bench_test.go` and must follow `func Benchmark<Function>(b *testing.B)`.
*   **Allocation Tracking:** Always call `b.ReportAllocs()` in benchmark functions.
*   **Timer Control:** Isolate setup code from measured code using `b.ResetTimer()` and `b.StopTimer()` / `b.StartTimer()`.

# Documentation Roles & Boundaries

*   **`README.md` (User-Facing):** Reserved strictly for public/consumer-facing documentation, quickstart examples, feature overviews, high-level capabilities, installation instructions, and package indices. It should **not** delve into internal agent/contributor rules.
*   **`AGENTS.md` (Architecture & Implementation Guidelines):** The canonical place to encode architecture, design invariants, implementation decisions, performance constraints, code guidelines, testing/benchmarking rules, and agent workflows.
*   **Documentation Synchronization:** Whenever any code change modifies, adds, or removes APIs, behaviors, or package structures, all relevant documentation (**Go doc comments**, **`README.md`**, and **`AGENTS.md`**) MUST be updated in the same change to keep documentation strictly in sync with the codebase.

# Code Style and Quality

*   **Formatting:** Always run `gofmt -s -w .` after making code changes.
*   **Documentation:** All exported packages, types, interfaces, constants, and functions must have comprehensive, up-to-date Go doc comments adhering to standard Go conventions (checked by `revive`).
*   **Linting:** Verify code with `golangci-lint run ./...` before completing tasks.

# Pull Request & Workflow Conventions

*   **Accurate & Detailed PR Descriptions:** PR descriptions must comprehensively capture all introduced features, constructors, public APIs, performance benchmarks, and any toolchain or CI workflow modifications made in the branch.
*   **PR Description Synchronization:** Whenever subsequent commits in a PR branch modify, fix, or expand upon the original changeset, the PR description must be proactively updated via `gh pr edit` to reflect the latest state of the PR.

# Continuous Learning & Self-Correction

*   **Codifying Corrections:** Whenever an agent is corrected by the user regarding an architectural decision, code style, testing pattern, or repository convention, the agent MUST update `AGENTS.md` (or the relevant package documentation) to codify the underlying principle so future agent sessions adhere to it automatically.

# Pre-Completion Verification Checklist

Before concluding any code modification task, agents MUST run and verify the following commands succeed without errors or warnings:
1. `gofmt -s -w .`
2. `go test ./...`
3. `golangci-lint run ./...`
