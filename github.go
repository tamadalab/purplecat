package purplecat

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const GITHUB_API_ENDPOINT = "https://api.github.com"

type GitHubRepository struct {
	UserName       string
	RepositoryName string
}

type githubApiResponse struct {
	License *License `json:"license"`
}

func NewGitHubRepository(userName, repoName string) *GitHubRepository {
	return &GitHubRepository{UserName: userName, RepositoryName: repoName}
}

func (repo *GitHubRepository) infoUrl() string {
	return fmt.Sprintf("%s/repos/%s/%s", GITHUB_API_ENDPOINT, repo.UserName, repo.RepositoryName)
}

func (repo *GitHubRepository) GetLicense(context *Context) (*License, error) {
	if !context.Allow(NETWORK_ACCESS) {
		return nil, fmt.Errorf("network access denide")
	}
	return findLicensesByGitHubApi(repo)
}

func findLicensesByGitHubApi(repo *GitHubRepository) (*License, error) {
	client := resty.New()
	request := client.NewRequest()
	request = request.SetResult(&githubApiResponse{})
	response, err := request.Get(repo.infoUrl())
	if err != nil {
		return nil, err
	}
	if response.StatusCode() == 404 {
		return nil, fmt.Errorf("%s/%s: repository not found", repo.UserName, repo.RepositoryName)
	}
	json := response.Result().(*githubApiResponse)
	if json.License == nil {
		return UNKNOWN_LICENSE, nil
	}
	return json.License, nil
}
