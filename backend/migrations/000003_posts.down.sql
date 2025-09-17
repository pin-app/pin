DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_place_comparisons_updated_at ON place_comparisons;
DROP TRIGGER IF EXISTS update_place_ratings_updated_at ON place_ratings;
DROP TRIGGER IF EXISTS update_post_images_updated_at ON post_images;
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP TRIGGER IF EXISTS set_comment_path_trigger ON comments;

DROP FUNCTION IF EXISTS set_comment_path();

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS place_comparisons;
DROP TABLE IF EXISTS place_ratings;
DROP TABLE IF EXISTS post_images;
DROP TABLE IF EXISTS posts;