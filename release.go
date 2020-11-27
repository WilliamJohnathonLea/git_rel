package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type releseGetResponse struct {
	TagName string `json:"tag_name"`
}

type releasePostRequest struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Draft   bool   `json:"draft"`
}

type release struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (r release) String() string {
	return fmt.Sprintf("v%d.%d.%d", r.Major, r.Minor, r.Patch)
}

func (r *release) IncMajor() {
	r.Major = r.Major + 1
	r.Minor = 0
	r.Patch = 0
}

func (r *release) IncMinor() {
	r.Minor = r.Minor + 1
	r.Patch = 0
}

func (r *release) IncPatch() {
	r.Patch = r.Patch + 1
}

func releaseFromString(in string, out *release) (err error) {
	r := regexp.MustCompile(`v?(?P<major>\d+).(?P<minor>\d+).(?P<patch>\d+)`)
	result := r.FindStringSubmatch(in)

	out.Major, err = strconv.ParseUint(result[1], 10, 64)
	out.Minor, err = strconv.ParseUint(result[2], 10, 64)
	out.Patch, err = strconv.ParseUint(result[3], 10, 64)
	return err
}
