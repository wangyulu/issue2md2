package github

import "time"

// Issue GitHub Issue 数据
type Issue struct {
	Title     string
	Body      string
	Author    string
	AuthorURL string
	CreatedAt time.Time
	Status    string // open, closed
	URL       string
	Reactions *Reactions
	Comments  []Comment
}

// PullRequest GitHub Pull Request 数据
type PullRequest struct {
	Title     string
	Body      string
	Author    string
	AuthorURL string
	CreatedAt time.Time
	Status    string // open, closed, merged
	URL       string
	Reactions *Reactions
	Comments  []Comment // 包含普通评论和 Review Comments
}

// Discussion GitHub Discussion 数据
type Discussion struct {
	Title     string
	Body      string
	Author    string
	AuthorURL string
	CreatedAt time.Time
	Status    string // open, closed
	URL       string
	Reactions *Reactions
	Comments  []Comment
}

// Comment 评论数据
type Comment struct {
	Author    string
	AuthorURL string
	Body      string
	CreatedAt time.Time
	Reactions *Reactions
	IsAnswer  bool // Discussion 特有
}

// Reactions 反应统计
type Reactions struct {
	ThumbsUp   int
	ThumbsDown int
	Laugh      int
	Hooray     int
	Confused   int
	Heart      int
	Rocket     int
	Eyes       int
}
