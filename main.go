
import (
	"log"
	"regexp"
	"time"

	"github.com/cego/docker-registry-pruner/config"
	"github.com/cego/docker-registry-pruner/registry"
)

func main() {
	log.Print("Configuration:\n" + config.Summary())
	log.Print("Starting...")

	api := registry.NewAPI(config.Host)
	if config.User != "" {
		api.SetCredentials(config.User, config.Pass)
	}

	repositories, err := api.GetRepositories()
	if err != nil {
		log.Fatal(err)
	}

	for _, repository := range repositories {
		if !config.Repos.MatchString(repository) {
			continue
		}

		log.Printf("Inspecting repository %v.", repository)

		tagsIndexedByDigest, err := api.GetTagsIndexedByDigest(repository)
		if err != nil {
			log.Print(err)
			continue
		}

		for digest, tags := range tagsIndexedByDigest {
			if !allMatch(tags, config.Tags) {
				continue
			}

			if config.MinAge > 0 {
				created, err := api.GetManifestCreated(repository, tags[0])
				if err != nil {
					log.Print(err)
					continue
				}

				latestAcceptableCreated := time.Now().Add(-1 * config.MinAge)
				if created.After(latestAcceptableCreated) {
					continue
				}
			}

			if config.Dry {
				log.Printf("Would have deleted %s manifest with digest %s and tags %v.", repository, digest, tags)
			} else {
				log.Printf("Deleting %s manifest with digest %s and tags %v.", repository, digest, tags)
				err = api.DeleteManifest(repository, digest)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}

	log.Printf("Done.")
}

// allMatch returns true if all strings in the given slice matches the given expr.
func allMatch(slice []string, expr *regexp.Regexp) bool {
	for _, s := range slice {
		if !expr.MatchString(s) {
			return false
		}
	}

	return true
}
