create table if not exists fidus_master(
	ID varchar(36) primary key,
	name varchar(100) not null,
	email varchar(100) unique not null,
	password_hash varchar(100)
);

create table if not exists fidus_password(
	key varchar(50) not null,
	master_id varchar(50) not null,
	password varchar(100)
);
