# Buscador de Canciones API

Esta API proporciona un servicio de búsqueda de canciones de diferentes proveedores, diseñada para ser utilizada por aplicaciones. Utiliza autenticación Bearer Token para asegurar el acceso.

## Instalación en Servidor

Para levantar el proyecto en un servidor, sigue estos pasos:

1. **Instala Git:**
```bash
sudo apt-get update
sudo apt-get install git
```

2. **Clona el repositorio:**
```bash
git clone https://github.com/japablazatww/songs-searcher.git
```

3. **Navega a la carpeta del proyecto:**
```bash
cd songs-searcher
```

4. **Haz ejecutable el script de instalación:**
```bash
chmod +x install.sh
```

5. **Ejecuta el script de instalación:**
```bash
sudo ./install.sh
```

## Dependencias para Desarrollo

Asegúrate de tener las siguientes dependencias instaladas para desarrollo local:
* [**Docker**](https://docs.docker.com/engine/install/)
* **Make:**  `sudo apt-get install make`
* [**Go Migrate**](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Levantar el Proyecto

Para levantar el proyecto, ejecuta los siguientes comandos:

```bash
make up_build
make migrate_db_up
```

## Documentación de la API

Esta API es un buscador de canciones que permite buscar por artista, canción y álbum. Requiere autenticación Bearer Token para su uso.

### Autenticación

1. **Registro de la Aplicación:**
   
   Para crear una aplicación, envía una solicitud POST con el nombre de la aplicación al endpoint `/auth/register/`:

```json
{
    "app_name": "TestingApp"
}
```

Este endpoint devolverá un `client_id` y un `client_secret`.

2. **Obtención del Token:**
   
   Envía una solicitud POST con el `client_id` y `client_secret` obtenidos al endpoint `/auth/token`:

```json
{
    "client_id": "el_client_id_obtenido",
    "client_secret": "el_client_secret_obtenido"
}
```

Este endpoint devolverá el token que deberás incluir en el header Authorization como un Bearer Token para poder realizar búsquedas.

### Endpoint de Búsqueda

El endpoint `/search/` requiere un Bearer Token en el header Authorization.

Acepta los siguientes parámetros de consulta (queries):
* `artist`: Nombre del artista
* `song`: Nombre de la canción
* `album`: Nombre del álbum

**Ejemplo:**
```
/search?artist=eminem&song=stan&album=The Marshall Mathers LP
```

Se requiere al menos un parámetro para realizar la búsqueda.

### Respuesta de la API

La API devuelve una respuesta en formato JSON con la siguiente estructura:

```json
{
    "totalSongs": 40,
    "songs": [
        {
            "id": 731939599,
            "name": "Stan",
            "artist": "Eminem",
            "duration": "06:45",
            "album": "Greatest Hits (Deluxe Edition)",
            "artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music125/v4/de/6d/3f/de6d3ff9-13c6-d70c-9010-0709e3780a84/886444202091.jpg/100x100bb.jpg",
            "price": "USD -1.00",
            "origin": "apple"
        }
    ]
}
```

### Algoritmo de Ordenamiento por coincidencia

La API utiliza un algoritmo de ordenamiento por coincidencia. Los resultados que coincidan más con la búsqueda aparecerán primero en la lista de canciones.