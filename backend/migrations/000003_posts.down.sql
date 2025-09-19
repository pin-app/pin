DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_place_comparisons_updated_at ON place_comparisons;
DROP TRIGGER IF EXISTS update_place_ratings_updated_at ON place_ratings;
DROP TRIGGER IF EXISTS update_post_images_updated_at ON post_images;
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP TRIGGER IF EXISTS set_comment_path_trigger ON comments;

DROP FUNCTION IF EXISTS set_comment_path();

DROP INDEX IF EXISTS idx_comments_created_at;
DROP INDEX IF EXISTS idx_comments_path_btree;
DROP INDEX IF EXISTS idx_comments_path;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;
DROP INDEX IF EXISTS idx_place_comparisons_worse_place;
DROP INDEX IF EXISTS idx_place_comparisons_better_place;
DROP INDEX IF EXISTS idx_place_comparisons_user_id;
DROP INDEX IF EXISTS idx_place_ratings_rating;
DROP INDEX IF EXISTS idx_place_ratings_place_id;
DROP INDEX IF EXISTS idx_place_ratings_user_id;
DROP INDEX IF EXISTS idx_post_images_sort_order;
DROP INDEX IF EXISTS idx_post_images_post_id;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_place_id;
DROP INDEX IF EXISTS idx_posts_user_id;

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS place_comparisons;
DROP TABLE IF EXISTS place_ratings;
DROP TABLE IF EXISTS post_images;
DROP TABLE IF EXISTS posts;