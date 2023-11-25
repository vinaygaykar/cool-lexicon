begin;

-- create table lexicon
create table if not exists lexicon(
    word varchar(100) character set utf8 collate utf8_unicode_ci,
    primary key (word)
);

commit;
