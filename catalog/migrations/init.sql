-- CREATE TABLE sellers (
--     id UUID PRIMARY KEY,
--     shop_name TEXT NOT NULL,
--     description TEXT,
--     user_id UUID NOT NULL
-- );

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price INT NOT NULL CHECK (price >= 0),
    seller_name TEXT  NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    comment TEXT NOT NULL
);



-- INSERT INTO sellers (id, shop_name, description, user_id) VALUES
--   ('', 'TechWorld', 'Гаджеты и электроника', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
--   ('', 'HomeStyle', 'Товары для дома', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'),
--   ('', 'FashionX', 'Одежда и аксессуары', 'cccccccc-cccc-cccc-cccc-cccccccccccc');

INSERT INTO categories (id, name, parent_id) VALUES
  ('10000000-0000-0000-0000-000000000001', 'Электроника', NULL),
  ('10000000-0000-0000-0000-000000000002', 'Бытовая техника', NULL),
  ('10000000-0000-0000-0000-000000000003', 'Одежда', NULL),
  ('10000000-0000-0000-0000-000000000004', 'Смартфоны', '10000000-0000-0000-0000-000000000001'),
  ('10000000-0000-0000-0000-000000000005', 'Мужская одежда', '10000000-0000-0000-0000-000000000003');

INSERT INTO products (id, name, description, price, seller_name, category_id, created_at) VALUES
  ('00000000-0000-0000-0000-000000000001', 'iPhone 14 Pro', 'Смартфон от Apple', 130000, 'TechWorld', '10000000-0000-0000-0000-000000000004', NOW()),
  ('00000000-0000-0000-0000-000000000002', 'Samsung Galaxy S23', 'Флагман от Samsung', 120000, 'TechWorld', '10000000-0000-0000-0000-000000000004', NOW()),
  ('00000000-0000-0000-0000-000000000003', 'Стиральная машина LG', 'Фронтальная загрузка, 7 кг', 40000, 'HomeStyle', '10000000-0000-0000-0000-000000000002', NOW()),
  ('00000000-0000-0000-0000-000000000004', 'Пылесос Dyson', 'Вертикальный беспроводной пылесос', 35000, 'HomeStyle', '10000000-0000-0000-0000-000000000002', NOW()),
  ('00000000-0000-0000-0000-000000000005', 'Кофеварка DeLonghi', 'Эспрессо-машина с капучинатором', 28000, 'HomeStyle', '10000000-0000-0000-0000-000000000002', NOW()),
  ('00000000-0000-0000-0000-000000000006', 'Мужская куртка The North Face', 'Тёплая и ветронепроницаемая', 16000, 'FashionX', '10000000-0000-0000-0000-000000000005', NOW()),
  ('00000000-0000-0000-0000-000000000007', 'Футболка Nike', 'Классическая мужская футболка', 2900, 'FashionX', '10000000-0000-0000-0000-000000000005', NOW()),
  ('00000000-0000-0000-0000-000000000008', 'Кроссовки Adidas Ultraboost', 'Удобные беговые кроссовки', 11000, 'FashionX', '10000000-0000-0000-0000-000000000005', NOW()),
  ('00000000-0000-0000-0000-000000000009', 'Ноутбук MacBook Air M2', 'Лёгкий и мощный', 150000, 'TechWorld', '10000000-0000-0000-0000-000000000001', NOW()),
  ('00000000-0000-0000-0000-000000000010', 'Планшет iPad Air', 'С поддержкой Apple Pencil', 75000, 'TechWorld', '10000000-0000-0000-0000-000000000001', NOW());
