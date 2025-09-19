DROP TRIGGER IF EXISTS update_place_relations_updated_at ON place_relations;
DROP TRIGGER IF EXISTS update_places_updated_at ON places;

DROP INDEX IF EXISTS idx_places_properties;
DROP INDEX IF EXISTS idx_place_relations_created_at;
DROP INDEX IF EXISTS idx_place_relations_type;
DROP INDEX IF EXISTS idx_place_relations_to_place;
DROP INDEX IF EXISTS idx_place_relations_from_place;
DROP INDEX IF EXISTS idx_places_created_at;
DROP INDEX IF EXISTS idx_places_name;
DROP INDEX IF EXISTS idx_places_geometry;

DROP TABLE IF EXISTS place_relations;
DROP TABLE IF EXISTS places;

DROP TYPE IF EXISTS place_relation_type;