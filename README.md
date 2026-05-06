# Flatter

Herramienta para copiar imágenes de múltiples carpetas a una carpeta destino.

## Instalación

```bash
go build -o flatter .
```

## Uso

```bash
./flatter [flags] <destino> <modo> <carpeta1> [<carpeta2>...]
# o sin compilar:
# go run main.go [flags] <destino> <modo> <carpeta1> [<carpeta2>...]
```

### Argumentos

| Argumento | Descripción |
|-----------|-------------|
| `destino` | Ruta de la carpeta destino |
| `modo` | Modo de copia: `copy` (renombra si existe) o `skip` (omite si existe) |
| `carpeta1`... | Carpeta/s fuente a escanear |

### Flags

| Flag | Descripción | Valor por defecto |
|------|-------------|-------------------|
| `-w` | Número de workers concurrentes | 8 |
| `-n` | Ignorar capturas de pantalla | false |
| `-i` | Formatos a ignorar (separados por coma, ej: webp,jpg) | (ninguno) |
| `-a` | Formatos adicionales a incluir (separados por coma, ej: mp4,pdf) | (ninguno) |

## Ejemplos

Copiar todas las imágenes de una carpeta (incluyendo subcarpetas):

```bash
./flatter /home/user/fotos copy /run/media/alexis/6BBD-E38E/Google\ Fotos
```
```bash
./flatter /home/user/fotos copy "/run/media/alexis/6BBD-E38E/Google Fotos"
#el uso de comillas es necesario cuando se pasa la ruta de una carpeta con espacios o usar el caracter \ para escapar los espacios.
```
Ejemplo de estructura en un dispositivo extraible:

```
run/media/alexis/6BBD-E38E/Google Fotos/
├── 2020/
│   ├── enero/
│   │   ├── foto1.jpg
│   │   ├── video.mp4        (ignorado) #en caso de usar -a mp4, no será ignorado.
│   │   └── captura.png
│   └── febrero/
│       └── imagen.webp
├── 2021/
│   ├── screenshots/
│   │   └── screenshot_001.png
│   └── wallpapers/
│       ├── fondo.jpg
│       └── logo.gif
└── 2022/
    └── mis_fotos/
        └── foto_recien_descargada.jpeg
```

Al pasar la carpeta `Google Fotos`, el programa busca **recursivamente** en todas las subcarpetas (`2020`, `2021`, `2022`, etc.) y copia solo los archivos de imagen válidos (`.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`), ignorando cualquier otro formato como `.mp4`, `.pdf`, `.txt`, etc.

Copiar de múltiples carpetas:

```bash
./flatter /home/user/fotos copy /home/user/descargas /home/user/imágenes
```

Ignorar capturas de pantalla:

```bash
./flatter -n /home/user/fotos copy /home/user/descargas
```

Usar 16 workers para mayor velocidad:

```bash
./flatter -w 16 /home/user/fotos copy /home/user/descargas
```

Omitir archivos que ya existen en el destino:

```bash
./flatter /home/user/fotos skip /home/user/descargas
```

Omitir ciertos formatos de imagen:

```bash
./flatter -i webp,jpg /home/user/fotos skip /home/user/descargas
```

Admitir formatos adicionales (ej. videos o documentos):

```bash
./flatter -a mp4,pdf /home/user/fotos copy /home/user/descargas
```

## Modos

- **`copy`**: Si el archivo ya existe, lo renombra agregando un número al final (ej: `imagen_1.jpg`)
- **`skip`**: Omite los archivos que ya existen en el destino

## Imágenes soportadas

El programa busca archivos con extensión: `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`

## Capturas de pantalla

Por defecto, el programa incluye capturas de pantalla. Usa el flag `-n` para ignorarlas. El programa detecta capturas de pantalla buscando patrones como:
- screenshots, screenshot, captura_de_pantalla, captura, screen shot, screen capture