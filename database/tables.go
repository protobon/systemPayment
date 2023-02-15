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

// PayerTableCreate - Payer
var PayerTableCreate = `CREATE TABLE IF NOT EXISTS payer
(
   id SERIAL,
   name VARCHAR(100) NOT NULL,
   email VARCHAR(100) NOT NULL,
   birth_date VARCHAR(10) NOT NULL,
   phone VARCHAR(20) NOT NULL,
   document VARCHAR(30) NOT NULL,
   user_reference VARCHAR(50) NOT NULL,
   created_at DATE NOT NULL,
   updated_at DATE,
   CONSTRAINT payer_pkey PRIMARY KEY (id)
)`

// AddressTableCreate - Address
var AddressTableCreate = `CREATE TABLE IF NOT EXISTS address
(
   id SERIAL,
   payer INT NOT NULL,
   state VARCHAR(100) NOT NULL,
   city VARCHAR(100) NOT NULL,
   zip_code VARCHAR(10) NOT NULL,
   street VARCHAR(100) NOT NULL,
   number VARCHAR(20) NOT NULL,
   created_at DATE NOT NULL,
   CONSTRAINT address_pkey PRIMARY KEY (id),
   CONSTRAINT fk_payer FOREIGN KEY(payer) REFERENCES payer(id)
)`

// SecureCardTableCreate - Address
var SecureCardTableCreate = `CREATE TABLE IF NOT EXISTS secure_card
(
   id SERIAL,
   payer INT NOT NULL,
   token VARCHAR(100) NOT NULL,
   created_at DATE NOT NULL,
   CONSTRAINT secure_card_pkey PRIMARY KEY (id),
   CONSTRAINT fk_payer FOREIGN KEY(payer) REFERENCES payer(id)
)`
