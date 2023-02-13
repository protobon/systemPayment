package database

// DummyTableCreate - Dummy
var DummyTableCreate = `CREATE TABLE IF NOT EXISTS dummy
(
   id SERIAL,
   name TEXT NOT NULL,
   price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
   created_at DATE NOT NULL,
   updated_at DATE,
   CONSTRAINT dummy_pkey PRIMARY KEY (id)
)`
