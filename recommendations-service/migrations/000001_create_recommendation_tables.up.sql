CREATE TABLE users (
    user_uid TEXT PRIMARY KEY,    
    interests TEXT[]                     
);

CREATE TABLE products (
    product_id BIGSERIAL PRIMARY KEY,
    tags TEXT[],                                  
    popularity BIGINT DEFAULT 0                          
);

CREATE TABLE recommendations (
    id BIGSERIAL PRIMARY KEY,  
    user_uid TEXT NOT NULL,        
    product_id BIGINT NOT NULL,         
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE, 
    FOREIGN KEY (user_uid) REFERENCES users(user_uid) ON DELETE CASCADE         
);

