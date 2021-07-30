CREATE TABLE IF NOT EXISTS users
(
    id text not null primary key,
    email text,
    login text not null,
    password text not null,
    name text,
    unique(email, login)
);
CREATE TABLE IF NOT EXISTS portfolios
(
    id text not null primary key,
    userid text not null,
    name text
);
CREATE TABLE IF NOT EXISTS transactions
(
    id text not null primary key,
    userid text not null,
    portfolioid text not null,
    date text not null,
    asset text not null,
    price real not null,
    quantity integer not null
);
