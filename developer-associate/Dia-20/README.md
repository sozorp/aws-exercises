# Día 20: Vamos a construir una arquitectura más avanzada que incluya AWS Step Functions, Lambda, y S3 para orquestar un flujo de procesamiento de datos

## Escenario

Tu sistema debe procesar archivos cargados en S3 de manera automática. El flujo debe incluir:

1. Validar si el archivo cumple con ciertos criterios.
2. Procesar el archivo si es válido.
3. Mover el archivo a una carpeta de **procesados** o de **errores** dependiendo del resultado.

Usaremos **Step Functions** para orquestar este flujo y manejar los casos de éxito y fallos.

## Objetivo

1. Configurar un bucket S3 para almacenar archivos.
2. Crear funciones Lambda para:

   - Validar el archivo.
   - Procesar el archivo.
   - Mover el archivo a la carpeta correspondiente.

3. Crear una máquina de estado en **Step Functions** que orqueste el flujo.
4. Probar el flujo cargando un archivo de prueba.

## Pasos Detallados

### Paso 1: Configurar el Bucket S3

1. Ve a la consola de **S3** y crea un bucket llamado `file-processing-{tu-nombre}`.
2. Dentro del bucket, crea las siguientes carpetas:

   - `uploads/`
   - `processed/`
   - `errors/`

### Paso 2: Crear las Funciones Lambda

1. Función 1: validate_file
   Esta función verifica que el archivo cumpla con ciertos criterios, como el tamaño o el formato.

   ```python

   import json
   import boto3

   def lambda_handler(event, context):
   s3 = boto3.client('s3')
   bucket = event['bucket']
   key = event['key']

       # Simulación de validación: archivo debe ser menor a 1 MB
       response = s3.head_object(Bucket=bucket, Key=key)
       file_size = response['ContentLength']

       if file_size > 1 * 1024 * 1024:  # 1 MB
       raise ValueError("File size exceeds the allowed limit")

       return {
       "status": "valid",
       "bucket": bucket,
       "key": key
       }
   ```

2. Función 2: process_file
   Procesa el archivo (en este caso, simplemente registra el nombre).

   ```python
   import json

   def lambda_handler(event, context):
        bucket = event['bucket']
        key = event['key']
        print(f"Processing file: s3://{bucket}/{key}")
        return {"status": "processed", "bucket": bucket, "key": key}
   ```

3. Función 3: move_file
   Mueve el archivo a la carpeta `processed/` o `errors/ según el resultado.

   ```python
   import boto3

   def lambda_handler(event, context):
       s3 = boto3.client('s3')
       source_bucket = event['bucket']
       source_key = event['key']
       destination_folder = event['destination']

       destination_key = f"{destination_folder}/{source_key.split('/')[-1]}"
       s3.copy_object(Bucket=source_bucket, CopySource={'Bucket': source_bucket, 'Key': source_key},    Key=destination_key)
       s3.delete_object(Bucket=source_bucket, Key=source_key)

       return {"status": "moved", "destination": destination_key}
   ```

### Paso 3: Crear la Máquina de Estado en Step Functions

1. Ve a la consola de **Step Functions** y crea una nueva máquina de estado.

2. Usa esta definición de estado en formato **JSON**:

   ```json
   {
     "Comment": "File processing workflow",
     "StartAt": "Validate File",
     "States": {
       "Validate File": {
         "Type": "Task",
         "Resource": "arn:aws:lambda:<region>:<account-id>:function:validate_file",
         "Next": "Process File",
         "Catch": [
           {
             "ErrorEquals": ["States.ALL"],
             "ResultPath": "$.error",
             "Next": "Move to Errors"
           }
         ]
       },
       "Process File": {
         "Type": "Task",
         "Resource": "arn:aws:lambda:<region>:<account-id>:function:process_file",
         "Next": "Move to Processed"
       },
       "Move to Processed": {
         "Type": "Task",
         "Resource": "arn:aws:lambda:<region>:<account-id>:function:move_file",
         "Parameters": {
           "bucket.$": "$.bucket",
           "key.$": "$.key",
           "destination": "processed"
         },
         "End": true
       },
       "Move to Errors": {
         "Type": "Task",
         "Resource": "arn:aws:lambda:<region>:<account-id>:function:move_file",
         "Parameters": {
           "bucket.$": "$.bucket",
           "key.$": "$.key",
           "destination": "errors"
         },
         "End": true
       }
     }
   }
   ```

### Paso 4: Probar el Flujo

1. Sube un archivo de prueba al bucket S3 en la carpeta `uploads/`:

   ```bash
   aws s3 cp test-file.txt s3://file-processing-{tu-nombre}/uploads/
   ```

2. Ejecuta la máquina de estado en **Step Functions** con el siguiente evento de inicio:

   ```json
   {
     "bucket": "file-processing-{tu-nombre}",
     "key": "uploads/test-file.txt"
   }
   ```

3. Verifica:

   - Si el archivo es válido, debe moverse a `processed/`.
   - Si no cumple con los criterios, debe moverse a `errors/`.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías flujos que incluyan más de dos ramas condicionales?
- ¿Qué beneficios ofrece Step Functions frente a orquestar los flujos directamente desde Lambda?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
