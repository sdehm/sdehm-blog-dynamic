package api

import (
	"fmt"

	"github.com/sdehm/sdehm-blog-dynamic/models"
)

const commentTemplate = `<div id="comment_%d" class="comment">
	<div class="comment-author"> %s </div>
	<div class="comment-date"> %s </div>
	<div class="comment-body"> %s </div>
</div>`

func RenderComment(c models.Comment) string {
	return fmt.Sprintf(commentTemplate, c.Id, c.Author, c.Timestamp, c.Body)
}

func RenderPostComments(p models.Post) string {
	var html string = "<div class=\"comments\">"
	for _, c := range p.Comments {
		html += RenderComment(c)
	}
	return html + "</div>"
}