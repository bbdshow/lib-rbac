package conf

import (
	"github.com/bbdshow/bkit/conf"
	"testing"
)

func TestConfigToFile(t *testing.T) {
	cfg := &Config{}
	if err := conf.UnmarshalDefaultVal(cfg); err != nil {
		t.Fatal(err)
	}
	if err := conf.MarshalToFile(cfg, "../../configs/config.release.toml"); err != nil {
		t.Fatal(err)
	}
}
