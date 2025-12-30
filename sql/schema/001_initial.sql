-- +goose Up
CREATE TYPE item_type AS ENUM ('PRODUCT', 'ASSET');
CREATE TYPE country AS ENUM ( 'AF', 'AX', 'AL', 'DZ', 'AS', 'AD', 'AO', 'AI',
'AQ', 'AG', 'AR', 'AM', 'AW', 'AU', 'AT', 'AZ', 'BS', 'BH', 'BD', 'BB', 'BY',
'BE', 'BZ', 'BJ', 'BM', 'BT', 'BO', 'BQ', 'BA', 'BW', 'BV', 'BR', 'IO', 'BN',
'BG', 'BF', 'BI', 'CV', 'KH', 'CM', 'CA', 'KY', 'CF', 'TD', 'CL', 'CN', 'CX',
'CC', 'CO', 'KM', 'CD', 'CG', 'CK', 'CR', 'CI', 'HR', 'CU', 'CW', 'CY', 'CZ',
'DK', 'DJ', 'DM', 'DO', 'EC', 'EG', 'SV', 'GQ', 'ER', 'EE', 'SZ', 'ET', 'FK',
'FO', 'FJ', 'FI', 'FR', 'GF', 'PF', 'TF', 'GA', 'GM', 'GE', 'DE', 'GH', 'GI',
'GR', 'GL', 'GD', 'GP', 'GU', 'GT', 'GG', 'GN', 'GW', 'GY', 'HT', 'HM', 'VA',
'HN', 'HK', 'HU', 'IS', 'IN', 'ID', 'IR', 'IQ', 'IE', 'IM', 'IL', 'IT', 'JM',
'JP', 'JE', 'JO', 'KZ', 'KE', 'KI', 'KP', 'KR', 'KW', 'KG', 'LA', 'LV', 'LB',
'LS', 'LR', 'LY', 'LI', 'LT', 'LU', 'MO', 'MK', 'MG', 'MW', 'MY', 'MV', 'ML',
'MT', 'MH', 'MQ', 'MR', 'MU', 'YT', 'MX', 'FM', 'MD', 'MC', 'MN', 'ME', 'MS',
'MA', 'MZ', 'MM', 'NA', 'NR', 'NP', 'NL', 'NC', 'NZ', 'NI', 'NE', 'NG', 'NU',
'NF', 'MP', 'NO', 'OM', 'PK', 'PW', 'PS', 'PA', 'PG', 'PY', 'PE', 'PH', 'PN',
'PL', 'PT', 'PR', 'QA', 'RE', 'RO', 'RU', 'RW', 'BL', 'SH', 'KN', 'LC', 'MF',
'PM', 'VC', 'WS', 'SM', 'ST', 'SA', 'SN', 'RS', 'SC', 'SL', 'SG', 'SX', 'SK',
'SI', 'SB', 'SO', 'ZA', 'GS', 'SS', 'ES', 'LK', 'SD', 'SR', 'SJ', 'SE', 'CH',
'SY', 'TW', 'TJ', 'TZ', 'TH', 'TL', 'TG', 'TK', 'TO', 'TT', 'TN', 'TR', 'TM',
'TC', 'TV', 'UG', 'UA', 'AE', 'GB', 'UM', 'US', 'UY', 'UZ', 'VU', 'VE', 'VN',
'VG', 'VI', 'WF', 'EH', 'YE', 'ZM', 'ZW' );
CREATE TYPE currency AS ENUM ( 'AED', 'AFN', 'ALL', 'AMD', 'AOA', 'ARS', 'AUD',
'AWG', 'AZN', 'BAM', 'BBD', 'BDT', 'BGN', 'BHD', 'BIF', 'BMD', 'BND', 'BOB',
'BOV', 'BRL', 'BSD', 'BTN', 'BWP', 'BYN', 'BZD', 'CAD', 'CDF', 'CHE', 'CHF',
'CHW', 'CLF', 'CLP', 'CNY', 'COP', 'COU', 'CRC', 'CUP', 'CVE', 'CZK', 'DJF',
'DKK', 'DOP', 'DZD', 'EGP', 'ERN', 'ETB', 'EUR', 'FJD', 'FKP', 'GBP', 'GEL',
'GHS', 'GIP', 'GMD', 'GNF', 'GTQ', 'GYD', 'HKD', 'HNL', 'HTG', 'HUF', 'IDR',
'ILS', 'INR', 'IQD', 'IRR', 'ISK', 'JMD', 'JOD', 'JPY', 'KES', 'KGS', 'KHR',
'KMF', 'KPW', 'KRW', 'KWD', 'KYD', 'KZT', 'LAK', 'LBP', 'LKR', 'LRD', 'LSL',
'LYD', 'MAD', 'MDL', 'MGA', 'MKD', 'MMK', 'MNT', 'MOP', 'MRU', 'MUR', 'MVR',
'MWK', 'MXN', 'MXV', 'MYR', 'MZN', 'NAD', 'NGN', 'NIO', 'NOK', 'NPR', 'NZD',
'OMR', 'PAB', 'PEN', 'PGK', 'PHP', 'PKR', 'PLN', 'PYG', 'QAR', 'RON', 'RSD',
'RUB', 'RWF', 'SAR', 'SBD', 'SCR', 'SDG', 'SEK', 'SGD', 'SHP', 'SLE', 'SOS',
'SRD', 'SSP', 'STN', 'SVC', 'SYP', 'SZL', 'THB', 'TJS', 'TMT', 'TND', 'TOP',
'TRY', 'TTD', 'TWD', 'TZS', 'UAH', 'UGX', 'USD', 'USN', 'UYI', 'UYU', 'UYW',
'UZS', 'VED', 'VES', 'VND', 'VUV', 'WST', 'XAD', 'XAF', 'XAG', 'XAU', 'XBA',
'XBB', 'XBC', 'XBD', 'XCD', 'XCG', 'XDR', 'XOF', 'XPD', 'XPF', 'XPT', 'XSU',
'XTS', 'XUA', 'XXX', 'YER', 'ZAR', 'ZMW', 'ZWG' );

CREATE TABLE users 
(
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,
    email           TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    name            TEXT
);

CREATE TABLE accounts 
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    country    country NOT NULL,
    deleted    BOOLEAN NOT NULL DEFAULT FALSE,
    nickname   TEXT,
    owner_id   UUID DEFAULT NULL REFERENCES users ON DELETE SET NULL
);

CREATE TABLE accounts_users 
(
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (account_id, user_id)
);

CREATE TABLE groups
(
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    description TEXT,
    name        TEXT NOT NULL,
    account_id  UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    parent_id   UUID DEFAULT NULL REFERENCES groups ON DELETE SET NULL
);

CREATE TABLE inventories 
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    in_stock   INTEGER NOT NULL DEFAULT 0,
    orderable  INTEGER,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE
);

CREATE TABLE items 
(
    id             UUID PRIMARY KEY,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,
    active         BOOLEAN NOT NULL DEFAULT TRUE,
    description    TEXT DEFAULT NULL,
    name           TEXT NOT NULL,
    price_amount   INTEGER DEFAULT NULL,
    price_currency currency DEFAULT NULL,
    type           item_type NOT NULL,
    variant        BOOLEAN NOT NULL DEFAULT FALSE,
    account_id     UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    group_id       UUID DEFAULT NULL REFERENCES groups ON DELETE SET NULL,
    inventory_id   UUID DEFAULT NULL REFERENCES inventories ON DELETE SET NULL
);

CREATE TABLE item_identifiers
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    ean        TEXT CHECK (char_length(ean) = 8 OR (char_length(ean) >= 12 AND char_length(ean) <= 14)),
    gtin       TEXT CHECK (char_length(gtin) = 8 OR (char_length(gtin) >= 12 AND char_length(gtin) <= 14)),
    isbn       TEXT CHECK (char_length(isbn) = 10 OR char_length(isbn) = 13),
    jan        TEXT CHECK (char_length(jan) = 8 OR char_length(jan) = 13),
    mpn        TEXT,
    nsn        TEXT CHECK (char_length(nsn) = 13),
    upc        TEXT CHECK (char_length(upc) = 12),
    qr         TEXT,
    sku        TEXT,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    item_id    UUID UNIQUE NOT NULL REFERENCES items ON DELETE CASCADE
);

CREATE TABLE item_variant_attributes 
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name       TEXT NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    item_id    UUID NOT NULL REFERENCES items ON DELETE CASCADE
);

CREATE TABLE item_variant_attribute_options
(
    id                        UUID PRIMARY KEY,
    created_at                TIMESTAMP NOT NULL,
    updated_at                TIMESTAMP NOT NULL,
    name                      TEXT NOT NULL,
    account_id                UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    item_variant_attribute_id UUID NOT NULL REFERENCES item_variant_attributes ON DELETE CASCADE
);

CREATE TABLE api_keys
(
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,
    expires_at      TIMESTAMP DEFAULT NULL,
    name            TEXT NOT NULL,
    note            TEXT,
    redacted_secret TEXT NOT NULL,
    secret          TEXT NOT NULL,
    account_id      UUID NOT NULL REFERENCES accounts ON DELETE CASCADE
);

-- +goose Down
DROP TABLE  api_keys;
DROP TABLE item_variant_attribute_options;
DROP TABLE item_variant_attributes;
DROP TABLE item_identifiers;
DROP TABLE items;
DROP TABLE inventories;
DROP TABLE groups;
DROP TABLE accounts_users;
DROP TABLE accounts;
DROP TABLE users;
DROP TYPE currency;
DROP TYPE country;
DROP TYPE item_type;
