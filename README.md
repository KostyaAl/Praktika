# Задание
Библиотека-обертка для API Github. Позволяет в более удобном виде получать доступ к API Github. 

## Пример использования:
Изначально нужно установить в переменные окружения персональный токен доступа к Github, который можно получить в настройках пользователя с названием ```TOKEN```
```Powershell
$env:TOKEN=YOUR_TOKEN_HERE
```
---
```go
ctx := context.Background()
service := NewGitHubService(ctx)
me, err := service.GetUserInfo("kill-your-soul")
fmt.Printf("me: %v\n", me)
fmt.Printf("err: %v\n", err)
```