CREATE TABLE people (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        surname VARCHAR(100) NOT NULL,
                        patronymic VARCHAR(100),
                        age INT NOT NULL,
                        nationality VARCHAR(100) NOT NULL,
                        gender VARCHAR(10) NOT NULL
);