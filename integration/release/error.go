// +build k8srequired

package release

import "github.com/giantswarm/microerror"

var releaseStatusNotMatchingError = microerror.New("release status not matching")

// IsReleaseStatusNotMatching asserts releaseStatusNotMatchingError
func IsReleaseStatusNotMatching(err error) bool {
	return microerror.Cause(err) == releaseStatusNotMatchingError
}
