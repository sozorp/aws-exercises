# Día 15: Sistema de event-driven processing usando EventBridge, Lambda, y SNS

## Escenario

Vas a crear un sistema de **event-driven** processing usando **EventBridge**, **Lambda**, y **SNS**. El objetivo es 1capturar eventos generados por un servicio ficticio de pedidos y enviar notificaciones a los usuarios cuando un pedido sea creado o cancelado.

## Objetivo

1. Configurar un **bus de eventos** en EventBridge.
2. Crear reglas en EventBridge para detectar eventos específicos (`OrderCreated` y `OrderCancelled`).
3. Configurar dos funciones Lambda:
   - Una que maneje el evento `OrderCreated`.
   - Otra que maneje el evento `OrderCancelled`.
4. Enviar una notificación por **SNS** cuando ocurra alguno de los eventos.
5. Probar el flujo enviando eventos manualmente a EventBridge.

## Pasos Detallados

1. Crear el Bus de Eventos:

   - Ve a la consola de **EventBridge**.
   - Crea un nuevo **bus de eventos** llamado `OrderEventBus`.

2. Crear el Tema SNS:

   - Ve a la consola de **SNS** y crea un tema llamado `OrderNotifications`.
   - Suscribe una dirección de correo electrónico al tema y verifica la suscripción.

3. Crear las Funciones Lambda:

   - Crea dos funciones Lambda con el siguiente código:

   **Función 1: order_created_handler.py**

   ```python
   import json
   import boto3

   sns = boto3.client('sns')
   SNS_TOPIC_ARN = 'arn:aws:sns:<region>:<account-id>:OrderNotifications'

   def lambda_handler(event, context):
   order_id = event['detail']['orderId']
   message = f"Order {order_id} has been created."
   sns.publish(
       TopicArn=SNS_TOPIC_ARN,
       Subject="Order Created",
       Message=message
   )
   return {
       "statusCode": 200,
       "body": json.dumps("Notification sent for OrderCreated event")
   }
   ```

   **Función 2: order_cancelled_handler.py**

   ```python
   import json
   import boto3

   sns = boto3.client('sns')
   SNS_TOPIC_ARN = 'arn:aws:sns:<region>:<account-id>:OrderNotifications'

   def lambda_handler(event, context):
   order_id = event['detail']['orderId']
   message = f"Order {order_id} has been cancelled."
   sns.publish(
       TopicArn=SNS_TOPIC_ARN,
       Subject="Order Cancelled",
       Message=message
   )
   return {
       "statusCode": 200,
       "body": json.dumps("Notification sent for OrderCancelled event")
   }
   ```

4. Configurar las Reglas en EventBridge:

   - Ve a la consola de **EventBridge**.
   - Crea una nueva regla llamada `OrderCreatedRule`:

     - Fuente del evento: `Event Bus -> OrderEventBus`.
     - Patrón del evento:

     ```json
     {
       "detail-type": ["OrderCreated"]
     }
     ```

     - Como destino, selecciona la función Lambda `order_created_handler`.

   - Crea otra regla llamada `OrderCancelledRule`:

     - Fuente del evento: `Event Bus -> OrderEventBus`.
     - Patrón del evento:

     ```json
     {
       "detail-type": ["OrderCancelled"]
     }
     ```

     - Como destino, selecciona la función Lambda `order_cancelled_handler`.

5. Probar el Flujo:

   - Ve a la consola de **EventBridge** y selecciona el bus de eventos `OrderEventBus`.
   - Envía un evento de prueba para `OrderCreated`:

   ```json
   {
     "source": "ecommerce.orders",
     "detail-type": "OrderCreated",
     "detail": {
       "orderId": "12345"
     }
   }
   ```

   - Envía otro evento de prueba para `OrderCancelled`:

   ```json
   {
     "source": "ecommerce.orders",
     "detail-type": "OrderCancelled",
     "detail": {
       "orderId": "67890"
     }
   }
   ```

   - Verifica que recibes los correos electrónicos de notificación correspondientes en tu bandeja de entrada.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías múltiples tipos de eventos en una misma función Lambda?
- ¿Qué beneficios ofrece EventBridge frente a soluciones personalizadas para manejo de eventos?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
