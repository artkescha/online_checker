CREATE TABLE users (
    id serial NOT NULL,
    name text NOT NULL,
	  password text NOT NULL,
    role_id integer DEFAULT 2,

  PRIMARY KEY(id),
	CONSTRAINT constraint_user_name UNIQUE (name),
  FOREIGN KEY (role_id) REFERENCES roles (Id) ON DELETE RESTRICT ON UPDATE CASCADE
);


CREATE TABLE roles (
    id serial NOT NULL,
    name text NOT NULL,
    PRIMARY KEY(id),
	  CONSTRAINT constraint_name1 UNIQUE (name)
);

select * from  prog_languages

UPDATE users SET role_id = 1 WHERE id = 82;

INSERT INTO roles(name) VALUES('user');


CREATE TABLE prog_languages (
    id serial NOT NULL,
    name varchar(200) NOT NULL,
    PRIMARY KEY(id),
	  CONSTRAINT constraint_lang_name UNIQUE (name)
);


INSERT INTO prog_languages(name) VALUES('golang 1.13');


CREATE TABLE tasks (
    id serial NOT NULL,
    title varchar(200) NOT NULL,
    description text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    test_path varchar(200) NOT NULL,
    PRIMARY KEY(id),
	  CONSTRAINT constraint_task_title UNIQUE (title),
    CONSTRAINT constraint_task_description UNIQUE (description)
);


INSERT INTO tasks(title, body, prog_lang_id, test_path) VALUES('task2', 'task2', 1, 'c:/');

select * from tasks


SELECT CURRENT_TIMESTAMP
SELECT CURRENT_TIMESTAMP(2);

SELECT CURRENT_TIME

SELECT current_time;

select extract(epoch from now());




/* helpers */


CREATE TABLE users (
    id serial NOT NULL,
    name text NOT NULL,
	  password text NOT NULL,
    role_id integer DEFAULT 2,

  PRIMARY KEY(id),
	CONSTRAINT constraint_user_name UNIQUE (name),
  FOREIGN KEY (role_id) REFERENCES roles (Id) ON DELETE RESTRICT ON UPDATE CASCADE
);


CREATE TABLE roles (
    id serial NOT NULL,
    name text NOT NULL,
    PRIMARY KEY(id),
	  CONSTRAINT constraint_name1 UNIQUE (name)
);

select * from  users

UPDATE users SET role_id = 1 WHERE id = 82;

INSERT INTO roles(name) VALUES('user');


CREATE TABLE prog_languages (
    id serial NOT NULL,
    name varchar(200) NOT NULL,
    PRIMARY KEY(id),
	  CONSTRAINT constraint_lang_name UNIQUE (name)
);


INSERT INTO prog_languages(name) VALUES('golang 1.13');


CREATE TABLE tasks (
    id serial NOT NULL,
    title varchar(200) NOT NULL,
    description text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    test_path varchar(200) NOT NULL,
    PRIMARY KEY(id),
	  CONSTRAINT constraint_task_title UNIQUE (title),
    CONSTRAINT constraint_task_description UNIQUE (description)
);


INSERT INTO tasks(title, body, prog_lang_id, test_path) VALUES('task2', 'task2', 1, 'c:/');

select * from tasks


SELECT * FROM tasks ORDER BY created_at DESC LIMIT 4 OFFSET 0


SELECT CURRENT_TIMESTAMP
SELECT CURRENT_TIMESTAMP(2);

SELECT CURRENT_TIME

SELECT current_time;

select extract(epoch from now());


select * from users



select * from attempts

select * from prog_languages

CREATE TABLE attempts (
    id serial NOT NULL,
    user_id integer NOT NULL,
    solution text NOT NULL,
    status smallint NOT NULL,
    description text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    task_id integer NOT NULL,
    language_id integer NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES users (Id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks (Id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (language_id) REFERENCES prog_languages (Id) ON DELETE RESTRICT ON UPDATE CASCADE
);



DROP TABLE attempts

SELECT * FROM attempts ORDER BY created_at DESC LIMIT 100 OFFSET 2


select * from users


DELETE FROM tasks where id = 2;

ALTER table attempts
ADD CONSTRAINT fk_name
  FOREIGN KEY (task_id)
  REFERENCES tasks (id)
  ON DELETE CASCADE;

REFERENCES parent_table_name(parent_column_name)
ON DELETE CASCADE;
ADD CONSTRAINT constaraint_task_id REFERENCES tasks (id) ON DELETE CASCADE