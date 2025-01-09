# Día 19: Vamos a construir una arquitectura de procesamiento de datos en tiempo real utilizando Amazon Kinesis, AWS Lambda, y Amazon Redshift

## Escenario

Tienes una aplicación que genera eventos de transacciones financieras. Necesitas procesar estos eventos en tiempo real, almacenarlos en **Amazon Redshift** para análisis posterior, y generar alertas en caso de transacciones inusuales.

## Objetivo

1. Configurar un **Kinesis Data Stream** que reciba eventos de transacciones.
2. Crear una función Lambda que procese los eventos y los almacene en **Amazon Redshift**.
3. Configurar un **Kinesis Data Analytics** para detectar patrones inusuales.
4. Probar el flujo completo enviando eventos de prueba.

## Pasos Detallados

### Paso 1: Configurar el Kinesis Data Stream

1. Ve a la consola de **Kinesis** y crea un nuevo **Data Stream** llamado `TransactionStream`.
   - Configura el stream con **1 shard** (puedes aumentarlo si simulas una gran carga de datos).

### Paso 2: Configurar Amazon Redshift

1. Ve a la consola de **Redshift** y crea un clúster.

   - Tipo de clúster: `dc2.large`.
   - Crea una base de datos llamada `transactions`.

2. Conéctate a Redshift usando una herramienta como **DBeaver** o el cliente SQL de Redshift.
3. Crea una tabla para almacenar las transacciones:

   ```sql
   CREATE TABLE transactions (
   transaction_id VARCHAR(50),
   amount DECIMAL(10, 2),
   transaction_type VARCHAR(20),
   timestamp TIMESTAMP
   );
   ```

### Paso 3: Crear la Función Lambda

1. Crea una nueva función Lambda con el siguiente código base en Python:

   ```python
   import json
   import boto3
   import psycopg2
   import os

   # Configura las variables de entorno

   REDSHIFT_HOST = os.environ['REDSHIFT_HOST']
   REDSHIFT_USER = os.environ['REDSHIFT_USER']
   REDSHIFT_PASSWORD = os.environ['REDSHIFT_PASSWORD']
   REDSHIFT_DB = os.environ['REDSHIFT_DB']
   REDSHIFT_PORT = 5439

   def lambda_handler(event, context):
        try: # Conectar a Redshift
        conn = psycopg2.connect(
        host=REDSHIFT_HOST,
        user=REDSHIFT_USER,
        password=REDSHIFT_PASSWORD,
        dbname=REDSHIFT_DB,
        port=REDSHIFT_PORT
        )
        cursor = conn.cursor()

           # Procesar cada registro del evento de Kinesis
           for record in event['Records']:
           payload = json.loads(record['kinesis']['data'])
           transaction_id = payload['transaction_id']
           amount = payload['amount']
           transaction_type = payload['transaction_type']
           timestamp = payload['timestamp']

               # Insertar el registro en Redshift
               cursor.execute(
               "INSERT INTO transactions (transaction_id, amount, transaction_type, timestamp) VALUES    (%s, %s, %s, %s)",
               (transaction_id, amount, transaction_type, timestamp)
               )

           conn.commit()
           cursor.close()
           conn.close()
           return {
           "statusCode": 200,
           "body": json.dumps("Transactions processed successfully")
           }

       except Exception as e:
       print(f"Error: {str(e)}")
       return {
       "statusCode": 500,
       "body": json.dumps("Error processing transactions")
       }
   ```

2. Configura la variables de entorno:

   - `REDSHIFT_HOST`
   - `REDSHIFT_USER`
   - `REDSHIFT_PASSWORD`
   - `REDSHIFT_DB`

3. Configura un **trigger** en Lambda para que se active cuando lleguen nuevos eventos al stream `TransactionStream`.

### Paso 4: Configurar Kinesis Data Analytics (Opcional)

1. Ve a la consola de **Kinesis Data Analytics** y crea una nueva aplicación.
2. Configura la aplicación para leer del stream `TransactionStream`.
3. Define una consulta SQL que detecte transacciones inusuales (por ejemplo, transacciones mayores a un cierto monto):

   ```sql
   SELECT STREAM
   transaction_id,
   amount,
   transaction_type,
   timestamp
   FROM
   "SOURCE_SQL_STREAM_001"
   WHERE
   amount > 10000;
   ```

4. Configura una acción para enviar una alerta a **SNS** si se detectan transacciones inusuales.

## Prueba del Flujo Completo

1. Usa el siguiente comando para enviar eventos al stream Kinesis:

   ```bash
   aws kinesis put-record \
   --stream-name TransactionStream \
   --partition-key "1" \
   --data '{"transaction_id": "txn123", "amount": 12000, "transaction_type": "purchase", "timestamp": "2024-12-21T10:00:00Z"}'
   ```

2. Verifica en Redshift que el registro se haya insertado correctamente.

3. Si configuraste Kinesis Data Analytics, verifica que se genere una alerta en caso de transacciones inusuales.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías el escalado del Kinesis Data Stream si el volumen de eventos aumenta significativamente?
- ¿Qué ventajas tiene Kinesis Data Analytics frente a soluciones personalizadas para detección de patrones?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
