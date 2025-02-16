INSERT INTO store (slug, title, price)
VALUES ('t-shirt', 'T-Shirt', 80),
       ('cup', 'Cup', 20),
       ('book', 'Book', 50),
       ('pen', 'Pen', 10),
       ('powerbank', 'Powerbank', 200),
       ('hoody', 'Hoody', 300),
       ('umbrella', 'Umbrella', 200),
       ('socks', 'Socks', 10),
       ('wallet', 'Wallet', 50),
       ('pink-hoody', 'Pink Hoody', 500)
ON CONFLICT (slug) DO NOTHING;
