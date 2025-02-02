CREATE TABLE actions (
    id SERIAL PRIMARY KEY,             
    action TEXT NOT NULL,              
    object_id TEXT NOT NULL,   
    time TIMESTAMP WITH TIME ZONE NOT NULL
);
