package vagrant

import (
	"bytes"
	"compress/flate"
	"github.com/mitchellh/packer/packer"
	"strings"
	"testing"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func testPP(t *testing.T) *PostProcessor {
	var p PostProcessor
	if err := p.Configure(testConfig()); err != nil {
		t.Fatalf("err: %s", err)
	}

	return &p
}

func testUi() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestPostProcessorPrepare_compressionLevel(t *testing.T) {
	var p PostProcessor

	// Default
	c := testConfig()
	delete(c, "compression_level")
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	if p.config.CompressionLevel != flate.DefaultCompression {
		t.Fatalf("bad: %#v", p.config.CompressionLevel)
	}

	// Set
	c = testConfig()
	c["compression_level"] = 7
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	if p.config.CompressionLevel != 7 {
		t.Fatalf("bad: %#v", p.config.CompressionLevel)
	}
}

func TestPostProcessorPrepare_outputPath(t *testing.T) {
	var p PostProcessor

	// Default
	c := testConfig()
	delete(c, "output")
	err := p.Configure(c)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Bad template
	c["output"] = "bad {{{{.Template}}}}"
	err = p.Configure(c)
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestPostProcessorPostProcess_badId(t *testing.T) {
	artifact := &packer.MockArtifact{
		BuilderIdValue: "invalid.packer",
	}

	_, _, err := testPP(t).PostProcess(testUi(), artifact)
	if !strings.Contains(err.Error(), "artifact type") {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderForName(t *testing.T) {
	if v, ok := providerForName("virtualbox").(*VBoxProvider); !ok {
		t.Fatalf("bad: %#v", v)
	}

	if providerForName("nope") != nil {
		t.Fatal("should be nil if bad provider")
	}
}
