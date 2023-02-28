// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package update

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

const (
	maxPerPage = 100 // https://docs.github.com/en/rest/reference/repos
	timeout    = 15 * time.Second
)

func gitHubTokenFromEnv() string {
	token := os.Getenv("BLDR_GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	return token
}

type gitHub struct {
	c *github.Client
}

func newGitHub(token string) *gitHub {
	c := new(http.Client)

	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		c = oauth2.NewClient(context.Background(), src)
	}

	c.Timeout = timeout

	return &gitHub{
		c: github.NewClient(c),
	}
}

// Latest returns information about available update.
func (g *gitHub) Latest(ctx context.Context, source string) (*LatestInfo, error) {
	sourceURL, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	if sourceURL.Host != "github.com" {
		panic(fmt.Sprintf("unexpected host %q", sourceURL.Host))
	}

	parts := strings.Split(sourceURL.Path, "/")
	owner, repo := parts[1], parts[2]

	v, err := extractVersion(source)
	if err != nil {
		return nil, err
	}

	considerPrereleases := v.Prerelease() != ""

	releases, err := g.getReleases(ctx, owner, repo)
	if err != nil {
		return nil, g.wrapGitHubError(err)
	}

	if len(releases) != 0 {
		return g.findLatestRelease(releases, sourceURL, considerPrereleases)
	}

	tags, err := g.getTags(ctx, owner, repo)
	if err != nil {
		return nil, g.wrapGitHubError(err)
	}

	return g.findLatestTag(ctx, tags, sourceURL, considerPrereleases)
}

// findLatestRelease returns information about latest released version.
func (g *gitHub) findLatestRelease(releases []*github.RepositoryRelease, sourceURL *url.URL, considerPrereleases bool) (*LatestInfo, error) {
	parts := strings.Split(sourceURL.Path, "/")
	owner, repo := parts[1], parts[2]

	var newest *github.RepositoryRelease

	// find newest release
	for _, release := range releases {
		if release.GetPrerelease() && !considerPrereleases {
			continue
		}

		if newest == nil || newest.CreatedAt.Before(release.CreatedAt.Time) {
			newest = release
		}
	}

	if newest == nil {
		return nil, fmt.Errorf("no release found")
	}

	res := &LatestInfo{
		BaseURL: fmt.Sprintf("https://github.com/%s/%s/releases/", owner, repo),
	}

	source := sourceURL.String()

	// update is available if the newest release doesn't have source in their assets download URLs
	for _, asset := range newest.Assets {
		if asset.GetBrowserDownloadURL() == source {
			res.LatestURL = source

			return res, nil
		}
	}

	// check default .tag.gz URL
	latestTarGz := g.getTagTarGZ(owner, repo, newest.GetTagName())
	if latestTarGz == source {
		res.LatestURL = source

		return res, nil
	}

	// we don't know correct asset if there are any
	if len(newest.Assets) == 0 {
		res.LatestURL = latestTarGz
	}

	res.HasUpdate = true

	return res, nil
}

// findLatestTag returns information about latest tagged version.
func (g *gitHub) findLatestTag(ctx context.Context, tags []*github.RepositoryTag, sourceURL *url.URL, considerPrereleases bool) (*LatestInfo, error) {
	parts := strings.Split(sourceURL.Path, "/")
	owner, repo := parts[1], parts[2]

	var (
		newest     *github.RepositoryTag
		newestDate time.Time
	)

	// find newest tag
	for _, tag := range tags {
		v, err := extractVersion(tag.GetName())
		if err != nil {
			return nil, err
		}

		if v.Prerelease() != "" && !considerPrereleases {
			continue
		}

		tagDate, err := g.getCommitTime(ctx, owner, repo, tag.GetCommit().GetSHA())
		if err != nil {
			return nil, err
		}

		if newest == nil || newestDate.Before(tagDate) {
			newest = tag
			newestDate = tagDate
		}
	}

	if newest == nil {
		return nil, fmt.Errorf("no tag found")
	}

	res := &LatestInfo{
		BaseURL:   fmt.Sprintf("https://github.com/%s/%s/releases/", owner, repo),
		LatestURL: g.getTagTarGZ(owner, repo, newest.GetName()),
	}

	// update is available if the newest tag doesn't have the same tarball URL
	res.HasUpdate = res.LatestURL != sourceURL.String()

	return res, nil
}

// getReleases returns all releases.
func (g *gitHub) getReleases(ctx context.Context, owner, repo string) ([]*github.RepositoryRelease, error) {
	opts := &github.ListOptions{
		PerPage: maxPerPage,
	}

	var res []*github.RepositoryRelease

	for {
		page, resp, err := g.c.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		res = append(res, page...)

		if resp.NextPage == 0 {
			return res, nil
		}

		opts.Page = resp.NextPage
	}
}

// getTags returns all tags.
func (g *gitHub) getTags(ctx context.Context, owner, repo string) ([]*github.RepositoryTag, error) {
	opts := &github.ListOptions{
		PerPage: maxPerPage,
	}

	var res []*github.RepositoryTag

	for {
		page, resp, err := g.c.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		res = append(res, page...)

		if resp.NextPage == 0 {
			return res, nil
		}

		opts.Page = resp.NextPage
	}
}

// getCommitTime returns commit's time.
func (g *gitHub) getCommitTime(ctx context.Context, owner, repo, sha string) (time.Time, error) {
	commit, _, err := g.c.Repositories.GetCommit(ctx, owner, repo, sha, &github.ListOptions{})
	if err != nil {
		return time.Time{}, err
	}

	t := commit.GetCommit().GetCommitter().GetDate()
	if t.IsZero() {
		return time.Time{}, fmt.Errorf("no commit date")
	}

	return t.Time, nil
}

// getTagTarGZ returns .tar.gz URL.
// API's GetTarballURL is not good enough.
func (g *gitHub) getTagTarGZ(owner, repo, name string) string {
	return fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.tar.gz", owner, repo, name)
}

func (g *gitHub) wrapGitHubError(err error) error {
	if err == nil {
		return nil
	}

	var ghe *github.RateLimitError

	if errors.As(err, &ghe) {
		err = fmt.Errorf("%w\nSet `BLDR_GITHUB_TOKEN` or `GITHUB_TOKEN` environment variable", err)
	}

	return err
}
