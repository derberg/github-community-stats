package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var token = os.Getenv("GITHUB_TOKEN")
var repoName = os.Getenv("GITHUB_REPO")
var orgName = os.Getenv("GITHUB_ORG")
var orgId = os.Getenv("GITHUB_ORG_ID")
var ctx = context.Background()

func main() {
	flag.Parse()

	client := graphqlClient()

	contributors, err := getContributors(repoName, orgName)
	if err != nil {
		panic(err)
	}

	for _, contributor := range contributors {

		username := contributor.GetAuthor().GetLogin()
		orgMember := false
		orgs, err := getUserOrgs(client, username)
		if err != nil {
			panic(err)
		}

		for _, org := range orgs {

			orgsJson, err := json.Marshal(org.Id)
			if err != nil {
				panic(err)
			}
			userOrgId, err := strconv.Unquote(string(orgsJson))

			if userOrgId == orgId {
				orgMember = true
			}

		}

		if orgMember == false {
			fmt.Printf("https://github.com/%v\n", username)
		}

	}
}

//Get a GraphQL Client to make calls to the API
func graphqlClient() *githubv4.Client {

	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	return client
}

func getCollaborators(client *githubv4.Client) error {

	// Query some details about a repository, an issue in it, and its comments.
	{
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
		err := client.Query(ctx, &getCollaborators, variables)
		if err != nil {
			return err
		}
		printJSON(getCollaborators)
	}
	return nil
}

//kyma org id MDEyOk9yZ2FuaXphdGlvbjM5MTUzNTIz
func getUserOrgs(client *githubv4.Client, username string) ([]struct{ Id githubv4.String }, error) {
	{
		/*
					  user(login:"derberg") {
			    		organizations(first:10) {
			      			totalCount
			      			nodes {
			        			id
			        			name
			      			}
			    		}
			  		  }
		*/

		var getUserOrgs struct {
			User struct {
				Organizations struct {
					Nodes []struct {
						Id githubv4.String
					}
				} `graphql:"organizations(first:$orgLimit)"`
			} `graphql:"user(login:$login)"`
		}

		variables := map[string]interface{}{
			"login":    githubv4.String(username),
			"orgLimit": githubv4.Int(10),
		}

		err := client.Query(ctx, &getUserOrgs, variables)
		if err != nil {
			return nil, err
		}
		orgs := getUserOrgs.User.Organizations.Nodes

		return orgs, err
	}

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

// rest api
func getContributors(repoName string, orgName string) ([]*github.ContributorStats, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := github.NewClient(httpClient)

	contributors, _, err := client.Repositories.ListContributorsStats(ctx, orgName, repoName)
	return contributors, err
}
