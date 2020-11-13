package purplecat

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const gitHubAPIEndPoint = "https://api.github.com"

// GitHubProject represents the repository name in GitHub.
type GitHubProject struct {
	// user name or organization name of the project in GitHub.
	UserName string
	// repository name of the project in GitHub.
	RepositoryName string
	license        *License
}

type githubAPIResponse struct {
	License *License `json:"license"`
}

// NewGitHubProject refers the repository on the GitHub,
// and fetch license information via GitHub API (call FetchLicense function).
func NewGitHubProject(userName, repoName string, context *Context) *GitHubProject {
	ghp := &GitHubProject{UserName: userName, RepositoryName: repoName}
	ghp.FetchLicense(context)
	return ghp
}

// Name returns the name of the receiver project.
func (repo *GitHubProject) Name() string {
	return fmt.Sprintf("%s/%s", repo.UserName, repo.RepositoryName)
}

func (repo *GitHubProject) infoURL() string {
	return fmt.Sprintf("%s/repos/%s", gitHubAPIEndPoint, repo.Name())
}

// Licenses returns the licenses of the receiver project, by fetching license names via GitHub API.
func (repo *GitHubProject) Licenses() []*License {
	if repo.license == nil {
		return []*License{}
	}
	return []*License{repo.license}
}

// Dependencies returns the dependency list of the receiver project, however, this project always returns the empty slice.
func (repo *GitHubProject) Dependencies() []Project {
	return []Project{}
}

// FetchLicense fetch license information via GitHub Rest API.
// If the given context deny the network access, this function returns error.
func (repo *GitHubProject) FetchLicense(context *Context) (*License, error) {
	if !context.Allow(NetworkAccessFlag) {
		return nil, fmt.Errorf("network access denied")
	}
	license, err := fetchLicensesByGitHubAPI(repo)
	if err != nil {
		return nil, err
	}
	repo.license = license
	return license, nil
}

func fetchLicensesByGitHubAPI(repo *GitHubProject) (*License, error) {
	client := resty.New()
	request := client.NewRequest()
	request = request.SetResult(&githubAPIResponse{})
	response, err := request.Get(repo.infoURL())
	if err != nil {
		return nil, err
	}
	if response.StatusCode() == 404 {
		return nil, fmt.Errorf("%s: repository not found", repo.Name())
	}
	json := response.Result().(*githubAPIResponse)
	if json.License == nil {
		return UnknownLicense, nil
	}
	return json.License, nil
}
