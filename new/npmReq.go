package new

import (
	"encoding/json"
	"net/http"
	"sort"
)

type NPMRes struct {
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

	var resJSON NPMRes
	json.NewDecoder(response.Body).Decode(&resJSON)
	availVers := filterVersions(resJSON.Versions, resJSON.DistTags)

	return availVers

}

func filterVersions(versions VersionsList, tags TagList) TagList {

	tags.Stables = make([]string, 0)

	for k := range versions {
		if len(k) != 5 || k == tags.Latest {
			continue
		}
		tags.Stables = append(tags.Stables, k)
	}
	sort.Slice(tags.Stables, func(i, j int) bool {
		return tags.Stables[i][2] < tags.Stables[j][2]
	})

	return tags
}
