-- migrate:up
CREATE TABLE users (
    id UUID NOT NULL,
    email TEXT NOT NULL,
    token TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    bio TEXT,
    image TEXT,

    CONSTRAINT pk$users PRIMARY KEY (id),
    CONSTRAINT unq$email UNIQUE (email),
    CONSTRAINT unq$username UNIQUE (username)
);

-- migrate:down

