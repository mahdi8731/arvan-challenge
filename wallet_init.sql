CREATE TABLE wallet (
    wallet_id UUID PRIMARY KEY,
   	phone_number varchar(15) UNIQUE,
    last_modified timestamp,
    inventory numeric(10,1)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    description varchar(256),
	date timestamp,
	amount numeric(10,1),
    wallet_id UUID REFERENCES wallet(wallet_id)
);