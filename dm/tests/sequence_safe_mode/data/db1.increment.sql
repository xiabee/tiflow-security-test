use sequence_safe_mode_test;
insert into t1 (uid, name) values (10003, 'Buenos Aires');
alter table t1 add column (age int, level int);
alter table t1 add index age(age);
alter table t1 add index level(level);
insert into t1 (uid, name, age) values (10005, 'Buenos Aires', 200);
insert into t2 (uid, name) values (20005, 'Aureliano José');
insert into t1 (uid, name, age) values (10006, 'Buenos Aires', 200);
alter table t2 add column (age int, level int);
alter table t2 add index age(age);
alter table t2 add index level(level);
insert into t1 (uid, name, age) values (10007, 'Buenos Aires', 300);
insert into t2 (uid, name, age) values (20006, 'Colonel Aureliano Buendía', 301);
insert into t1 (uid, name, age) values (10008, 'Buenos Aires', 400);
insert into t2 (uid, name, age) values (20007, 'Colonel Aureliano Buendía', 401);
update t1 set age = 404 where uid = 10005;
