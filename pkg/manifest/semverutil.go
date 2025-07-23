package manifest

import "github.com/Masterminds/semver/v3"

func ParseConstraint(expr string) (*semver.Constraints, error) {
	if expr == "" {
		return semver.NewConstraint("*")
	}
	return semver.NewConstraint(expr)
}
