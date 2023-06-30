package npm_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/patrickdevivo/go-dep-apis/pkg/npm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var httpClient = http.DefaultClient

func TestMain(m *testing.M) {
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       path.Join("testdata/fixtures"),
		Mode:               recorder.ModeReplayWithNewEpisodes,
		SkipRequestLatency: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	httpClient.Transport = r

	code := m.Run()
	if err = r.Stop(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
func TestGetMeta(t *testing.T) {
	meta, _, err := npm.NewClient(npm.WithHttpClient(httpClient)).GetMeta(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "registry", meta.DBName)
}

func TestGetPackage(t *testing.T) {
	pkg, _, err := npm.NewClient().GetPackage(context.Background(), "jquery")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jquery", pkg.Name)
}
