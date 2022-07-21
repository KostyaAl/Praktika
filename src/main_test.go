package main

import (
	// "err"
	"context"
	"fmt"
	"testing"
)

var ()

func TestGetUsernInfo(t *testing.T) {
	service := NewGitHubService(context.Background())
	got, err := service.GetUserInfo("kill-your-soul")

	want := User{
		UserName:       "kill-your-soul",
		FullName:       "kill_your_soul",
		FollowersCount: 4,
		FollowingCount: 5,
	}
	if *got != want {
		t.Errorf("got %q, wanted %q\n error: %q", got, want, err)
	} else {
		fmt.Println("Success")
	}
}

func main(){
	ctx := context.Background()
	service := NewGitHubService(ctx)
	me, _ := service.GetUserInfo("kill-your-soul")
	fmt.Printf("me: %v\n", me)
	fmt.Printf("err: %v\n", err)
	repo, _ := service.GetRepositoryByName("kill-your-soul", "XakepParser")
	fmt.Printf("repo: %v\n", repo)
	branches, _ := service.GetRepositoryBranches("XakepParser")
	fmt.Printf("branches: %v\n", branches)
	err := service.CreateRepository("test")
	fmt.Printf("err: %v\n", err)
	err = service.DeleteTag("test", "testing_new_features")
	fmt.Printf("err: %v\n", err)
	user, _ := service.GetRepositoryContributors("XakepParser")
	fmt.Printf("user: %v\n", user)
	err = service.DeleteBranch("XakepParser", "test")
	fmt.Printf("err: %v\n", err)
	com, _ := service.GetBranchCommits("kill-your-soul", "XakepParser", "main")
	fmt.Printf("com: %v\n", com)
	pull, _ := service.GetRepositoryPullRequests("XakepParser")
	fmt.Printf("pull: %v\n", pull)
	repos, _ := service.GetUserRepositories("kill-your-soul")
	fmt.Printf("repos: %v\n", repos)
}


