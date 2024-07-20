CREATE EXTENSION IF NOT EXISTS postgis;

CREATE INDEX idx_properties_location ON properties USING GIST (ST_GeographyFromText('SRID=4326;POINT(' || longitude || ' ' || latitude || ')'));

