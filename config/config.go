package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/namsral/flag"
)

// Public variables used by other code requiring configuration parameters.
// See the flags in init() for a note on each one.
var (
	Host   string
	User   string
	Pass   string
	Repos  *regexp.Regexp
	Tags   *regexp.Regexp
	MinAge time.Duration
	Dry    bool
)

func init() {
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "DOCKER_REGISTRY_PRUNER", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Printf("Flags accepted by %s:\n\n", os.Args[0])
		fs.PrintDefaults()
		fmt.Print("\nNote that each flag can also be given as an environment variable with the prefix \"DOCKER_REGISTRY_PRUNER_\".\n")
		fmt.Printf("\nFor example, DOCKER_REGISTRY_PRUNER_USER=bert %s\n", os.Args[0])
	}

	fs.StringVar(&Host, "host", "http://localhost:5000", "Address of the registry.")
	fs.StringVar(&User, "user", "", "Username for the registry API. Omit to skip authentication.")
	fs.StringVar(&Pass, "pass", "", "Password for the registry API.")
	fs.DurationVar(&MinAge, "minage", 24*30*time.Hour, "Minimum age of manifests to delete. Manifests newer than this will not be deleted. This flag accepts any input valid for the time.ParseDuration function of go (eg. '4h30m', '23s'), see https://golang.org/pkg/time/#ParseDuration.")
	fs.BoolVar(&Dry, "dry", false, "Run in dry mode where nothing is deleted.")

	repos := fs.String("repos", ".*", "Regular expression matching all repositories to delete from.")
	tags := fs.String("tags", ".*", "Regular expression matching all tags to delete.")

	fs.Parse(os.Args[1:])

	Repos = regexp.MustCompile(*repos)
	Tags = regexp.MustCompile(*tags)
}

// Summary returns a human-readable description of the effect of read configuration parameters.
func Summary() string {
	var result string

	result += fmt.Sprintf("Connecting to registry at '%s'.\n", Host)
	if User != "" {
		result += fmt.Sprintf("Authenticating as user '%s'.\n", User)
	} else {
		result += "Not authenticating.\n"
	}
	result += fmt.Sprintf("Manifest from repositories matching '%s' where all tags of the manifest match '%s' are up for deletion.\n", Repos, Tags)
	if MinAge > 0 {
		result += fmt.Sprintf("Manifests created with the last %s will not be deleted.\n", MinAge)
	}
	if Dry {
		result += "This is a dry-run, so nothing will be deleted.\n"
	}

	return result
}
