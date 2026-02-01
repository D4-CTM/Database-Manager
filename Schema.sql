-- This script creates a database schema for an Oracle database.
-- It is designed to be run by the 'user' user, which is created by the
-- docker-compose setup with the password 'user_pass'.

--
-- Table: DEPARTMENTS
--
CREATE TABLE DEPARTMENTS (
    department_id   NUMBER(10) NOT NULL,
    department_name VARCHAR2(50) NOT NULL,
    CONSTRAINT pk_departments PRIMARY KEY (department_id)
);

--
-- Table: EMPLOYEES
--
CREATE TABLE EMPLOYEES (
    employee_id   NUMBER(10) NOT NULL,
    first_name    VARCHAR2(50) NOT NULL,
    last_name     VARCHAR2(50) NOT NULL,
    email         VARCHAR2(100) NOT NULL,
    department_id NUMBER(10),
    CONSTRAINT pk_employees PRIMARY KEY (employee_id),
    CONSTRAINT fk_employees_departments FOREIGN KEY (department_id) REFERENCES DEPARTMENTS(department_id)
);

--
-- Table: SALARIES
--
CREATE TABLE SALARIES (
    salary_id   NUMBER(10) NOT NULL,
    employee_id NUMBER(10) NOT NULL,
    amount      NUMBER(10, 2) NOT NULL,
    from_date   DATE NOT NULL,
    to_date     DATE,
    CONSTRAINT pk_salaries PRIMARY KEY (salary_id),
    CONSTRAINT fk_salaries_employees FOREIGN KEY (employee_id) REFERENCES EMPLOYEES(employee_id)
);

--
-- Indexes
--
CREATE INDEX idx_employees_last_name ON EMPLOYEES(last_name);
CREATE INDEX idx_salaries_employee_id ON SALARIES(employee_id);

--
-- Add some sample data
--

-- Departments
INSERT INTO DEPARTMENTS (department_id, department_name) VALUES (1, 'Human Resources');
INSERT INTO DEPARTMENTS (department_id, department_name) VALUES (2, 'Engineering');
INSERT INTO DEPARTMENTS (department_id, department_name) VALUES (3, 'Sales');

-- Employees
INSERT INTO EMPLOYEES (employee_id, first_name, last_name, email, department_id) VALUES (1, 'John', 'Doe', 'john.doe@example.com', 2);
INSERT INTO EMPLOYEES (employee_id, first_name, last_name, email, department_id) VALUES (2, 'Jane', 'Smith', 'jane.smith@example.com', 2);
INSERT INTO EMPLOYEES (employee_id, first_name, last_name, email, department_id) VALUES (3, 'Peter', 'Jones', 'peter.jones@example.com', 3);

-- Salaries
INSERT INTO SALARIES (salary_id, employee_id, amount, from_date, to_date) VALUES (1, 1, 80000, TO_DATE('2023-01-01', 'YYYY-MM-DD'), NULL);
INSERT INTO SALARIES (salary_id, employee_id, amount, from_date, to_date) VALUES (2, 2, 85000, TO_DATE('2023-01-01', 'YYYY-MM-DD'), NULL);
INSERT INTO SALARIES (salary_id, employee_id, amount, from_date, to_date) VALUES (3, 3, 70000, TO_DATE('2023-01-01', 'YYYY-MM-DD'), NULL);

COMMIT;

