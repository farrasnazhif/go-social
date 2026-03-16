CREATE EXTENSION IF NOT EXISTS pg_trgm;


-- trigram indexes for text search
CREATE INDEX IF NOT EXISTS idx_comments_content
ON comments USING gin (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title
ON posts USING gin (title gin_trgm_ops);

-- array index
CREATE INDEX IF NOT EXISTS idx_posts_tags
ON posts USING gin (tags);

-- unique value, using b-tree index not gin

CREATE INDEX IF NOT EXISTS idx_users_username
ON users (username);

-- relational lookup indexes (B-tree)
CREATE INDEX IF NOT EXISTS idx_posts_user_id
ON posts (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id
ON comments (post_id);