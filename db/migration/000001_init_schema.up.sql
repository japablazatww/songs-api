CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY, -- ID de la canción, recibida desde iTunes
    name VARCHAR(255) NOT NULL, -- Nombre de la canción
    artist VARCHAR(255) NOT NULL, -- Nombre del artista
    duration VARCHAR(10), -- Duración en formato HH:MM o MM:SS
    album VARCHAR(255) NOT NULL, -- Nombre del álbum
    artwork TEXT, -- URL de la portada del álbum
    price VARCHAR(20), -- Precio de la canción
    origin VARCHAR(50) -- Plataforma de origen (por ejemplo, "apple")
);

-- Índice para búsquedas por artista
CREATE INDEX idx_artist ON songs (artist);

-- Índice para búsquedas por álbum
CREATE INDEX idx_album ON songs (album);
