package new

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type _NPMRes struct {
	Versions VersionsList `json:"versions"`
	DistTags TagList      `json:"dist-tags"`
}

type VersionsList map[string]struct {
	Version string `json:"version"`
}

type TagList struct {
	Latest  string `json:"latest"`
	Beta    string `json:"beta"`
	Preview string `json:"preview"`
	RC      string `json:"rc"`
	Stables []string
}

func getAvailableVersions() TagList {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://registry.npmjs.com/@minecraft/server", nil)
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var resJSON _NPMRes
	json.NewDecoder(response.Body).Decode(&resJSON)
	availVers := filterVersions(resJSON.Versions, resJSON.DistTags)

	return availVers

}

func filterVersions(versions VersionsList, tags TagList) TagList {

	tags.Stables = make([]string, 0)

	currentVer := []string{"0", "0", "0"}
	for k := range versions {
		if len(k) == 5 {
			tags.Stables = append(tags.Stables, k)
			continue
		}
		if strings.Contains(k, "internal") || regexp.MustCompile("\\d{4}").Match([]byte(k)) {
			continue
		}

		var tag = parseType(k)

		switch tag.DistTag {
		case "beta-stable":
			mcVer := strings.Split(tag.MinecraftVersion, ".")
			if tag.Version[1] <= currentVer[1] ||
				mcVer[2] <= currentVer[2] {
				continue
			}
			fmt.Println(k)
			currentVer = mcVer
			tags.Latest = k

		default:
			continue
		}

	}
	sort.Slice(tags.Stables, func(i, j int) bool {
		return tags.Stables[i][2] < tags.Stables[j][2]
	})

	return tags
}

func parseType(version string) (tag struct {
	Version          []string
	DistTag          string
	MinecraftVersion string
}) {
	split := strings.SplitN(version, "-", 2)
	tag.Version = strings.Split(split[0], ".")
	metadata := split[1]

	switch metadata[0] {
	case 'b':
		if strings.Contains(metadata, "stable") {
			tag.DistTag = "beta-stable"
		} else {
			tag.DistTag = "beta-preview"
		}
		tag.MinecraftVersion = metadata[5:]

	case 'r':
		tag.DistTag = "rc"
		tag.MinecraftVersion = metadata[3:]
	}

	return tag
}
