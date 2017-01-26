package config

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/git-lfs/git-lfs/lfsapi"
)

var (
	GitCommit   string
<<<<<<< HEAD
<<<<<<< HEAD
=======
	Version     = "1.5.3"
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
=======
	Version     = "1.5.5"
>>>>>>> refs/remotes/origin/release-1.5
	VersionDesc string
)

const (
	Version = "1.5.0"
)

func init() {
	gitCommit := ""
	if len(GitCommit) > 0 {
		gitCommit = "; git " + GitCommit
	}
	VersionDesc = fmt.Sprintf("git-lfs/%s (GitHub; %s %s; go %s%s)",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		strings.Replace(runtime.Version(), "go", "", 1),
		gitCommit,
	)

	lfsapi.UserAgent = VersionDesc
}
