CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price INT NOT NULL CHECK (price >= 0),
    seller_name TEXT NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO categories (id, name)
VALUES 
  ('11111111-1111-1111-1111-111111111111', 'Смартфоны'),
  ('22222222-2222-2222-2222-222222222222', 'Бытовая техника'),
  ('33333333-3333-3333-3333-333333333333', 'Мужская одежда'),
  ('44444444-4444-4444-4444-444444444444', 'Электроника');

INSERT INTO products (id, name, description, price, seller_name, category_id)
VALUES 
  (
    gen_random_uuid(),
    'iPhone 14 Pro',
    'Смартфон от Apple с передовой камерой и процессором A16',
    130000,
    'TechWorld',
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    gen_random_uuid(),
    'Samsung Galaxy S23',
    'Флагманский смартфон с лучшим дисплеем и камерой',
    120000,
    'TechWorld',
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    gen_random_uuid(),
    'Стиральная машина LG',
    'Фронтальная загрузка, 7 кг, эффективное удаление пятен',
    40000,
    'HomeStyle',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    gen_random_uuid(),
    'Пылесос Dyson',
    'Вертикальный беспроводной пылесос с циклонной технологией',
    35000,
    'HomeStyle',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    gen_random_uuid(),
    'Кофеварка DeLonghi',
    'Эспрессо-машина с капучинатором и автоматическим отключением',
    28000,
    'HomeStyle',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    gen_random_uuid(),
    'Мужская куртка The North Face',
    'Тёплая и ветронепроницаемая куртка для холодного сезона',
    16000,
    'FashionX',
    '33333333-3333-3333-3333-333333333333'
  ),
  (
    gen_random_uuid(),
    'Футболка Nike',
    'Классическая мужская футболка из дышащего материала',
    2900,
    'FashionX',
    '33333333-3333-3333-3333-333333333333'
  ),
  (
    gen_random_uuid(),
    'Laptop Pro 15',
    'Высокопроизводительный ноутбук для профессионалов',
    120000,
    'TechWorld',
    '44444444-4444-4444-4444-444444444444'
  ),
  (
    gen_random_uuid(),
    'Gaming Mouse X500',
    'Эргономичная игровая мышь с RGB-подсветкой',
    4500,
    'ClickZone',
    '44444444-4444-4444-4444-444444444444'
  ),
  (
    gen_random_uuid(),
    'Smartphone Ultra 12',
    'Новейший смартфон с OLED-экраном и быстрой зарядкой',
    85000,
    'LaptopHub',
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    gen_random_uuid(),
    'Wireless Keyboard K7',
    'Компактная беспроводная клавиатура с длительным сроком работы',
    3800,
    'KeyMasters',
    '44444444-4444-4444-4444-444444444444'
  ),
  (
    gen_random_uuid(),
    'Noise Cancelling Headphones',
    'Наушники с активным шумоподавлением и высоким качеством звука',
    15000,
    'AudioShop',
    '44444444-4444-4444-4444-444444444444'
  ),
  (
    gen_random_uuid(),
    'Smartwatch FitX',
    'Фитнес-часы с отслеживанием сна и пульсометром',
    12000,
    'SportStore',
    '44444444-4444-4444-4444-444444444444'
  ),
  (
    gen_random_uuid(),
    'Men''s Casual Jacket',
    'Стильная мужская куртка для повседневной носки',
    8700,
    'UrbanWear',
    '33333333-3333-3333-3333-333333333333'
  );



