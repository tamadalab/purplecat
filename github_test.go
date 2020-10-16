package purplecat

import "testing"

func TestGetLicenses(t *testing.T) {
	testdata := []struct {
		user        string
		repo        string
		isError     bool
		wontLicense string
	}{
		{"tamada", "pochi", false, "Apache-2.0"},
		{"tamada", "9rules", false, "Apache-2.0"},
		{"tamada", "unknown_repo", true, ""},
		{"tamada", "tamada.github.io", false, "unknown"},
	}
	for _, td := range testdata {
		repo := NewGitHubRepository(td.user, td.repo)
		license, err := repo.GetLicense(NewContext(false, "json", 1))
		if td.isError == (err == nil) {
			t.Errorf("wont error: %v, got %v", td.isError, err == nil)
		}
		if license != nil && license.SpdxId != td.wontLicense {
			t.Errorf("%s/%s: wont license %s, got %s", repo.UserName, repo.RepositoryName, td.wontLicense, license.SpdxId)
		}
	}
}
