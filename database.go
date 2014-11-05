package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings" 
	"github.com/russross/meddler"
	. "github.com/drone/drone/shared/model"
)

const (
	commitRecentTemplate = `
SELECT		r.repo_host, r.repo_owner, r.repo_name,
			c.commit_status, c.commit_started, c.commit_finished, c.commit_finished, 
			c.commit_sha, c.commit_branch, c.commit_pr, c.commit_author, 
			c.commit_gravatar, c.commit_timestamp, c.commit_message, c.commit_created, c.commit_updated
FROM		repos r, commits c
WHERE		r.repo_id = c.repo_id
ORDER BY	c.created desc
LIMIT 20
`
	commitTeamStmt = `
SELECT      r.repo_owner, r.repo_owner, r.repo_name,
            c.commit_status, c.commit_started, c.commit_finished, c.commit_finished, 
            c.commit_sha, c.commit_branch, c.commit_pr, c.commit_author, 
            c.commit_gravatar, c.commit_timestamp, c.commit_message, c.commit_created, 
            c.commit_updated
FROM        repos r, commits c
WHERE       r.repo_id = c.repo_id
ORDER BY    c.created DESC
LIMIT 20
`
)

var (
	db *sql.DB
)

func commitRecentStmt() (string, []interface{}) {
	list := strings.Split(*repos, ",")

	repoList := make([]interface{}, len(list))
	for i, v := range list {
		repoList[i] = interface{}(v)
	}

	var params bytes.Buffer
	for i := range repoList {
		if i < len(repoList)-1 {
			params.WriteString(fmt.Sprintf("?, "))
		} else {
			params.WriteString(fmt.Sprintf("?"))
		}
	}

	//return fmt.Sprintf(commitRecentTemplate, params.String()), repoList
	return commitRecentTemplate, repoList
}

func listRepoWallCommits() ([]*CommitRepo, error) {
	var commits []*CommitRepo

	stmt, repoList := commitRecentStmt()

	err := meddler.QueryAll(db, &commits, stmt, repoList...)
	return commits, err
}

func listTeamWallCommits() ([]*CommitRepo, error) {
	var commits []*CommitRepo
	err := meddler.QueryAll(db, &commits, commitTeamStmt, *team)
	return commits, err
}
