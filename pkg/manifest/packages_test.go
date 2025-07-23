package manifest

import (
	"github.com/Masterminds/semver/v3"
	"testing"
)

func TestReadPackages(t *testing.T) {
	pf, err := ReadPackages("testdata/packages.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(pf.Packages) != 3 || pf.Packages[0].Name != "packet-1" {
		t.Fatalf("unexpected %+v", pf.Packages)
	}
}

func TestConstraint(t *testing.T) {
	c, err := ParseConstraint(">=1.10")
	if err != nil {
		t.Fatal(err)
	}

	ok := c.Check(semver.MustParse("1.12.0"))
	if !ok {
		t.Fatalf("should match 1.12")
	}
}
