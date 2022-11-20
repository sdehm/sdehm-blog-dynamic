package models

import (
	"fmt"
	"time"
)

const commentTemplate = `<div id="comment_%d" class="comment">
	<div class="comment-author"> %s </div>
	<div class="comment-date"> %s </div>
	<div class="comment-body"> %s </div>
</div>`

type Comment struct {
	Id int
	Author string
	Body string
	Timestamp time.Time
}

type Post struct {
	Path string
	Comments []Comment
}

func renderComment(c Comment) string {
	return fmt.Sprintf(commentTemplate, c.Id, c.Author, c.Timestamp, c.Body)
}

func renderPostComments(p Post) string {
	var html string
	for _, c := range p.Comments {
		html += renderComment(c)
	}
	return html
}