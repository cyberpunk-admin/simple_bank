CREATE TABLE "sessions" (
     "id" uuid PRIMARY KEY,
     "user_name" varchar NOT NULL,
     "refresh_token" varchar NOT NULL,
     "user_agent" varchar NOT NULL,
     "client_ip" varchar NOT NULL,
     "is_blocked" boolean NOT NULL DEFAULT false,
     "expires_at" timestamptz NOT NULL,
     "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_name") REFERENCES "users" ("user_name");
