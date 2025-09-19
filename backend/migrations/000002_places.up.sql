CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TYPE place_relation_type AS ENUM ('CONTAINS', 'PART_OF', 'OVERLAPS');

-- vertices
CREATE TABLE IF NOT EXISTS places (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    geometry GEOMETRY(GEOMETRY, 4326),
    properties JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- edges between places
CREATE TABLE IF NOT EXISTS place_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    to_place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    relation_type place_relation_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(from_place_id, to_place_id, relation_type)
);

CREATE INDEX IF NOT EXISTS idx_places_geometry ON places USING GIST (geometry);
CREATE INDEX IF NOT EXISTS idx_places_name ON places(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_places_created_at ON places(created_at);
CREATE INDEX IF NOT EXISTS idx_places_properties ON places USING GIN (properties);

CREATE INDEX IF NOT EXISTS idx_place_relations_from_place ON place_relations(from_place_id);
CREATE INDEX IF NOT EXISTS idx_place_relations_to_place ON place_relations(to_place_id);
CREATE INDEX IF NOT EXISTS idx_place_relations_type ON place_relations(relation_type);
CREATE INDEX IF NOT EXISTS idx_place_relations_created_at ON place_relations(created_at);
CREATE INDEX IF NOT EXISTS idx_places_properties ON places USING GIN (properties);

CREATE TRIGGER update_places_updated_at BEFORE UPDATE ON places
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_place_relations_updated_at BEFORE UPDATE ON place_relations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
