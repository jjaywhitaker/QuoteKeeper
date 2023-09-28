create table if not exists quotes (
id serial primary key, 
body varchar(1000) not null, 
author varchar(255),
created_date timestamp not null 
);

create table if not exists categories (
id serial primary key, 
category varchar (255) not null,
create_date timestamp not null 
);

create table if not exists quotes_category (
quote_id int references quotes, 
category_id int references categories, 
primary key (quote_id, category_id)
);