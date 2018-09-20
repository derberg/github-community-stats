package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/derberg/github-community-stats/internal/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var token = os.Getenv("GITHUB_TOKEN")
var repoName = os.Getenv("GITHUB_REPO")
var orgName = os.Getenv("GITHUB_ORG")
var orgId = os.Getenv("GITHUB_ORG_ID")
var ctx = context.Background()

type User struct {
	Name, Email, Company, Location string
	IsHireable                     bool
	Organizations                  struct {
		Nodes []struct {
			Id   string
			Name string
		}
	}
}

type User1 struct {
	Name, Email, Company, Location string
	IsHireable                     bool
	Commits                        int
	Organizations                  []string
}

type MyBox struct {
	Items []User1
}

func main() {
	flag.Parse()

	client := graphqlClient()
	users := []User1{}
	box := MyBox{users}

	contributors, err := github.GetContributors(repoName, orgName, token)
	if err != nil {
		panic(err)
	}

	for _, contributor := range contributors {

		username := contributor.GetAuthor().GetLogin()
		totalContrib := contributor.GetTotal()
		orgMember := false
		orgNames := []string{}
		orgs := getUserOrgs(client, username)
		if err != nil {
			panic(err)
		}
		var raw User

		json.Unmarshal([]byte(orgs), &raw)

		for _, org := range raw.Organizations.Nodes {

			orgsJson, err := json.Marshal(org.Id)
			if err != nil {
				panic(err)
			}
			userOrgId, err := strconv.Unquote(string(orgsJson))

			if userOrgId == orgId {
				orgMember = true
			}

			orgNames = append(orgNames, org.Name)

		}

		if orgMember == false {
			user := User1{
				Name:          fmt.Sprint("https://github.com/", username),
				Email:         raw.Email,
				Company:       raw.Company,
				Location:      raw.Location,
				IsHireable:    raw.IsHireable,
				Commits:       totalContrib,
				Organizations: orgNames,
			}
			box.Items = append(box.Items, user)
			//users = append(users.Name, user)
			//users[i].Name = user
		}

	}

	printJSON(box.Items)
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

//kyma org id MDEyOk9yZ2FuaXphdGlvbjM5MTUzNTIz
func getUserOrgs(client *githubv4.Client, username string) string {
	{
		/*
								  user(login:"derberg") {
									name
			    					email
			    					company
			    					location
			    					isHireable
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
				Name          githubv4.String
				Email         githubv4.String
				Company       githubv4.String
				Location      githubv4.String
				IsHireable    githubv4.Boolean
				Organizations struct {
					Nodes []struct {
						Id   githubv4.String
						Name githubv4.String
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
			panic(err)
		}
		orgs := getUserOrgs.User
		//printJSON(getUserOrgs.User)

		c, err := json.Marshal(orgs)
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(c))
		//userData := json.Marshal(orgs)
		return string(c)
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

// query for getting forks and issues per repo
func getIssuesAndForks(client *githubv4.Client) error {

	// Query some details about a repository, an issue in it, and its comments.
	{
		/*
						{
			  organization(login: "kyma-project") {
			    repositories(first: 10) {
			      edges {
			        node {
			          name
			          issues(first:100 after:"Y3Vyc29yOnYyOpHOFTLViw==") {
						      pageInfo {
						        endCursor
						        hasNextPage
						      }
			            edges {
						        cursor
						        node {
						          number
						          state
						          author {
						            login
						          }
						        }
						      }
			          }
			          forks(first: 100) {
						      totalCount
						      pageInfo {
						        endCursor
						        hasNextPage
						      }
						      edges {
						        node {
						          owner {
						            login
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

/*
// query for getting forks and issues per repo
func getIssuesAndForks(client *githubv4.Client, username string) string {
	{
		/*
			{
			  repository(name: "kyma", owner: "kyma-project") {
			    issues(first:100 after:"Y3Vyc29yOnYyOpHOFTLViw==") {
			      pageInfo {
			        endCursor
			        hasNextPage
			      }
			      edges {
			        cursor
			        node {
			          number
			          state
			          author {
			            login
			          }
			        }
			      }
			    }
			    forks(first: 100) {
			      totalCount
			      pageInfo {
			        endCursor
			        hasNextPage
			      }
			      edges {
			        node {
			          owner {
			            login
			          }
			        }
			      }
			    }
			  }
			}

*/
/*
		var getIssuesAndForks struct {
			Repositories struct
			Issue struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.String
				}
				Edges struct {
					Cursor githubv4.String
					Node   struct {
						Number githubv4.Int
						State  githubv4.String
						Author struct {
							Login githubv4.String
						}
					}
				}
			} `graphql:"issues(first:$orgLimit after:$afterCursor)"`
			Fork struct {
				TotalCount githubv4.Int
				PageInfo   struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.String
				}
				Edges struct {
					Cursor githubv4.String
					Node   struct {
						Owner struct {
							Login githubv4.String
						}
					}
				}
			} `graphql:"forks(first:$orgLimit after:$afterCursor)"`
		}

		variables := map[string]interface{}{
			"login":    githubv4.String(username),
			"orgLimit": githubv4.Int(10),
		}

		err := client.Query(ctx, &getUserOrgs, variables)
		if err != nil {
			panic(err)
		}
		orgs := getUserOrgs.User
		//printJSON(getUserOrgs.User)

		c, err := json.Marshal(orgs)
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(c))
		//userData := json.Marshal(orgs)
		return string(c)
	}

}
*/
