# Flatter

Herramienta para copiar imágenes de múltiples carpetas a una carpeta destino.

## Instalación

```bash
go build -o flatter .
```

## Uso

```bash
go run main.go <destino> <modo> <carpeta_raiz> o <capeta1> <carpeta2...
```

### Argumentos

| Argumento | Descripción |
|-----------|-------------|
| `destino` | Ruta de la carpeta destino |
| `modo` | Modo de copia: `copy` (renombra si existe) o `skip` (omite si existe) |
| `carpeta1`... | Carpetas fuente a escanear |

### Flags

| Flag | Descripción | Valor por defecto |
|------|-------------|-------------------|
| `-workers` | Número de workers concurrentes | 8 |
| `-nsc` | Ignorar capturas de pantalla | false |

## Ejemplos

Copiar todas las imágenes de una carpeta (incluyendo subcarpetas):

```bash
go run main.go /home/user/fotos copy /run/media/alexis/6BBD-E38E/Google\ Fotos
```

Ejemplo de estructura en un dispositivo extraible:

```
run/media/alexis/6BBD-E38E/Google Fotos/
├── 2020/
│   ├── enero/
│   │   ├── foto1.jpg
│   │   ├── video.mp4        (ignorado)
│   │   └── captura.png
│   └── febrero/
│       └── imagen.webp
├── 2021/
│   ├── screenshots/
│   │   └── screenshot_001.png
│   └──wallpapers/
│       ├── fondo.jpg
│       └── logo.gif
└── 2022/
    └── mis_fotos/
        └── foto_recien_descargada.jpeg
```

Al pasar la carpeta `Google Fotos`, el programa busca **recursivamente** en todas las subcarpetas (`2020`, `2021`, `2022`, etc.) y copia solo los archivos de imagen válidos (`.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`), ignorando cualquier otro formato como `.mp4`, `.pdf`, `.txt`, etc.

Copiar de múltiples carpetas:

```bash
go run main.go /home/user/fotos copy /home/user/descargas /home/user/imágenes
```

Ignorar capturas de pantalla:

```bash
go run main.go -nsc /home/user/fotos copy /home/user/descargas
```

Usar 16 workers para mayor velocidad:

```bash
go run main.go -workers 16 /home/user/fotos copy /home/user/descargas
```

Omitir archivos que ya existen en el destino:

```bash
go run main.go /home/user/fotos skip /home/user/descargas
```

## Modos

- **`copy`**: Si el archivo ya existe, lo renombra agregando un número al final (ej: `imagen_1.jpg`)
- **`skip`**: Omite los archivos que ya existen en el destino

## Imágenes soportadas

El programa busca archivos con extensión: `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`

## Capturas de pantalla

Por defecto, el programa incluye capturas de pantalla. Usa el flag `-nsc` para ignorarlas. El programa detecta capturas de pantalla buscando patrones como:
- screenshots, screenshot, captura_de_pantalla, captura, screen shot, screen capture