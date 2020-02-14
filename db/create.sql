create table tweets (
    id bigint(20) unsigned not null primary key,
    lang char(2) not null ,
    tweet text not null ,
    username varchar(100) not null,
    created_at datetime not null,
    since_id int(11) unsigned default 0,
    screen_name varchar(255) default null,
    key idx_lang(lang),
    key idx_username(username)
) engine=innodb default charset utf8mb4;

create table tweet_topics (
    id int(11) unsigned not null auto_increment primary key,
    lang char(2) not null,
    topic varchar(255) not null ,
    score double default 0.0,
    calculated_date date default null,
    key idx_lang(lang)
) engine=innodb default charset utf8mb4;

create table tweet_pickup_users (
    id int(11) unsigned not null auto_increment primary key,
    screen_name varchar(255) not null,
    lang char(2) not null,
    score double,
    calculated_date date default null,
    key idx_lang(lang)
) engine=innodb default charset utf8mb4;