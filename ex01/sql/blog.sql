CREATE TABLE IF NOT EXISTS blog (
    id         integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    message    varchar(280) NOT NULL
);

SELECT * FROM blog;

DROP TABLE blog;