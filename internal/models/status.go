package models

type Status struct {
	UsersCount   int32 `json:"user"`
	ForumsCount  int32 `json:"forum"`
	ThreadsCount int32 `json:"thread"`
	PostsCount   int32 `json:"post"`
}
