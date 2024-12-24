# Día 4: Vamos a incorporar más servicios para un flujo más complejo, añadiendo SNS (Simple Notification Service) junto con SQS y Lambda

## Escenario

Estás desarrollando un sistema de notificaciones que necesita distribuir mensajes a diferentes subsistemas. Para lograrlo:

1. Se publican mensajes en un tema de SNS.
2. Un subsistema recibe esos mensajes a través de una cola SQS para procesamiento adicional.
3. Otro subsistema recibe esos mensajes directamente como un correo electrónico.

## Objetivo

1. Crear un flujo donde:

   - Se publique un mensaje en un tema de SNS llamado `NotificationTopic`.
   - La cola SQS llamada `ProcessingQueue` reciba automáticamente los mensajes del tema SNS.
   - Un correo electrónico de notificación se envíe a una dirección de suscripción configurada.

2. Probar el flujo enviando mensajes al tema SNS.

## Pasos para Resolver

1. Configurar el Tema SNS:

   - Ve a la consola de SNS.
   - Crea un tema llamado `NotificationTopic`.
   - Copia el ARN del tema para usarlo más adelante.

2. Crear y Configurar la Cola SQS:

   - Crea una cola llamada `ProcessingQueue`.
   - En la configuración de permisos de la cola, permite que el tema SNS `NotificationTopic` envíe mensajes a esta cola.
   - Asocia esta cola al tema `NotificationTopic` como suscriptor.

3. Configurar una Suscripción por Correo Electrónico:

   - En el tema `NotificationTopic`, agrega una suscripción del tipo email.
   - Ingresa una dirección de correo electrónico donde quieras recibir notificaciones.
   - Verifica la suscripción mediante el correo que te enviará SNS.

4. Crear la Función Lambda:

   - La función Lambda leerá mensajes de la cola SQS y realizará un procesamiento básico (por ejemplo, imprimir el mensaje en los logs de CloudWatch).
   - Configura un trigger en Lambda para que procese mensajes de la cola `ProcessingQueue`.

5. Probar el Flujo:
   - Publica un mensaje en el tema SNS usando la consola o AWS CLI:

   ```bash
   aws sns publish --topic-arn <ARN_DE_NOTIFICATIONTOPIC> --message "Test notification message"
   ```

   - Verifica
     - Que el mensaje llegue a la cola SQS (ProcessingQueue) y que la función Lambda lo procese.
     - Que el mensaje se reciba en el correo electrónico configurado.

## Preguntas de Práctica Asociadas

- ¿Cuál es la diferencia entre SNS y SQS en términos de patrones de mensajería?
- ¿Cómo podrías escalar este flujo para manejar grandes volúmenes de mensajes?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
