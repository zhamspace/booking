do
$$
    begin
        execute 'ALTER DATABASE ' || current_database() || ' SET timezone = ''+05''';
    end;
$$;

create table if not exists audit
(
    id          text        default gen_random_uuid() primary key,
    created_at  timestamptz default now() not null,
    user_id     text                      not null,
    change_type text                      not null,
    resource    text                      not null,
    object_id   text                      not null,
    object      jsonb                     not null,
    metadata    jsonb                     null
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_audit_resource_object_id_created_at ON audit (resource, object_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_resource_created_at_desc ON audit (resource, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_user_created_at_desc ON audit (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_change_type_created_at_desc ON audit (change_type, created_at DESC);