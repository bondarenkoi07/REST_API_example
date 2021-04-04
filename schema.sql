create table users(
    id serial not null,
    login varchar(32) not null unique,
    password varchar(32) not null,
    PRIMARY KEY(id)
);



create table developers(
    id integer not null,
    org_name varchar(128) unique not null,
    section varchar(128) not null,
    FOREIGN KEY (id) references users(id) on delete cascade,
    PRIMARY KEY(id)
);

create table markets(
    id serial not null,
    name varchar(128) not null,
    max_products integer,
    PRIMARY KEY(id)
);

create table product(
    id serial not null,
    name varchar(128) not null,
    cost smallint not null,
    count smallint  not null,
    developerId integer not null,
    marketId integer,
    FOREIGN KEY (developerId) references developers(id) on delete cascade,
    FOREIGN KEY (marketId ) references markets(id) on delete cascade
);
