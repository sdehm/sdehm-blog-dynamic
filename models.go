package main

import (
	"fmt"
	"time"
)

const commentTemplate = `<div id="comment_%d" class="comment">
	<div class="comment-author"> %s </div>
	<div class="comment-date"> %s </div>
	<div class="comment-body"> %s </div>
</div>`

type comment struct {
	Id int
	Author string
	Body string
	Timestamp time.Time
}

type post struct {
	Path string
	Comments []comment
}

func renderComment(c comment) string {
	return fmt.Sprintf(commentTemplate, c.Id, c.Author, c.Timestamp, c.Body)
}

func renderPostComments(p post) string {
	var html string
	for _, c := range p.Comments {
		html += renderComment(c)
	}
	return html
}