# Día 1: Crea una aplicación básica que use AWS Lambda y DynamoDB

## Escenario

Eres un desarrollador que necesita registrar visitas en un sitio web. Diseña una solución que haga lo siguiente:

1. Usa una función Lambda que se active mediante una solicitud HTTP a través de API Gateway.

2. La función debe guardar un registro de la visita en una tabla DynamoDB. Cada registro debe incluir:
    - ID único para la visita (UUID).
    - Timestamp de la visita.
    - User Agent del navegador del visitante.

## Instrucciones

1. Crea una tabla DynamoDB con los siguientes campos:
    - Partition Key: visitId (String).
    - Atributos adicionales: timestamp, userAgent.

2. Configura un API Gateway que acepte solicitudes POST y active tu función Lambda.

3. Prueba tu API usando Postman o curl y verifica que los datos se guarden correctamente en DynamoDB.

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
