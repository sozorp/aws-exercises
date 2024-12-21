# Día 3: Vamos a trabajar con SQS y AWS Lambda para procesar eventos de manera asincrónica

## Escenario

Tienes una aplicación que genera tareas que necesitan ser procesadas de manera independiente. Los mensajes de las tareas se colocan en una cola de SQS, y una función Lambda se encargará de procesarlos cuando lleguen.

## Objetivo

1. Configurar una cola SQS llamada `TaskQueue`.
2. Crear una función Lambda que lea mensajes de la cola y procese cada tarea simulando un "trabajo pesado" (por ejemplo, imprimir información en los logs de CloudWatch).
3. Después de procesar el mensaje, eliminarlo de la cola.
4. Probar el flujo enviando mensajes a SQS usando AWS CLI o SDK.

## Pasos a resolver

1. Configura una cola SQS:

   - Ve a la consola de SQS y crea una cola llamada TaskQueue.
     -Asegúrate de usar el tipo de cola estándar.

2. Crea una función Lambda para procesar los mensajes

3. Configura un trigger en Lambda:

   - Ve a la configuración de la función Lambda.
   - Agrega un trigger para que la Lambda se active automáticamente cuando haya mensajes en la cola `TaskQueue`.

4. Prueba el flujo:

   - Envia mensajes a la cola SQS con el siguiente comando de AWS CLI o usa la consola:

   ```bash
   aws sqs send-message --queue-url https://sqs.<region>.amazonaws.com/<account-id>/TaskQueue --message-body "Task 1"
   ```

   - Verifica los logs de la función Lambda en **CloudWatch Logs** para confirmar que se procesaron los mensajes correctamente.

5. **(Opcional)**: Configura un tiempo de espera (visibility timeout) en la cola para que, si Lambda falla al procesar un mensaje, este vuelva a estar disponible después de un tiempo.

## Preguntas de Práctica Asociadas

- ¿Qué sucede si una función Lambda falla al procesar un mensaje de SQS?
- ¿Cómo manejarías los mensajes que no se pueden procesar después de múltiples intentos? (Pista: Dead-Letter Queue).

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
