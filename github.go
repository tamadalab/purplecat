package purplecat

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const gitHubAPIEndPoint = "https://api.github.com"

/*
 * GitHubRepository represents the repository name in the GitHub.
 */
type GitHubRepository struct {
	UserName       string
	RepositoryName string
}

type githubAPIResponse struct {
	License *License `json:"license"`
}

func NewGitHubRepository(userName, repoName string) *GitHubRepository {
	return &GitHubRepository{UserName: userName, RepositoryName: repoName}
}

func (repo *GitHubRepository) infoURL() string {
	return fmt.Sprintf("%s/repos/%s/%s", gitHubAPIEndPoint, repo.UserName, repo.RepositoryName)
}

func (repo *GitHubRepository) GetLicense(context *Context) (*License, error) {
	if !context.Allow(NetworkAccessFlag) {
		return nil, fmt.Errorf("network access denide")
	}
	return findLicensesByGitHubAPI(repo)
}

func findLicensesByGitHubAPI(repo *GitHubRepository) (*License, error) {
	client := resty.New()
	request := client.NewRequest()
	request = request.SetResult(&githubAPIResponse{})
	response, err := request.Get(repo.infoURL())
	if err != nil {
		return nil, err
	}
	if response.StatusCode() == 404 {
		return nil, fmt.Errorf("%s/%s: repository not found", repo.UserName, repo.RepositoryName)
	}
	json := response.Result().(*githubAPIResponse)
	if json.License == nil {
		return UnknownLicense, nil
	}
	return json.License, nil
}
