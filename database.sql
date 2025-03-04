-- Create sex type
CREATE TYPE sex AS ENUM ('male', 'female');

CREATE TABLE staff
(
    id            INTEGER PRIMARY KEY NOT NULL,
    password      TEXT                NOT NULL,
    department_id SMALLSERIAL         NOT NULL
);

CREATE TABLE students
(
    id                    INTEGER PRIMARY KEY NOT NULL,
    password              TEXT                NOT NULL,
    enrollment_start_time TIMESTAMPTZ         NOT NULL,
    max_units             SMALLINT            NOT NULL,
    remaining_actions     SMALLINT            NOT NULL,
    department_id         SMALLSERIAL         NOT NULL,
    entry_year            SMALLINT            NOT NULL,
    gender                sex                 NOT NULL
);

CREATE TABLE departments
(
    id   SMALLSERIAL PRIMARY KEY NOT NULL,
    name TEXT                    NOT NULL
);

CREATE TABLE courses
(
    course_id        INTEGER     NOT NULL,
    group_id         INTEGER     NOT NULL,
    for_department   SMALLSERIAL NOT NULL,
    name             TEXT        NOT NULL,
    lecturer         TEXT        NOT NULL,
    units            SMALLINT    NOT NULL,
    capacity         INTEGER     NOT NULL,
    reserve_capacity INTEGER     NOT NULL,
    exam_time        TIMESTAMPTZ,
    class_time       INTEGER     NOT NULL,
    sex_lock         sex,
    notes            TEXT        NOT NULL,
    PRIMARY KEY (course_id, group_id)
);

CREATE TABLE enrolled_courses
(
    id         SERIAL PRIMARY KEY NOT NULL,
    course_id  INTEGER            NOT NULL,
    group_id   INTEGER            NOT NULL,
    student_id INTEGER            NOT NULL,
    reserved   BOOLEAN            NOT NULL
);

ALTER TABLE staff
    ADD CONSTRAINT staff_department_id_department_id FOREIGN KEY (department_id) REFERENCES departments (id);
ALTER TABLE students
    ADD CONSTRAINT students_department_id_department_id FOREIGN KEY (department_id) REFERENCES departments (id);
ALTER TABLE courses
    ADD CONSTRAINT courses_for_department_department_id FOREIGN KEY (for_department) REFERENCES departments (id);
ALTER TABLE enrolled_courses
    ADD CONSTRAINT enrolled_courses_course_id_courses_course_id FOREIGN KEY (course_id, group_id) REFERENCES courses (course_id, group_id);
ALTER TABLE enrolled_courses
    ADD CONSTRAINT enrolled_courses_student_id_users_id FOREIGN KEY (student_id) REFERENCES students (id);