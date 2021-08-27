package changelog

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	releaseNoteRegexp = regexp.MustCompile("(?s)```release-note(.+?)```")
	typologyRegexp    = regexp.MustCompile(`(?m)(.+?)(\((.+)\))?: ?(.*)`)
)

const defaultGitHubBaseURI = "https://github.com"

// ReleaseNote ...
type ReleaseNote struct {
	Typology    string
	Scope       string
	Description string
	URI         string
	Num         int
	Author      string
	AuthorURL   string
	MergeDay    string
}

type Client struct {
	c *github.Client
}

// NewClient ...
func NewClient(token string) *Client {
	client := github.NewClient(nil)

	// Eventually create an authenticated client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	}

	return &Client{
		c: client,
	}
}

// Get returns the list of release notes found for the given parameters.
func (c *Client) Get(org, repo, branch string, lastPrNr int) ([]ReleaseNote, error) {
	ctx := context.Background()
	listingOpts := &github.PullRequestListOptions{
		State:     "closed",
		Base:      branch,
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 1000,
		},
	}
	prs, _, err := c.c.PullRequests.List(ctx, org, repo, listingOpts)
	if _, ok := err.(*github.RateLimitError); ok {
		return nil, fmt.Errorf("hit rate limiting")
	}
	if err != nil {
		return nil, err
	}

	releaseNotes := []ReleaseNote{}
	for _, p := range prs {
		num := p.GetNumber()
		if num <= lastPrNr {
			break
		}

		isMerged, _, err := c.c.PullRequests.IsMerged(ctx, org, repo, num)
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("hit rate limiting")
		}
		if err != nil {
			return nil, fmt.Errorf("error detecting if pr %d is merged or not", num)
		}
		if !isMerged {
			// It means PR has been closed but not merged in
			continue
		}

		res := releaseNoteRegexp.FindStringSubmatch(p.GetBody())
		note := "NONE"
		if len(res) > 0 {
			note = strings.TrimSpace(res[1])
		}
		if note == "NONE" || note == "none" {
			rn := ReleaseNote{
				Typology:    "none",
				Scope:       "",
				Description: p.GetTitle(),
				URI:         fmt.Sprintf("%s/%s/%s/pull/%d", defaultGitHubBaseURI, org, repo, num),
				Num:         num,
				Author:      fmt.Sprintf("@%s", p.GetUser().GetLogin()),
				AuthorURL:   p.GetUser().GetHTMLURL(),
				MergeDay:    p.GetMergedAt().Format("2000-01-01"),
			}
			releaseNotes = append(releaseNotes, rn)
			continue
		}
		notes := strings.Split(note, "\n")
		for _, n := range notes {
			n = strings.Trim(n, "\r")
			matches := typologyRegexp.FindStringSubmatch(n)
			if len(matches) < 5 {
				return nil, fmt.Errorf("error extracting type from release note, pr: %d", num)
			}

			rn := ReleaseNote{
				Typology:    matches[1],
				Scope:       matches[3],
				Description: n,
				URI:         fmt.Sprintf("%s/%s/%s/pull/%d", defaultGitHubBaseURI, org, repo, num),
				Num:         num,
				Author:      fmt.Sprintf("@%s", p.GetUser().GetLogin()),
				AuthorURL:   p.GetUser().GetHTMLURL(),
				MergeDay:    p.GetMergedAt().Format("2000-01-01"),
			}
			releaseNotes = append(releaseNotes, rn)
		}
	}

	return releaseNotes, nil
}
