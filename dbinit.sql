DROP TABLE comments;
DROP TABLE posts;

CREATE TABLE posts (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid() ,
    path VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE comments (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    author VARCHAR(20) NOT NULL,
    body VARCHAR(255) NOT NULL,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    INDEX (post_id)
);
