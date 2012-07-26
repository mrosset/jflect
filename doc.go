// Copyright 2012 The jflect Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Jflect takes json from stdin and outputs go structs to stdout.

Installing

    go get github.com/str1ngs/jflect

Examples

Web service:

    curl https://api.github.com/users/str1ngs | jflect -s GithubUser

Output:

type GithubUser struct {
	AvatarUrl   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	GravatarId  string `json:"gravatar_id"`
	Hireable    bool   `json:"hireable"`
	HtmlUrl     string `json:"html_url"`
	Id          int    `json:"id"`
	Location    string `json:"location"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicGists int    `json:"public_gists"`
	PublicRepos int    `json:"public_repos"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}

From file:

	jflect < foo.json 
	

Saving is just a matter of redirecting stdout to a file.

	jflect < foo.json > foo.go
	curl https://api.github.com/users/str1ngs | jflect -s GihubUser > githuser.go


Notes

Jflect is primarily designed to fast prototype go structs. Its not
intended for automation. It does handle nested json objects. But does
currently handle arbitrary json arrays

*/
package documentation
