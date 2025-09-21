ALTER TABLE charging_applications
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
-- пользователи
INSERT INTO users (login, password, is_moderator)
VALUES
  ('admin', 'admin', true),
  ('client1', 'client1', false),
  ('client2', 'client2', false);

-- тарифы зарядки
INSERT INTO charging_tariffs (nameof_tariff, description, image_url, price_per_hour, power, is_deleted)
VALUES
  ('Быстрая зарядка DC (будни)', 'Зарядка постоянным током 50-150 кВт', 'http://127.0.0.1:9000/charging-images/image.png', 12.0, 150.0,  false),
  ('Быстрая зарядка DC (выходные)', 'Зарядка постоянным током 50-150 кВт', 'http://127.0.0.1:9000/charging-images/image.png', 15.0, 150.0, false),
  ('AC зарядка Level 2 (будни)', 'Зарядка переменным током 22 кВт', 'http://127.0.0.1:9000/charging-images/image.png', 8.0, 22.0, false),
  ('AC зарядка Level 2 (выходные)', 'Зарядка переменным током 22 кВт', 'http://127.0.0.1:9000/charging-images/image.png', 10.0, 22.0, false),
  ('Медленная зарядка Level 1', 'Домашняя зарядка 3.7-7.4 кВт', 'http://127.0.0.1:9000/charging-images/image.png', 5.0, 7.4, false),
  ('Ультрабыстрая зарядка DC', 'Зарядка 350 кВт (Tesla Supercharger)', 'http://127.0.0.1:9000/charging-images/image.png', 18.0, 350.0, false);

-- заявки
INSERT INTO charging_applications (total_price, creator_id, amount_of_orders, status)
VALUES
  (648.0, 3, 2, 'canceled'),
  (216.0, 2, 1, 'deleted'),
  (0.0, 3, 1, 'draft');


-- заказы (связь M-M)
INSERT INTO charging_orders (application_id, tariff_id, battery_capacity, current_percent, start_time, estimated_time, calculated_price)
VALUES
  (1, 1, 75.5, 20, '2025-02-22 14:30:00', 2.5, 300.0),
  (2, 6, 100.0, 10, '2025-01-23 15:00:00', 1.0, 180.0),
  (2, 2, 85.0, 40, '2025-01-23 15:30:00', 1.5, 468.0),
  (3, 1, 65.0, 25, '2025-01-24 16:10:00', 1.5, 216.0);
