package main

// repo interface
type repo interface {
	getComment(id int) (*comment, error)
	getPostComments(path string) ([]comment, error)
	addComment(string, comment) error
	deleteComment(id int) error
}

// mock implementation of the data interface for testing
type dataMock struct { 
	posts map[string]post
}

func (d *dataMock) getComment(id int) (*comment, error) {
	for _, p := range d.posts {
		for _, c := range p.Comments {
			if c.Id == id {
				return &c, nil
			}
		}
	}
	return nil, nil
}

func (d *dataMock) getPostComments(path string) ([]comment, error) {
	return d.posts[path].Comments, nil
}

func (d *dataMock) addComment(p string, c comment) error {
	post := d.posts[p]
	post.Comments = append(post.Comments, c)
	d.posts[p] = post
	return nil
}

func (d *dataMock) deleteComment(id int) error {
	for _, p := range d.posts {
		for i, c := range p.Comments {
			if c.Id == id {
				p.Comments = append(p.Comments[:i], p.Comments[i+1:]...)
				return nil
			}
		}
	}
	return nil
}
