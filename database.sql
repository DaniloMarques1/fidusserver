create table if not exists fidus_master(
	ID varchar(36) primary key,
	name varchar(100) not null,
	email varchar(100) unique not null,
	password_hash varchar(100),
	created_at TIMESTAMP NOT NULL
);

create table if not exists fidus_password(
	key varchar(50) not null,
	master_id varchar(50) not null,
	password bytea,
	primary key(key, master_id),
	created_at TIMESTAMP NOT NULL,
	constraint fk_master foreign key(master_id) references fidus_master(ID)
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- this should be executed inside the database server
-- ALTER TABLE fidus_master ADD created_at TIMESTAMP NOT NULL DEFAULT NOW(); 
-- ALTER TABLE fidus_password ADD created_at TIMESTAMP NOT NULL DEFAULT NOW(); 

-- expiration for the master password
-- ALTER TABLE fidus_master ADD password_expiration_date TIMESTAMP; 
-- UPDATE fidus_master SET password_expiration_date = '2024-08-25 02:32:22.978876';
-- ALTER TABLE fidus_master ALTER COLUMN password_expiration_date SET NOT NULL;
