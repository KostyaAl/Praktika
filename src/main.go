package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

// User хранит краткую информацию о пользователе
type User struct {
	UserName       string // GitHub username пользователя
	FullName       string // Полное имя пользователя
	FollowersCount int    // Количество подписчиков
	FollowingCount int    // Количество подписок
}

// Repository хранит информацию о репозиториях пользователя
type Repository struct {
	Name            string    // Название репозитория
	Description     string    // Краткое описание репозитория
	Link            string    // Ссылка на репозиторий
	IsPrivate       bool      // Приватный репозиторий или открытый
	StarsCount      int       // Количество звезд
	ForksCount      int       // Количество форков (ответвлений, сделанных другими пользователями)
	LastUpdatedTime time.Time // Время последнего изменения

	// programmingLanguage хранит информацию об используемых языках программирования
	programmingLanguage []struct {
		Name           string  // Название языка программирования
		PercentOfUsage float64 // Процент использования в репозитории
	}
}

type Branch struct {
	Name      string    // Название ветки
	UpdatedAt time.Time // Дата последнего обновления
}

type Commit struct {
	Hash      string    // SHA коммита
	Title     string    // Сообщение коммита
	CreatedAt time.Time // Дата создания коммита
}

type Issue struct {
	Title                   string    // Тема issue
	IsClosed                bool      // Актуальная или разрешенная проблема
	ResolvedPullRequestLink string    // Ссылка на PR, в котором разрешена проблема
	CreatedAt               time.Time // Дата создания
	UpdatedAt               time.Time // Дата обновления
}

type PullRequest struct {
	ID           int    // Номер запроса на слияние (отображен в url как /pulls/{id})
	Title        string // Название запроса на слияние
	SourceBranch string // Название ветки-источника
	TargetBranch string // Название ветки-назначения
	IsClosed     bool   // Закрыт или открыт
}

type Thread struct {
	IsResolved bool // Закрытое или открытое обсуждение
}

type Tag struct {
	Title       string    // Название тега
	Hash        string    // SHA, хэш
	Description string    // Описание тега
	ZipLink     string    // Ссылка на скачивание архива
	CreatedAt   time.Time // Дата создания
}

type GitServiceIFace interface {
	// GetUserInfo получает основную информацию о пользователе
	GetUserInfo(userName string) (*User, error)

	// GetUserRepositories получает список всех репозиториев пользователя
	GetUserRepositories(userName string) ([]*Repository, error)

	// GetRepositoryByName получает информацию об указанном репозитории
	GetRepositoryByName(userName, repositoryName string) (*Repository, error)

	// CreateRepository создает репозиторий с указанным именем
	CreateRepository(repositoryName string) error

	// GetRepositoryBranches получает список всех веток репозитория
	GetRepositoryBranches(repositoryName string) ([]*Branch, error)

	// CreateBranch создает новую ветку
	CreateBranch(repoName, branchName string) error

	// DeleteBranch удаляет указанную ветку
	DeleteBranch(repoName, branchName string) error

	// GetBranchCommits возвращает коммиты указанной ветки
	GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error)

	// GetRepositoryPullRequests получает информацию о запросах на слияние
	GetRepositoryPullRequests(repositoryName string) ([]*PullRequest, error)

	// CreatePullRequest создает новый запрос на слияние
	CreatePullRequest(sourceBranch, destBranch, title string) error

	// GetThreadsInfo получает информацию об обсуждениях конкретного запроса на слияние
	GetThreadsInfo(repositoryName string, pullRequestID int) ([]*Thread, error)

	// GetIssues получает информацию об опубликованных проблемах репозитория
	GetIssues(repositoryName string) ([]*Issue, error)

	// GetRepositoryContributors получает список соавторов репозитория
	GetRepositoryContributors(repositoryName string) ([]*User, error)

	// GetRepositoryTags возвращает информацию о тегах репозитория
	GetRepositoryTags(userName, repositoryName string) ([]*Tag, error)

	// CreateTag создает новый тег
	CreateTag(title string) error

	// DeleteTag удаляет тег по имени
	DeleteTag(repositoryName, tagName string) error

	// SetAccessToRepository предоставляет доступ к репозиторию указанному пользователю
	SetAccessToRepository(oppoUserName, repositoryName string) error

	// DenyAccessToRepository закрывает доступ к репозиторию указанному пользователю
	DenyAccessToRepository(oppoUserName, repositoryName string) error
}

// Структура, реализующая интерфейс GitServiceIFace
type gitHubService struct {
	client  *github.Client
	context context.Context
}

// Структура для доступа к токену из переменных окружения

// NewGitHubService - конструктор gitHubService
func NewGitHubService(ctx context.Context) GitServiceIFace {

	token := os.Getenv("TOKEN")
	// fmt.Printf("creds.token: %v\n", token)
	ts := oauth2.StaticTokenSource(
		// Передаем Oauth2.0-токен, который можно получить в настройках профиля GitHub
		// Токен необходимо передавать из переменных окружения!
		// Пример библиотеки: https://github.com/caarlos0/env
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Запросы к GitHub API будут отправлены от имени аутентифицированного пользователя
	client := github.NewClient(tc)

	return &gitHubService{
		client:  client,
		context: ctx,
	}
}

// Необходимо реализовать нижепредставленные методы в соответствии со структурой интерфейса
//                                   |
//                                   |
//                                   |
//                                   V

func (ghs *gitHubService) GetUserInfo(userName string) (*User, error) {
	user_tmp, _, err := ghs.client.Users.Get(ghs.context, userName)
	if user_tmp.Login == nil || user_tmp.Name == nil || user_tmp.Followers == nil || user_tmp.Following == nil {
		return nil, fmt.Errorf("error while parsing")
	}
	user := User{
		UserName:       *user_tmp.Login,
		FullName:       *user_tmp.Name,
		FollowersCount: *user_tmp.Followers,
		FollowingCount: *user_tmp.Following,
	}
	return &user, err
}

func (ghs *gitHubService) GetUserRepositories(userName string) ([]*Repository, error) {
	var repositories []*Repository
	tmp, _, err := ghs.client.Repositories.List(ghs.context, userName, &github.RepositoryListOptions{})
	for _, item := range tmp {
		langs, _, _ := ghs.client.Repositories.ListLanguages(ghs.context, userName, *item.Name)
		sum := 0.0
		for lang := range langs {
			sum += float64(langs[lang])

		}
		repo := Repository{
			Name:            *item.FullName,
			Description:     *item.Description,
			Link:            *item.HTMLURL,
			IsPrivate:       *item.Private,
			StarsCount:      *item.StargazersCount,
			ForksCount:      *item.ForksCount,
			LastUpdatedTime: item.UpdatedAt.Time,
			programmingLanguage: []struct {
				Name           string
				PercentOfUsage float64
			}{},
		}
		for lang := range langs {
			repo.programmingLanguage = append(repo.programmingLanguage, struct {
				Name           string
				PercentOfUsage float64
			}{Name: lang, PercentOfUsage: ((float64(langs[lang]) / sum) * 100)})
		}
		repositories = append(repositories, &repo)
	}
	return repositories, err
}

func (ghs *gitHubService) GetRepositoryByName(userName, repositoryName string) (*Repository, error) {
	repo_tmp, _, err := ghs.client.Repositories.Get(ghs.context, userName, repositoryName)
	langs, _, _ := ghs.client.Repositories.ListLanguages(ghs.context, userName, repositoryName)
	sum := 0.0
	for lang := range langs {
		sum += float64(langs[lang])

	}
	repo := Repository{
		Name:            *repo_tmp.FullName,
		Description:     *repo_tmp.Description,
		Link:            *repo_tmp.HTMLURL,
		IsPrivate:       *repo_tmp.Private,
		StarsCount:      *repo_tmp.StargazersCount,
		ForksCount:      *repo_tmp.ForksCount,
		LastUpdatedTime: repo_tmp.UpdatedAt.Time,
		programmingLanguage: []struct {
			Name           string
			PercentOfUsage float64
		}{},
	}
	for lang := range langs {
		repo.programmingLanguage = append(repo.programmingLanguage, struct {
			Name           string
			PercentOfUsage float64
		}{Name: lang, PercentOfUsage: ((float64(langs[lang]) / sum) * 100)})
	}
	return &repo, err
}

func (ghs *gitHubService) CreateRepository(repositoryName string) error {
	_, _, err := ghs.client.Repositories.Create(ghs.context, "", &github.Repository{
		Name: &repositoryName,
	})
	return err
}

func (ghs *gitHubService) GetRepositoryBranches(repositoryName string) ([]*Branch, error) {
	var branch []*Branch
	user, _, err := ghs.client.Users.Get(ghs.context, "")
	username := user.Login
	branches, _, _ := ghs.client.Repositories.ListBranches(ghs.context, *username, repositoryName, &github.BranchListOptions{})
	for _, item := range branches {
		time, _, _ := ghs.client.Repositories.GetBranch(ghs.context, *username, repositoryName, *item.Name, true)
		branch = append(branch, &Branch{Name: *item.Name, UpdatedAt: *time.Commit.Commit.Committer.Date})
	}
	fmt.Printf("branches: %v\n", branches)
	return branch, err
}

func (ghs *gitHubService) CreateBranch(repoName, branchName string) error {
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	ref := "heads/" + branchName
	branch, _, _ := ghs.client.Repositories.GetBranch(ghs.context, *owner, repoName, "main", true)
	sha := branch.Commit.SHA
	_, _, err := ghs.client.Git.CreateRef(ghs.context,
		*owner, repoName,
		&github.Reference{Ref: &ref, Object: &github.GitObject{SHA: sha}})
	return err
}

func (ghs *gitHubService) DeleteBranch(repoName, branchName string) error {
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	_, err := ghs.client.Git.DeleteRef(ghs.context, *owner, repoName, "heads/"+branchName)
	return err
}

func (ghs *gitHubService) GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error) {
	var commits []*Commit
	branch, _, _ := ghs.client.Repositories.GetBranch(ghs.context, userName, repositoryName, branchName, true)
	sha := branch.Commit.SHA
	comms, _, err := ghs.client.Repositories.ListCommits(ghs.context, userName, repositoryName, &github.CommitsListOptions{SHA: *sha})
	fmt.Printf("commits: %v\n", comms)
	for _, item := range comms {
		commits = append(commits, &Commit{Hash: *item.SHA, Title: *item.Commit.Message, CreatedAt: *item.Commit.Committer.Date})
	}
	return commits, err
}

func (ghs *gitHubService) GetRepositoryPullRequests(repositoryName string) ([]*PullRequest, error) {
	var pullrequests []*PullRequest
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	tmp, _, err := ghs.client.PullRequests.List(ghs.context, *owner, repositoryName, &github.PullRequestListOptions{State: "all"})
	// fmt.Printf("tmp: %v\n", tmp)
	for _, item := range tmp {
		pullrequests = append(pullrequests, &PullRequest{ID: int(*item.ID), Title: *item.Title, SourceBranch: *item.GetBase().Repo.BranchesURL, TargetBranch: *item.State, IsClosed: *item.Merged})
	}
	return pullrequests, err
}

func (ghs *gitHubService) CreatePullRequest(sourceBranch, destBranch, title string) error {
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	_, _, err := ghs.client.PullRequests.Create(ghs.context, *owner, "test", &github.NewPullRequest{})
	// fmt.Printf("tmp: %v\n", tmp)
	return err
}

func (ghs *gitHubService) GetThreadsInfo(repositoryName string, pullRequestID int) ([]*Thread, error) {
	var threads []*Thread
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	tmp, _, err := ghs.client.Repositories.Threads(ghs.context, owner, repositoryName, pullRequestID)
	for _, item := range tmp{
		threads = append(threads, &Thread{IsResolved: item.IsResolved})
	}
	return threads, err
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetIssues(repositoryName string) ([]*Issue, error) {
	var issues []*Issue
	tmp, _, err := ghs.client.Issues.List(ghs.context, true, &github.IssueListOptions{})
	for _, item := range tmp {
		if item.Repository.Name == &repositoryName {
			issues = append(issues, &Issue{Title: *item.Title, IsClosed: *item.Locked, ResolvedPullRequestLink: *item.PullRequestLinks.URL, CreatedAt: *item.CreatedAt, UpdatedAt: *item.UpdatedAt})
		}
	}
	return issues, err
}

func (ghs *gitHubService) GetRepositoryContributors(repositoryName string) ([]*User, error) {
	var user []*User
	_user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := _user.Login
	cont, _, err := ghs.client.Repositories.ListCollaborators(ghs.context, *owner, repositoryName, &github.ListCollaboratorsOptions{})
	for _, item := range cont {
		tmp, _ := ghs.GetUserInfo(*item.Login)
		user = append(user, tmp)
	}
	return user, err
}

func (ghs *gitHubService) GetRepositoryTags(userName, repositoryName string) ([]*Tag, error) {
	var tags []*Tag
	tmp, _, err := ghs.client.Repositories.ListTags(ghs.context, userName, repositoryName, &github.ListOptions{})
	for _, item := range tmp {
		tags = append(tags, &Tag{Title: *item.Name, Hash: *item.Commit.SHA, Description: *item.Commit.Message, ZipLink: *item.ZipballURL, CreatedAt: *item.Commit.Committer.Date})
	}
	return tags, err
}

func (ghs *gitHubService) CreateTag(title string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DeleteTag(repositoryName, tagName string) error {
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	username := user.Login
	_, err := ghs.client.Git.DeleteRef(ghs.context, *username, repositoryName, "tags/"+tagName)
	return err
}

func (ghs *gitHubService) SetAccessToRepository(oppoUserName, repositoryName string) error {
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	_, _, err := ghs.client.Repositories.AddCollaborator(ghs.context, *owner, repositoryName, oppoUserName, &github.RepositoryAddCollaboratorOptions{})
	return err
}

func (ghs *gitHubService) DenyAccessToRepository(oppoUserName, repositoryName string) error {
	user, _, _ := ghs.client.Users.Get(ghs.context, "")
	owner := user.Login
	_, err := ghs.client.Repositories.RemoveCollaborator(ghs.context, *owner, repositoryName, oppoUserName)
	return err
}

func main() {
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
