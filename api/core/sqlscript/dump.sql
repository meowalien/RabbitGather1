use core;
# 豋入後即擁有的權限
insert into role (name) value ('login');
# 定義權限
insert into restful_rbac_model (p_type, v0, v1, v2) value ('p', 'login', '/member/login', 'POST');



























