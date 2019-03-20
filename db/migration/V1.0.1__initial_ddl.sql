create table log
(
  id          bigint                                    not null auto_increment primary key,
  log_type    varchar(20)                               not null,
  content     text                                      null,
  create_time timestamp(3) default current_timestamp(3) not null
) ENGINE = InnoDB default CHARSET = utf8mb4;

create index idx_log_log_type on log (log_type);