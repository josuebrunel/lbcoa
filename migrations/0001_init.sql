create table if not exists stat (
    qs text primary key,
    hits integer default 1
);

create index if not exists idx_stat_qs on stat(qs);