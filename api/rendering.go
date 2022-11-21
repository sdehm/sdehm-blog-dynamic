package api

import (
	"fmt"

	"github.com/sdehm/sdehm-blog-dynamic/models"
)

const commentTemplate = `<div>
<hr class="border-dotted border-neutral-300 dark:border-neutral-600">
<div id="comment_%d" class="comment">
	<div class="comment-author font-bold text-xs text-neutral-500 dark:text-neutral-400"> %s </div>
	<span class="comment-date mt-[0.1rem] text-xs text-neutral-500 dark:text-neutral-400"> 
	  <time datetime="%s"> %s </time> 
	</span>
	<div class="comment-body"> %s </div>
</div>
</div>`

const postTemplate = `<div id="comments">
<p>Comments</p>
<form id="comment-form" class="w-full max-w-xs" action="#">
	<div class="form-control">
		<label class="label block text-sm mb-1 text-neutral-500 dark:text-neutral-400">
			<span class="label-text">Name</span>
		</label>
		<input required type="text" placeholder="Name" name="name" class="rounded bg-transparent appearance-none focus:outline-dotted focus:outline-2 focus:outline-transparent">
	</div>
	<div class="form-control">

		<label class="label block text-sm mb-1 text-neutral-500 dark:text-neutral-400">
			<span class="label-text">Comment</span>
		</label>
		<textarea required name="comment" class="rounded bg-transparent appearance-none focus:outline-dotted focus:outline-2 focus:outline-transparent h-24 w-full" placeholder="Comment"></textarea>
	</div>
	<div class="form-control">
		<button class="border-2 border-neutral-300 dark:border-neutral-600 font-bold py-2 px-2 rounded hover:border-4 mb-2">Submit</button>
	</div>
</form>
<div id="comments" class="comments">
%s
</div>`

func RenderComment(c models.Comment) string {
	tf := c.Timestamp.Format("2 January 2006")
	return fmt.Sprintf(commentTemplate, c.Id, c.Author, c.Timestamp, tf, c.Body)
}

func RenderPostComments(p models.Post) string {
	var html string
	for _, c := range p.Comments {
		// prepend rendered comment to html
		html = RenderComment(c) + html
	}
	return fmt.Sprintf(postTemplate, html)
}
