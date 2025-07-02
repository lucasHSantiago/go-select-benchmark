CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INTEGER DEFAULT 1
);

INSERT INTO orders (customer_name) VALUES
    ('John Doe'),
    ('Jane Smith');

INSERT INTO order_items (order_id, product_name, price, quantity) VALUES
    (1, 'Laptop', 999.99, 1),
    (1, 'Mouse', 29.99, 1),
    (1, 'Keyboard', 79.99, 1),
    (2, 'Mouse', 29.99, 2),
    (2, 'Keyboard', 79.99, 1);

-- Function to seed the database with sample data
CREATE OR REPLACE FUNCTION seed_database(num_orders INTEGER DEFAULT 100)
RETURNS VOID AS $$
DECLARE
    i INTEGER;
    order_id INTEGER;
    product_names TEXT[] := ARRAY['Laptop', 'Mouse', 'Keyboard', 'Monitor', 'Webcam', 'Speaker', 'Headphones', 'Tablet', 'Phone', 'Charger'];
    customer_names TEXT[] := ARRAY['John Doe', 'Jane Smith', 'Bob Johnson', 'Alice Brown', 'Charlie Wilson', 'Diana Davis', 'Eve Miller', 'Frank Garcia', 'Grace Lee', 'Henry Martinez'];
    num_items INTEGER;
    j INTEGER;
BEGIN
    -- Clear existing data
    DELETE FROM order_items;
    DELETE FROM orders;
    
    -- Reset sequences
    ALTER SEQUENCE orders_id_seq RESTART WITH 1;
    ALTER SEQUENCE order_items_id_seq RESTART WITH 1;
    
    -- Insert orders
    FOR i IN 1..num_orders LOOP
        INSERT INTO orders (customer_name) 
        VALUES (customer_names[1 + (i % array_length(customer_names, 1))])
        RETURNING id INTO order_id;
        
        -- Each order gets exactly 5 items
        num_items := 5;
        
        FOR j IN 1..num_items LOOP
            INSERT INTO order_items (order_id, product_name, price, quantity)
            VALUES (
                order_id,
                product_names[1 + (random() * (array_length(product_names, 1) - 1))::INTEGER],
                (10 + random() * 990)::DECIMAL(10,2), -- Random price between 10-1000
                1 + (random() * 3)::INTEGER -- Random quantity 1-4
            );
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Database seeded with % orders and % order items', 
        num_orders, 
        (SELECT COUNT(*) FROM order_items);
END;
$$ LANGUAGE plpgsql;

-- Call the seed function with default data
SELECT seed_database(50000);