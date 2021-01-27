package main

import (
	"bytes"
	"github.com/mattfarina/hypper/internal/test"
	"helm.sh/helm/v3/pkg/chartutil"
	kubefake "helm.sh/helm/v3/pkg/kube/fake"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"helm.sh/helm/v3/pkg/time"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mattfarina/hypper/pkg/cli"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	helmAction "helm.sh/helm/v3/pkg/action"
)

func testTimestamper() time.Time { return time.Unix(242085845, 0).UTC() }

func init() {
	helmAction.Timestamper = testTimestamper
}

func runTestActionCmd(t *testing.T, tests []cmdTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetEnv()()

			store := storageFixture()
			for _, rel := range tt.rels {
				_ = store.Create(rel)
			}
			_, out, err := executeActionCommandC(store, tt.cmd)
			if (err != nil) != tt.wantError {
				t.Errorf("expected error, got '%v'", err)
			}
			if tt.golden != "" {
				test.AssertGoldenString(t, out, tt.golden)
			}
		})
	}
}

func storageFixture() *storage.Storage {
	return storage.Init(driver.NewMemory())
}

func executeActionCommandC(store *storage.Storage, cmd string) (*cobra.Command, string, error) {
	return executeActionCommandStdinC(store, nil, cmd)
}

func executeActionCommandStdinC(store *storage.Storage, in *os.File, cmd string) (*cobra.Command, string, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)

	actionConfig := &helmAction.Configuration{
		Releases:     store,
		KubeClient:   &kubefake.PrintingKubeClient{Out: ioutil.Discard},
		Capabilities: chartutil.DefaultCapabilities,
		Log:          func(format string, v ...interface{}) {},
	}

	root, err := newRootCmd(actionConfig, logger, args)
	if err != nil {
		return nil, "", err
	}

	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	oldStdin := os.Stdin
	if in != nil {
		root.SetIn(in)
		os.Stdin = in
	}

	if mem, ok := store.Driver.(*driver.Memory); ok {
		mem.SetNamespace(settings.HelmSettings.Namespace())
	}
	c, err := root.ExecuteC()
	result := buf.String()
	os.Stdin = oldStdin

	return c, result, err
}

// cmdTestCase describes a test case that works with releases.
type cmdTestCase struct {
	name      string
	cmd       string
	golden    string
	wantError bool
	// Rels are the available releases at the start of the test.
	rels []*release.Release
	// Number of repeats (in case a feature was previously flaky and the test checks
	// it's now stably producing identical results). 0 means test is run exactly once.
	// TODO(itxaka): Disabled for now, we are not using it but we may want to keep it 1:1 with helm?
	//repeat int
}

func resetEnv() func() {
	origEnv := os.Environ()
	return func() {
		os.Clearenv()
		for _, pair := range origEnv {
			kv := strings.SplitN(pair, "=", 2)
			os.Setenv(kv[0], kv[1])
		}
		settings = cli.New(logger)
	}
}

func executeCommandStdinC(cmd string) (*cobra.Command, string, error) {

	buf := new(bytes.Buffer)
	args, err := shellwords.Parse(cmd)
	actionConfig := new(helmAction.Configuration)

	if err != nil {
		return nil, "", err
	}

	root, err := newRootCmd(actionConfig, logger, args)
	if err != nil {
		return nil, "", err
	}

	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	oldStdin := os.Stdin

	c, err := root.ExecuteC()
	result := buf.String()
	os.Stdin = oldStdin

	return c, result, err
}
