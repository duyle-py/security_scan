create table repository (
    id SERIAL primary key,
    user_id text not null,
    repo_url text not null,
    repo_name text not null
);
CREATE INDEX ON repository (user_id);

create table security_scan (
    id  serial primary key,
    status text not null,
    user_id text not null,
    repo_url text not null,
    repo_name text not null,
    findings bytea,
    queued_at timestamptz,
    scanning_at timestamptz,
    finished_at timestamptz
);
CREATE INDEX ON security_scan (user_id);
