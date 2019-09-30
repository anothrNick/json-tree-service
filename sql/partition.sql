CREATE TABLE app_projects (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    slug VARCHAR NOT NULL UNIQUE,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE project_resource_objects (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    resource_path VARCHAR NOT NULL,
    creator_type VARCHAR NOT NULL,
    creator uuid,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    data JSONB
);

CREATE OR REPLACE FUNCTION create_partition_and_insert() RETURNS trigger AS
  $BODY$
    DECLARE
      partition TEXT;
    BEGIN
      partition := TG_RELNAME || '_' || MD5(NEW.project_id::VARCHAR);
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
        RAISE NOTICE 'A partition has been created %',partition;
        EXECUTE 'CREATE TABLE ' || partition || ' (check (project_id = ''' || NEW.project_id || ''')) INHERITS (' || TG_RELNAME || ');';
      END IF;
      EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING id;';
      RETURN NULL;
    END;
  $BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE TRIGGER testing_partition_insert_trigger
BEFORE INSERT ON project_resource_objects
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();