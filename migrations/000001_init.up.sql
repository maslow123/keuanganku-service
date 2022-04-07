CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "email" varchar(45) NOT NULL,
  "password" varchar(255) NOT NULL,
  "photo" varchar(100) DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "pos" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "name" varchar(255) NOT NULL,
  "type" int NOT NULL,
  "total" int DEFAULT 0,
  "color" varchar(10) DEFAULT '#FFFFFF',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
); 

CREATE TABLE "transactions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "pos_id" int NOT NULL,
  "total" int NOT NULL,
  "details" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "transactions" ADD FOREIGN KEY ("pos_id") REFERENCES "pos" ("id") ON DELETE CASCADE;
ALTER TABLE "pos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE INDEX ON "transactions" ("user_id", "pos_id");
CREATE INDEX ON "pos" ("user_id");

INSERT INTO users (name, email, password) 
VALUES
('user 1', 'user1@gmail.com', '$2a$05$wQ8lYAdEw7ZzF3OSzWeCKee8wc0KWxbBqfJpNu.lb.f1rvuSyy/I2');