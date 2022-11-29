alter table if exists users
drop column created_date;

alter table if exists users
drop column deleted_date;

alter table if exists users
drop column updated_date;