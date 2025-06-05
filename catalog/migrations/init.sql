CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price INT NOT NULL CHECK (price >= 0),
    seller_name TEXT NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    comment TEXT NOT NULL
);

INSERT INTO products (id, name, description, price, seller_name, category)
VALUES 
  (
    gen_random_uuid(),
    'iPhone 14 Pro',
    'Смартфон от Apple с передовой камерой и процессором A16',
    130000,
    'TechWorld',
    'Смартфоны'
  ),
  (
    gen_random_uuid(),
    'Samsung Galaxy S23',
    'Флагманский смартфон с лучшим дисплеем и камерой',
    120000,
    'TechWorld',
    'Смартфоны'
  ),
  (
    gen_random_uuid(),
    'Стиральная машина LG',
    'Фронтальная загрузка, 7 кг, эффективное удаление пятен',
    40000,
    'HomeStyle',
    'Бытовая техника'
  ),
  (
    gen_random_uuid(),
    'Пылесос Dyson',
    'Вертикальный беспроводной пылесос с циклонной технологией',
    35000,
    'HomeStyle',
    'Бытовая техника'
  ),
  (
    gen_random_uuid(),
    'Кофеварка DeLonghi',
    'Эспрессо-машина с капучинатором и автоматическим отключением',
    28000,
    'HomeStyle',
    'Бытовая техника'
  ),
  (
    gen_random_uuid(),
    'Мужская куртка The North Face',
    'Тёплая и ветронепроницаемая куртка для холодного сезона',
    16000,
    'FashionX',
    'Мужская одежда'
  ),
  (
    gen_random_uuid(),
    'Футболка Nike',
    'Классическая мужская футболка из дышащего материала',
    2900,
    'FashionX',
    'Мужская одежда'
  ),
  (
    gen_random_uuid(),
    'Laptop Pro 15',
    'Высокопроизводительный ноутбук для профессионалов',
    120000,
    'TechWorld',
    'Электроника'
  ),
  (
    gen_random_uuid(),
    'Gaming Mouse X500',
    'Эргономичная игровая мышь с RGB-подсветкой',
    4500,
    'ClickZone',
    'Электроника'
  ),
  (
    gen_random_uuid(),
    'Smartphone Ultra 12',
    'Новейший смартфон с OLED-экраном и быстрой зарядкой',
    85000,
    'LaptopHub',
    'Смартфоны'
  ),
  (
    gen_random_uuid(),
    'Wireless Keyboard K7',
    'Компактная беспроводная клавиатура с длительным сроком работы',
    3800,
    'KeyMasters',
    'Электроника'
  ),
  (
    gen_random_uuid(),
    'Noise Cancelling Headphones',
    'Наушники с активным шумоподавлением и высоким качеством звука',
    15000,
    'AudioShop',
    'Электроника'
  ),
  (
    gen_random_uuid(),
    'Smartwatch FitX',
    'Фитнес-часы с отслеживанием сна и пульсометром',
    12000,
    'SportStore',
    'Электроника'
  ),
  (
    gen_random_uuid(),
    'Men''s Casual Jacket',
    'Стильная мужская куртка для повседневной носки',
    8700,
    'UrbanWear',
    'Мужская одежда'
  );



