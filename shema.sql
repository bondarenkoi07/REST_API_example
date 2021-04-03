create table users(
    id serial not null,
    login varchar(32) not null unique,
    password varchar(32) not null,
    role varchar (32) not null
);



create table developers(
    id integer not null,
    org_name varchar(128) unique not null,
    section varchar(128) not null,
    FOREIGN KEY (id) references users(id) on delete cascade,
    PRIMARY KEY(id)
);

create table product(
    id serial not null,
    name varchar(128) not null,
    cost smallint not null,
    count smallint  not null,
    developerId integer unique,
    FOREIGN KEY (developerId) references developers(id) on delete cascade
);