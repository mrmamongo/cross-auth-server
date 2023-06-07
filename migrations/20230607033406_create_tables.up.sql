create table public.user (
    id serial2 primary key unique,
    username text not null,
    telegram_username text not null
)