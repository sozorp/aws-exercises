# Día 8: Vamos a incorporar un flujo de procesamiento event-driven más complejo, utilizando Kinesis, Lambda, y DynamoDB

## Escenario

Tu aplicación debe procesar un flujo constante de eventos de datos. Los datos llegan en tiempo real a través de un stream de Kinesis, donde cada evento se registra y se guarda en una tabla DynamoDB para análisis posterior.

## Objetivo

1. Configurar un **Kinesis Data Stream** para recibir eventos en tiempo real.
2. Crear una **función Lambda** que lea eventos desde el stream y guarde los datos en DynamoDB.
3. Probar el flujo completo enviando eventos al stream y verificando que se guarden en la tabla DynamoDB.

## Pasos para Resolver

1. Configurar un Stream de Kinesis:

   - Ve a la consola de Kinesis y crea un **Data Stream** llamado `EventStream`.
   - Configura el stream con **1 shard** (suficiente para este ejercicio).

2. Configurar DynamoDB:

   - Crea una tabla DynamoDB llamada `EventLogs`.
   - Configura los atributos:
     - **Partition Key**: `eventId` (String).
   - Opcional: Añade un índice secundario para búsquedas avanzadas.

3. Crear la Función Lambda:

   - La función debe:
     - Leer eventos del stream.
     - Procesar los datos y guardar cada evento en la tabla DynamoDB.
   - Usa este código base en Python:

   ```python
   import json
   import boto3
   import uuid
   from datetime import datetime


   dynamodb = boto3.resource('dynamodb')
   table = dynamodb.Table('EventLogs')

   def lambda_handler(event, context):
       for record in event['Records']: # Obtener datos del evento
       payload = json.loads(record['kinesis']['data'])
       event_id = str(uuid.uuid4())
       timestamp = datetime.utcnow().isoformat()

        # Guardar en DynamoDB
        table.put_item(
            Item={
                'eventId': event_id,
                'timestamp': timestamp,
                'data': payload
            }
        )
        print(f"Processed event: {payload}")

    return {
        "statusCode": 200,
        "body": json.dumps("All events processed successfully")
    }
   ```

4. Configurar el Trigger en Lambda:

   - Ve a la consola de **Lambda** y crea una nueva función con el código anterior.
   - Configura un trigger desde el stream de Kinesis (`EventStream`) a la función Lambda.

5. Probar el Flujo:

   - Usa el siguiente comando de AWS CLI para enviar eventos al stream:

   ```bash
   aws kinesis put-record \
   --stream-name EventStream \
   --partition-key "1" \
   --data '{"sensor": "temperature", "value": 22.5}'
   ```

- Envía varios eventos con diferentes datos.
- Verifica en DynamoDB que cada evento esté guardado con su eventId único y los datos procesados.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías una gran cantidad de datos simultáneos en el stream?
- ¿Qué harías si un evento no puede ser procesado por Lambda?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
