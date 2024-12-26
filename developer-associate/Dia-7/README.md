# Día 7: Vamos a construir un flujo más completo que combine varios servicios de AWS

## Escenario

Estás desarrollando una aplicación para generar y distribuir reportes personalizados. Cada vez que un usuario solicita un reporte:

1. Una solicitud se envía a un endpoint de API Gateway.
2. Se activa una función Lambda que procesa la solicitud, genera un archivo, y lo guarda en S3.
3. Una notificación se envía al usuario (a través de SNS) con un enlace para descargar el reporte.

## Objetivo

1. Crear un flujo que integre **API Gateway**, **Lambda**, **S3**, y **SNS**.
2. Implementar la función Lambda para:

   - Recibir una solicitud HTTP con parámetros como `userId` y `reportType`.
   - Generar un archivo (puede ser un archivo de texto o JSON).
   - Guardar el archivo en un bucket S3 con un nombre único.
   - Enviar una notificación al usuario (a través de SNS) con el enlace al archivo generado.

3. Probar el flujo completo.

## Pasos para Resolver

1. Configurar el Bucket S3:

   - Crea un bucket llamado `report-generator-{tu-nombre}`.
   - Habilita acceso público para objetos específicos (opcional si quieres generar un enlace público).

2. Configurar SNS:

   - Crea un tema SNS llamado `ReportNotifications`.
   - Suscribe una dirección de correo electrónico al tema.
   - Verifica la suscripción.

3. Crear la Función Lambda:

   - La función Lambda debe realizar las siguientes tareas:
     - Validar los parámetros de entrada (`userId` y `reportType`).
     - Generar un archivo (simulado).
     - Subir el archivo al bucket S3.
     - Publicar un mensaje en el tema SNS con el enlace al archivo.
     - Código base en Python:

Código base en Python:

```python
import json
import boto3
from datetime import datetime

s3 = boto3.client('s3')
sns = boto3.client('sns')

BUCKET_NAME = "report-generator-yourname"
SNS_TOPIC_ARN = "arn:aws:sns:<region>:<account-id>:ReportNotifications"

def lambda_handler(event, context):
    # Obtener parámetros de entrada
    user_id = event['queryStringParameters'].get('userId')
    report_type = event['queryStringParameters'].get('reportType')

    if not user_id or not report_type:
        return {
            "statusCode": 400,
            "body": json.dumps({"message": "Missing userId or reportType"})
        }

    # Generar el archivo
    file_content = f"Report for {user_id}\nType: {report_type}\nGenerated at: {datetime.utcnow()}"
    file_name = f"{user_id}-{report_type}-{int(datetime.utcnow().timestamp())}.txt"

    # Subir el archivo a S3
    s3.put_object(Bucket=BUCKET_NAME, Key=file_name, Body=file_content)
    file_url = f"https://{BUCKET_NAME}.s3.amazonaws.com/{file_name}"

    # Enviar notificación a SNS
    sns.publish(
        TopicArn=SNS_TOPIC_ARN,
        Subject="Your report is ready",
        Message=f"Hello {user_id}, your report is ready: {file_url}"
    )

    return {
        "statusCode": 200,
        "body": json.dumps({"message": "Report generated successfully", "reportUrl": file_url})
    }
```

4. Configurar API Gateway:

   - Crea un endpoint `GET` que invoque la función Lambda.
   - Asegúrate de pasar los parámetros `userId` y `reportType` desde la URL.

5. Probar el Flujo:
   - Envía una solicitud al endpoint con parámetros, por ejemplo:

```bash
curl "https://<tu-api-gateway-id>.execute-api.<region>.amazonaws.com/<stage>?userId=123&reportType=summary"
```

- Verifica:
  - Que el archivo se genere en S3.
  - Que se reciba un correo electrónico con el enlace al archivo.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías el acceso seguro a los archivos en S3?
- ¿Qué alternativas podrías usar en lugar de SNS para notificar a los usuarios?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
