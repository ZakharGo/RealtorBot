CREATE TABLE flats(
    id serial not null unique,
    flat varchar(10) unique ,
    record int
);
CREATE TABLE records(
    id serial not null unique,
    date TIMESTAMP,
    count INT
);
CREATE TABLE flat_records(
    id serial not null unique,
    flat_id int references flats(id) on delete cascade not null,
    record_id int references records(id) on delete cascade not null
);

CREATE TABLE scripts(
    script varchar(255) not null
)


