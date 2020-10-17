package purplecat

import "testing"

func TestGetLicenses(t *testing.T) {
	testdata := []struct {
		user        string
		repo        string
		denyNetwork bool
		isError     bool
		wontLicense string
	}{
		{"tamada", "pochi", false, false, "Apache-2.0"},
		{"tamada", "9rules", false, false, "Apache-2.0"},
		{"tamada", "9rules", true, true, ""},
		{"tamada", "unknown_repo", true, true, ""},
		{"tamada", "unknown_repo", false, true, ""},
		{"tamada", "MonsterPorker", false, true, "unknown"},
		{"tamada", "MonsterPorker", true, true, "unknown"},
	}
	for _, td := range testdata {
		repo := NewGitHubRepository(td.user, td.repo)
		license, err := repo.GetLicense(NewContext(td.denyNetwork, "json", 1))
		if td.isError == (err == nil) {
			t.Errorf("%s/%s: wont error %v, got %v", td.user, td.repo, td.isError, err == nil)
		}
		if license != nil && license.SpdxID != td.wontLicense {
			t.Errorf("%s/%s: wont license %s, got %s", repo.UserName, repo.RepositoryName, td.wontLicense, license.SpdxID)
		}
	}
}
