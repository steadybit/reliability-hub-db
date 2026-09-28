// Command sync regenerates the reliability-hub-db YAML files for one
// extension from its --describe JSON output.
//
// Usage:
//
//	sync --extension <id> [--describe-file <path>] [--image <url>] [--dry-run]
//
// One of --describe-file or --image is required. With --describe-file the
// caller has already captured a --describe dump (offline mode); with --image
// the tool runs `docker run --rm <image> --describe` to produce one.
//
// The extension's non-binary metadata (GitHub/GHCR coords, license, tags,
// releaseDate, etc.) is read from sync.yml at the repo root.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/steadybit/reliability-hub-db/internal/describe"
	"github.com/steadybit/reliability-hub-db/internal/source"
	"github.com/steadybit/reliability-hub-db/internal/syncconfig"
	"github.com/steadybit/reliability-hub-db/internal/syncer"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "sync:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var (
		extension    = fs.String("extension", "", "extension ID (required)")
		describeFile = fs.String("describe-file", "", "path to a captured --describe JSON file (offline mode)")
		image        = fs.String("image", "", "GHCR image to run for --describe (online mode; overrides image field in sync.yml)")
		syncCfg      = fs.String("sync-config", "sync.yml", "path to sync.yml registry")
		root         = fs.String("root", ".", "path to the reliability-hub-db checkout")
		dryRun       = fs.Bool("dry-run", false, "print what would change without writing")
		keepOrphans  = fs.Bool("keep-orphans", false, "report orphans but do not delete them")
		minRoutes    = fs.Int("min-routes", 2, "abort if --describe has fewer than this many routes (sanity guard)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *extension == "" {
		fs.Usage()
		return fmt.Errorf("--extension is required")
	}
	if *describeFile == "" && *image == "" {
		return fmt.Errorf("one of --describe-file or --image is required")
	}

	cfgPath, err := filepath.Abs(*syncCfg)
	if err != nil {
		return err
	}
	cfg, err := syncconfig.Load(cfgPath)
	if err != nil {
		return err
	}
	ext := cfg.Lookup(*extension)
	if ext == nil {
		return fmt.Errorf("extension %q not found in %s", *extension, cfgPath)
	}

	var payload []byte
	if *describeFile != "" {
		payload, err = source.FromFile(*describeFile)
	} else {
		imageRef := *image
		if imageRef == "" {
			imageRef = ext.Image
		}
		if imageRef == "" {
			return fmt.Errorf("--image not set and sync.yml has no image for %s", ext.ID)
		}
		fmt.Fprintf(stdout, "running docker for %s: %s\n", ext.ID, imageRef)
		payload, err = source.FromImage(imageRef)
	}
	if err != nil {
		return err
	}

	out, err := describe.Parse(payload)
	if err != nil {
		return err
	}

	rootAbs, err := filepath.Abs(*root)
	if err != nil {
		return err
	}

	report, err := syncer.Sync(out, *ext, syncer.Options{
		Root:        rootAbs,
		DryRun:      *dryRun,
		KeepOrphans: *keepOrphans,
		MinRoutes:   *minRoutes,
	})
	if err != nil {
		return err
	}

	printReport(stdout, report, *dryRun)
	return nil
}

func printReport(w io.Writer, r *syncer.Report, dryRun bool) {
	prefix := ""
	if dryRun {
		prefix = "[dry-run] "
	}
	fmt.Fprintf(w, "%sextension:    %s\n", prefix, r.ExtensionID)
	fmt.Fprintf(w, "%sactions:      %d projected\n", prefix, len(r.ActionsWritten))
	fmt.Fprintf(w, "%stargetTypes:  %d projected\n", prefix, len(r.TargetTypesWritten))
	fmt.Fprintf(w, "%sadvice:       %d projected\n", prefix, len(r.AdviceWritten))
	fmt.Fprintf(w, "%smdx updated:  %d %v\n", prefix, len(r.MDXUpdated), r.MDXUpdated)
	fmt.Fprintf(w, "%smdx unchanged:%d\n", prefix, len(r.MDXUnchanged))
	if len(r.MDXNoHeading) > 0 {
		fmt.Fprintf(w, "%smdx no-heading:    %v\n", prefix, r.MDXNoHeading)
	}
	if len(r.MDXNoTable) > 0 {
		fmt.Fprintf(w, "%smdx no-table:      %v\n", prefix, r.MDXNoTable)
	}
	if len(r.MDXMissingFile) > 0 {
		fmt.Fprintf(w, "%smdx missing-file:  %v\n", prefix, r.MDXMissingFile)
	}
	if len(r.OrphanActions) > 0 {
		verb := "deleted"
		if !r.OrphansDeleted {
			verb = "would delete"
		}
		fmt.Fprintf(w, "%s%s orphan actions:      %v\n", prefix, verb, r.OrphanActions)
	}
	if len(r.OrphanTargetTypes) > 0 {
		verb := "deleted"
		if !r.OrphansDeleted {
			verb = "would delete"
		}
		fmt.Fprintf(w, "%s%s orphan targetTypes:  %v\n", prefix, verb, r.OrphanTargetTypes)
	}
	if len(r.OrphanAdvice) > 0 {
		verb := "deleted"
		if !r.OrphansDeleted {
			verb = "would delete"
		}
		fmt.Fprintf(w, "%s%s orphan advice:       %v\n", prefix, verb, r.OrphanAdvice)
	}
}
