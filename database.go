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
SELECT		r.slug, r.host, r.owner, r.name,
			c.status, c.started, c.finished, c.duration, 
			c.hash, c.branch, c.pull_request, c.author, 
			c.gravatar, c.timestamp, c.message, c.created, c.updated
FROM		repos r, commits c
WHERE		r.id = c.repo_id
AND			r.slug IN (%s)
ORDER BY	c.created desc
LIMIT 20
`
	commitTeamStmt = `
SELECT      r.slug, r.host, r.owner, r.name,
            c.status, c.started, c.finished, c.duration, 
            c.hash, c.branch, c.pull_request, c.author, 
            c.gravatar, c.timestamp, c.message, c.created, 
            c.updated
FROM        repos r, commits c
WHERE       r.id = c.repo_id
AND         r.team_id IN (SELECT	id 
						  FROM	    teams 
						  WHERE	    slug = ?) 
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

	return fmt.Sprintf(commitRecentTemplate, params.String()), repoList
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
