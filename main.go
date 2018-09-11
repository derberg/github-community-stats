package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Println(err)
	}
}

func run() error {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_GRAPHQL_TEST_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	// Query some details about a repository, an issue in it, and its comments.
	{
		type githubV4Actor struct {
			Login     githubv4.String
			AvatarURL githubv4.URI `graphql:"avatarUrl(size:72)"`
			URL       githubv4.URI
		}

		/*
			organization(login: "kyma-project") {
				repositories(first: 10) {
				  edges {
					node {
					  name
					  collaborators(first: 100 affiliation: ALL) {
						totalCount
						edges {
						  node {
							name
							login
							organizations(first: 5) {
							  edges {
								node {
								  name
								}
							  }
							}
						  }
						}
					  }
					}
				  }
				}
			  }
		*/

		var getCollaborators struct {
			Organization struct {
				Repositories struct {
					Edges []struct {
						Node struct {
							Name          githubv4.String
							Collaborators struct {
								TotalCount githubv4.Int
								Edges      []struct {
									Node struct {
										Name          githubv4.String
										Login         githubv4.String
										Organizations struct {
											Edges []struct {
												Node struct {
													Name githubv4.String
												}
											}
										} `graphql:"organizations(first:$orgLimit)"`
									}
								}
							} `graphql:"collaborators(first:$collaboratorsLimit)"`
						}
					}
				} `graphql:"repositories(first:$repoLimit)"`
			} `graphql:"organization(login:$repositoryOwner)"`
		}

		variables := map[string]interface{}{
			"repositoryOwner":    githubv4.String("kyma-project"),
			"repoLimit":          githubv4.Int(10),
			"orgLimit":           githubv4.Int(10),
			"collaboratorsLimit": githubv4.Int(100),
		}
		err := client.Query(context.Background(), &getCollaborators, variables)
		if err != nil {
			return err
		}
		printJSON(getCollaborators)
	}

	return nil
}

// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
