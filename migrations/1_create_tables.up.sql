create table repository (
    id SERIAL primary key,
    user_id text not null,
    repo_url text not null,
    repo_name text not null
);

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