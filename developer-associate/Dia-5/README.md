# Día 5: Vamos a escalar un poco más el nivel de dificultad en este Día 5 al introducir AWS Step Functions para orquestar un flujo de trabajo

## Escenario

Estás desarrollando un sistema de procesamiento de pedidos que necesita coordinar varias tareas. Por ejemplo:

1. Validar el pedido.
2. Procesar el pago.
3. Enviar una notificación al cliente.

Necesitas garantizar que cada paso se ejecute en orden y manejar posibles fallos.

## Objetivo

1. Crear un flujo de trabajo orquestado con AWS Step Functions que realice las siguientes tareas:

   - Paso 1: Validar el pedido (simulado por una función Lambda que verifica los datos del pedido).
   - Paso 2: Procesar el pago (simulado por una función Lambda que retorna éxito o error aleatoriamente).
   - Paso 3: Notificar al cliente (simulado por una función Lambda que imprime el resultado en los logs).

2. Manejar fallos en el paso de procesamiento de pagos para que el flujo registre el error y detenga el proceso.

## Pasos para Resolver

1. Crear las Funciones Lambda:

   - Función 1: Validar Pedido

   ```python
    def lambda_handler(event, context):
      print("Validating order:", event)
      if "orderId" not in event:
        raise ValueError("Invalid order: Missing orderId")
      return {"status": "Order Validated", "orderId": event["orderId"]}
   ```

   - Función 2: Procesar Pago

   ```python
    import random

    def lambda_handler(event, context):
      print("Processing payment for order:", event)
      if random.choice([True, False]):
        raise ValueError("Payment processing failed")
      return {"status": "Payment Successful", "orderId": event["orderId"]}
   ```

   - Función 3: Notificar al Cliente

   ```python
    def lambda_handler(event, context):
      print("Notifying customer for order:", event)
      return {"status": "Notification Sent", "orderId": event["orderId"]}
   ```

2. Crear el Flujo de Step Functions:
   - Diseña un flujo que:
   - Llame primero a la función de validación.
   - Luego llame a la función de procesamiento de pago.
   - Si el pago falla, capture el error y registre el fallo.
   - Si el pago es exitoso, continúe notificando al cliente.

Ejemplo de definición en formato JSON (puedes usar el editor visual de Step Functions en la consola de AWS):

```json
{
  "Comment": "Order processing workflow",
  "StartAt": "ValidateOrder",
  "States": {
    "ValidateOrder": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:<region>:<account-id>:function:<ValidateOrderLambda>",
      "Next": "ProcessPayment"
    },
    "ProcessPayment": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:<region>:<account-id>:function:<ProcessPaymentLambda>",
      "Catch": [
        {
          "ErrorEquals": ["States.ALL"],
          "ResultPath": "$.error",
          "Next": "PaymentFailed"
        }
      ],
      "Next": "NotifyCustomer"
    },
    "PaymentFailed": {
      "Type": "Fail",
      "Error": "PaymentFailure",
      "Cause": "The payment could not be processed."
    },
    "NotifyCustomer": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:<region>:<account-id>:function:<NotifyCustomerLambda>",
      "End": true
    }
  }
}
```

3. Probar el Flujo:
   - Inicia una ejecución del flujo de trabajo con un evento de entrada

```json
{
  "orderId": "12345",
  "customer": "example@example.com"
}
```

- Observa el comportamiento del flujo en la consola de Step Functions:
  - Si el pago falla, el flujo debería detenerse con el error registrado.
  - Si el pago tiene éxito, debería continuar con la notificación.

## Preguntas de Práctica Asociadas

- ¿Qué ventajas tienen Step Functions sobre usar solo Lambda y SQS para flujos complejos?
- ¿Cómo se implementan tareas paralelas en un flujo de Step Functions?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
