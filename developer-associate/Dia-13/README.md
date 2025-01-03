# Día 13: Vamos a construir un sistema de procesamiento de imágenes con escalabilidad automática. Este ejercicio incluirá S3, Lambda, DynamoDB, y CloudWatch Events

## Escenario

Estás construyendo un sistema que procesa imágenes cargadas en un bucket S3. Cada vez que un usuario sube una imagen, el sistema genera una miniatura y guarda información sobre la imagen en DynamoDB. También configuraremos una alerta que se dispare si hay errores en el procesamiento.

## Objetivo

1. Configurar un bucket S3 para cargar imágenes.
2. Crear una función Lambda que:
   - Reciba el evento de carga.
   - Genere una miniatura de la imagen.
   - Guarde la miniatura en S3.
   - Registre la información de la imagen y su miniatura en DynamoDB.
3. Configurar una alerta con **CloudWatch Events** para notificar errores.
4. Probar el flujo cargando imágenes al bucket.

## Pasos para Resolver

1. Configurar el Bucket S3:

   - Crea un bucket llamado `image-processing-{tu-nombre}`.
   - Configura una carpeta `uploads/` para las imágenes originales y `thumbnails/` para las miniaturas.

2. Configurar DynamoDB:

   - Crea una tabla llamada `ImageMetadata` con:
     - **Partition Key**: `imageId` (String).
   - Añade un índice secundario global opcional para búsquedas basadas en otros atributos como `uploadDate`.

3. Crear la Función Lambda:

   - Usa el siguiente código base en Python (requiere la biblioteca `Pillow` para manipular imágenes):

   ```python
   import boto3
   import os
   import uuid
   from PIL import Image
   from datetime import datetime

   s3 = boto3.client('s3')
   dynamodb = boto3.resource('dynamodb')
   table = dynamodb.Table('ImageMetadata')

   def lambda_handler(event, context):
    try:
    for record in event['Records']:
    bucket = record['s3']['bucket']['name']
    key = record['s3']['object']['key']

            # Descargar la imagen desde S3
            download_path = f"/tmp/{key.split('/')[-1]}"
            s3.download_file(bucket, key, download_path)

            # Crear una miniatura
            thumbnail_path = f"/tmp/thumbnail-{uuid.uuid4().hex}.jpg"
            with Image.open(download_path) as img:
                img.thumbnail((128, 128))
                img.save(thumbnail_path)

            # Subir la miniatura a S3
            thumbnail_key = f"thumbnails/{os.path.basename(thumbnail_path)}"
            s3.upload_file(thumbnail_path, bucket, thumbnail_key)

            # Guardar metadata en DynamoDB
            table.put_item(
                Item={
                    "imageId": str(uuid.uuid4()),
                    "originalImage": key,
                    "thumbnailImage": thumbnail_key,
                    "uploadDate": datetime.utcnow().isoformat()
                }
            )
            print(f"Processed {key} and created {thumbnail_key}")
    except Exception as e:
        print(f"Error processing file: {str(e)}")
        raise e
   ```

4. Configurar la Lambda:

   - Crea una nueva función Lambda y sube el código.
   - Usa **Layer** para incluir la biblioteca `Pillow` o empaqueta el código con la dependencia.
   - Configura un trigger de evento en el bucket S3 para que se active en la carpeta `uploads/`.

5. Configurar la Alerta con CloudWatch:

   - Crea una métrica personalizada para rastrear errores de Lambda:
     - En el código, asegúrate de registrar los errores.
   - Configura una alerta en **CloudWatch Alarms**:
     - Métrica: Invocaciones fallidas de Lambda.
     - Acción: Enviar un correo electrónico a través de SNS.

6. Probar el Flujo:
   - Sube una imagen a la carpeta `uploads/` del bucket S3:

   ```bash
   aws s3 cp test-image.jpg s3://image-processing-{tu-nombre}/uploads/
   ```

   -Verifica:
   - Que se crea la miniatura en `thumbnails/`.
   - Que los datos de la imagen se registran en la tabla DynamoDB.
   - Que se genera una alerta si la Lambda falla.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías grandes cargas de imágenes para evitar límites de Lambda?
- ¿Cómo optimizarías el acceso concurrente a DynamoDB desde Lambda?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
