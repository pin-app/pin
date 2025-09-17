CREATE EXTENSION IF NOT EXISTS "ltree";

CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS post_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    caption TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- for relative ranking system
CREATE TABLE IF NOT EXISTS place_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, place_id)
);

-- rankings derived from relative comparisons
CREATE TABLE IF NOT EXISTS place_comparisons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    better_place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    worse_place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, better_place_id, worse_place_id)
);

-- hierarchical ltree for comments
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    path LTREE NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_place_id ON posts(place_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);

CREATE INDEX IF NOT EXISTS idx_post_images_post_id ON post_images(post_id);
CREATE INDEX IF NOT EXISTS idx_post_images_sort_order ON post_images(post_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_place_ratings_user_id ON place_ratings(user_id);
CREATE INDEX IF NOT EXISTS idx_place_ratings_place_id ON place_ratings(place_id);
CREATE INDEX IF NOT EXISTS idx_place_ratings_rating ON place_ratings(rating) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_place_comparisons_user_id ON place_comparisons(user_id);
CREATE INDEX IF NOT EXISTS idx_place_comparisons_better_place ON place_comparisons(better_place_id);
CREATE INDEX IF NOT EXISTS idx_place_comparisons_worse_place ON place_comparisons(worse_place_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_path ON comments USING GIST (path);
CREATE INDEX IF NOT EXISTS idx_comments_path_btree ON comments USING BTREE (path);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);

-- automatically set comment path based on hierarchy
CREATE OR REPLACE FUNCTION set_comment_path()
RETURNS TRIGGER AS $$
DECLARE
    parent_path LTREE;
BEGIN
    -- if this is a top-level comment, set path to the comment ID
    IF NEW.parent_id IS NULL THEN
        NEW.path := NEW.id::TEXT::LTREE;
    ELSE
        -- get parent path and append this comment's ID
        SELECT path INTO parent_path FROM comments WHERE id = NEW.parent_id;
        NEW.path := parent_path || NEW.id::TEXT::LTREE;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_comment_path_trigger BEFORE INSERT ON comments
    FOR EACH ROW EXECUTE FUNCTION set_comment_path();


CREATE TRIGGER update_posts_updated_at BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_post_images_updated_at BEFORE UPDATE ON post_images
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_place_ratings_updated_at BEFORE UPDATE ON place_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_place_comparisons_updated_at BEFORE UPDATE ON place_comparisons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();