CREATE TABLE IF NOT EXISTS courses (
    id            BIGSERIAL PRIMARY KEY,
    title         TEXT NOT NULL,
    slug          TEXT NOT NULL UNIQUE,
    description   TEXT NOT NULL DEFAULT '',
    price_paise   BIGINT NOT NULL CHECK (price_paise > 0),
    whatsapp_link TEXT NOT NULL,
    pdf_s3_key    TEXT NOT NULL,
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS orders (
    id                  BIGSERIAL PRIMARY KEY,
    course_id           BIGINT NOT NULL REFERENCES courses(id),
    buyer_name          TEXT NOT NULL,
    buyer_email         TEXT NOT NULL,
    razorpay_order_id   TEXT NOT NULL UNIQUE,
    razorpay_payment_id TEXT,
    amount_paise        BIGINT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'CREATED', -- CREATED, PAID
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS enrollments (
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT NOT NULL UNIQUE REFERENCES orders(id),
    course_id    BIGINT NOT NULL REFERENCES courses(id),
    buyer_name   TEXT NOT NULL,
    buyer_email  TEXT NOT NULL,
    access_token TEXT NOT NULL UNIQUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reviews (
    id            BIGSERIAL PRIMARY KEY,
    enrollment_id BIGINT NOT NULL UNIQUE REFERENCES enrollments(id),
    course_id     BIGINT NOT NULL REFERENCES courses(id),
    rating        INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment       TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reviews_course ON reviews(course_id);

-- test course (price paise mein: 49900 = Rs 499)
INSERT INTO courses (title, slug, description, price_paise, whatsapp_link, pdf_s3_key)
VALUES ('Java Backend Notes', 'java-backend', 'Complete notes + live classes',
        49900, 'https://chat.whatsapp.com/XXXX', 'notes/java-backend.pdf')
ON CONFLICT (slug) DO NOTHING;