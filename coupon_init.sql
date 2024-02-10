CREATE TABLE Coupons (
    coupon_id UUID PRIMARY KEY,
    code varchar(50) UNIQUE,
    expire_date timestamp,
    charge_amount NUMERIC,
	allowed_times NUMERIC
);

CREATE TABLE couponsـused (
    id UUID PRIMARY KEY,
    phone_number varchar(15) UNIQUE,
    coupon_id UUID REFERENCES coupons(coupon_id)
);

CREATE TABLE outbox (
    id SERIAL PRIMARY KEY,
    phone_number varchar(15),
    amount NUMERIC
);