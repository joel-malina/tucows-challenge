CREATE TABLE orders (
                        order_id UUID PRIMARY KEY,
                        customer_id UUID NOT NULL,
                        order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        status VARCHAR(50),
                        total_price DECIMAL(10, 2)
);

CREATE TABLE products (
                          product_id UUID PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price DECIMAL(10, 2) NOT NULL,
                          stock INT NOT NULL
);

CREATE TABLE order_items (
                             item_id UUID PRIMARY KEY,
                             order_id UUID NOT NULL,
                             product_id UUID NOT NULL,
                             quantity INT NOT NULL,
                             price DECIMAL(10, 2),
                             FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
                             FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);
