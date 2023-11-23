-- begin and commit are not needed as per https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/database/sqlite/README.md#L5
-- create table lexicon
create table if not exists lexicon(
    word varchar(100) collate nocase,
    primary key (word)
);
