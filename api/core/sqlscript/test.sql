insert into user (type, uuid, account, password_hash ,password_salt, last_login) value (?,?,?,?,?,?);


insert into user_info(user_id , nick_name,email,phone , gender , photo_url)value (?,?,?,?,?,?);



select * from user;
select * from user_info;