package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gojektech/heimdall/v6/httpclient"
)

const (
	patch = "patch"
	minor = "minor"
	major = "major"

	reqTimeout = 2000 * time.Millisecond
	githubAPIBaseURI = "https://api.github.com"
)

func main() {

	versionFlgPtr :=
		flag.String("version", patch, "what kind of release to make (major / minor / patch)")
	flag.Parse()

	repoName := flag.Args()[0]

	if validateVersionInput(*versionFlgPtr) {
		rel, tagErr := createOrUpdateTag(repoName, *versionFlgPtr)
		if tagErr != nil {
			panic(tagErr)
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("publish new release? (%s) [y/N]: ", rel.String())
			scanner.Scan()
			confirm := strings.ToLower(scanner.Text())

			if confirmRelease(confirm) {
				pubErr := publishRelease(repoName, rel)
				if pubErr != nil {
					panic(pubErr)
				} else {
					fmt.Printf("successfully published release: %s\n", rel.String())
				}
			} else {
				fmt.Println("Aborting publish of release")
			}
		}
	} else {
		fmt.Println("version must be major, minor or patch")
	}

}

func confirmRelease(in string) bool {
	return in == "y"
}

func validateVersionInput(version string) bool {
	return version == patch || version == minor || version == major
}

func createOrUpdateTag(repo, version string) (rel release, err error) {
	client := createHTTPClient()
	uri := releaseURI(repo)
	headers := createHeaders()

	resp, err := client.Get(uri, headers)
	if err != nil || resp.StatusCode != 200 {
		return rel, fmt.Errorf("Couldn't get the current release version")
	}

	defer resp.Body.Close()
	target := []releseGetResponse{} // Github returns an array
	err = json.NewDecoder(resp.Body).Decode(&target)
	if len(target) > 0 {
		// Get a release struct from the current tag
		releaseFromString(target[0].TagName, &rel)
	}
	fmt.Printf("current version: %s\n", rel)
	incrementRelease(version, &rel)

	fmt.Printf("new version: %s\n", rel)
	return rel, err
}

func publishRelease(repo string, rel release) error {
	client := createHTTPClient()
	uri := releaseURI(repo)
	headers := createHeaders()
	relReq := releasePostRequest{
		TagName: rel.String(),
		Name: rel.String(),
		Draft: false,
	}

	js, err := json.Marshal(relReq)
	if err != nil {
		return err
	}

	resp, err := client.Post(
		uri,
		bytes.NewReader(js),
		headers,
	)

	if err != nil || resp.StatusCode != 201 {
		return fmt.Errorf("Couldn't publish new release version")
	}

	return nil
}

func createHTTPClient() *httpclient.Client {
	return httpclient.NewClient(httpclient.WithHTTPTimeout(reqTimeout))
}

func releaseURI(repo string) string {
	return githubAPIBaseURI + "/repos/" + repo + "/releases"
}

func createHeaders() http.Header {
	token := "token " + os.Getenv("GITHUB_TOKEN")
	var headers http.Header = make(map[string][]string)
	headers.Add("Authorization", token)
	headers.Add("Accept", "application/vnd.github.v3+json")
	return headers
}

func incrementRelease(version string, rel *release) {
	if version == patch {
		rel.IncPatch()
	} else if version == minor {
		rel.IncMinor()
	} else {
		rel.IncMajor()
	}
}
