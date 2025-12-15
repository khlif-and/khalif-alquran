-- Fungsi otomatis untuk mengupdate kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

--SEPARATOR--

-- Tabel Surahs (Data Statis Al-Quran)
CREATE TABLE IF NOT EXISTS surahs (
    number INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    latin_name VARCHAR(100) NOT NULL,
    english_name VARCHAR(100) NOT NULL,
    indonesian_name VARCHAR(100) NOT NULL, -- KOLOM BARU
    revelation_type VARCHAR(20) NOT NULL, -- Meccan / Medinan
    total_ayahs INT NOT NULL
);

--SEPARATOR--

-- Tabel Ayahs (Data Ayat)
CREATE TABLE IF NOT EXISTS ayahs (
    id SERIAL PRIMARY KEY,
    surah_id INT NOT NULL,
    number INT NOT NULL, -- Nomor ayat dalam surat
    number_in_quran INT NOT NULL, -- Nomor ayat global (1-6236)
    text_arabic TEXT NOT NULL,
    text_latin TEXT NOT NULL,
    translation TEXT NOT NULL,
    audio_url TEXT,
    CONSTRAINT fk_surah FOREIGN KEY (surah_id) REFERENCES surahs(number) ON DELETE CASCADE,
    UNIQUE (surah_id, number) -- Mencegah duplikasi ayat dalam satu surat
);

--SEPARATOR--

-- Tabel Bookmarks (Data User - Perlu updated_at)
CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL, -- Bisa UUID atau String ID dari Auth
    surah_id INT NOT NULL,
    ayah_number INT NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_bookmark_surah FOREIGN KEY (surah_id) REFERENCES surahs(number) ON DELETE CASCADE,
    UNIQUE (user_id, surah_id, ayah_number) -- Satu user hanya bisa bookmark satu ayat sekali
);

--SEPARATOR--

-- Trigger: Auto update updated_at untuk tabel Bookmarks
DROP TRIGGER IF EXISTS update_bookmarks_modtime ON bookmarks;

--SEPARATOR--

CREATE TRIGGER update_bookmarks_modtime 
BEFORE UPDATE ON bookmarks 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

--SEPARATOR--

-- Setup Readonly User (Untuk keperluan reporting/debugging aman)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'readonly_user') THEN
        CREATE ROLE readonly_user WITH LOGIN PASSWORD 'readonly_password';
    END IF;
END
$$;

--SEPARATOR--

GRANT CONNECT ON DATABASE postgres TO readonly_user;

--SEPARATOR--

GRANT USAGE ON SCHEMA public TO readonly_user;

--SEPARATOR--

GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;

--SEPARATOR--

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readonly_user;