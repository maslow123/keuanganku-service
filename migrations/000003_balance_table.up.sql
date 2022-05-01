CREATE TABLE "balance" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "type" int NOT NULL, -- 0: transfer, 1: cash
  "total" int DEFAULT 0,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  
  UNIQUE(user_id, type)
); 
CREATE INDEX ON "balance" ("id", "user_id");
-- CREATE UNIQUE INDEX balance_note_idx on balance (user_id, type);

INSERT INTO balance 
("user_id", "type", "total")
VALUES
(1, 0, 0), (1, 1, 0);