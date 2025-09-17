DROP TRIGGER IF EXISTS update_place_relations_updated_at ON place_relations;
DROP TRIGGER IF EXISTS update_places_updated_at ON places;

DROP TABLE IF EXISTS place_relations;
DROP TABLE IF EXISTS places;

DROP TYPE IF EXISTS place_relation_type;