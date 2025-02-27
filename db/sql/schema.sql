CREATE TABLE Books (
    ID SERIAL PRIMARY KEY,
    ISBN VARCHAR(15) UNIQUE,
    Title VARCHAR(1024),
    Author VARCHAR(512),
    Publisher VARCHAR(255),
    Publication_Year VARCHAR(50),
    Volume VARCHAR(50),
    Image_URL VARCHAR(512),
    Description TEXT,
    Recommendation TEXT,
    Toc TEXT,
    Source VARCHAR(50),
    Url VARCHAR(512),
    Vector_Search BOOLEAN DEFAULT False
);

CREATE TABLE BookEmbedding (
    ID SERIAL PRIMARY KEY,
    ISBN VARCHAR(15) NOT NULL UNIQUE,
    embedding vector (1536),
    FOREIGN KEY (ISBN) REFERENCES Books (ISBN) ON DELETE CASCADE
);

CREATE TABLE Libraries (
    ID SERIAL PRIMARY KEY,
    Lib_Code VARCHAR(20) UNIQUE,
    Lib_Name VARCHAR(100),
    address VARCHAR(255),
    Tel VARCHAR(100),
    Latitude FLOAT,
    Longitude FLOAT,
    Homepage VARCHAR(100),
    Closed VARCHAR(512),
    Operating_Time VARCHAR(512)
);

CREATE TABLE LibsBooks (
    ID SERIAL PRIMARY KEY,
    Lib_Code VARCHAR(20),
    ISBN VARCHAR(15),
    Class_Num VARCHAR(255),
    scrap BOOLEAN,
    FOREIGN KEY (LibCode) REFERENCES Libraries (LibCode) ON DELETE CASCADE,
    FOREIGN KEY (ISBN) REFERENCES Books (ISBN) ON DELETE CASCADE,
    UNIQUE (LibCode, ISBN)
);