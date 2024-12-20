# Día 2: Hoy trabajaremos con Amazon S3 y AWS Lambda, integrándolos para crear una funcionalidad común en la certificación

## Escenario

Eres un desarrollador encargado de procesar imágenes subidas a un bucket S3. Cada vez que un usuario sube una imagen, necesitas generar una entrada de log que se guarde en otra carpeta dentro del mismo bucket.

## Objetivo

1. Configurar un bucket S3 con dos carpetas virtuales:
    - uploads/ para las imágenes subidas.
    - logs/ para guardar archivos JSON con información de cada subida.

2. Configurar un evento de S3 que active una función Lambda cuando se suba un archivo a uploads/.

3. La función Lambda debe generar un archivo JSON con la siguiente información:
    - Nombre del archivo subido.
    - Tamaño del archivo (en bytes).
    - Fecha y hora de la subida.

4. Guardar este archivo JSON en la carpeta logs/ dentro del mismo bucket.

## Pasos para Resolver

1. Configura un bucket S3:
    - Crea un bucket llamado image-processor-{tu-nombre}.
    - Crea dos carpetas: uploads/ y logs/.

2. Configura un evento S3:
    - Agrega un evento de bucket para que cualquier archivo subido a uploads/ active una función Lambda.

3. Prueba el flujo:
    - Sube un archivo a la carpeta uploads/ usando la consola de AWS o aws s3 cp.
    - Verifica que el archivo JSON con los detalles del log se guarde en logs/.

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
