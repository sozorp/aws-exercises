# Día 12: Vamos a trabajar con Amazon Aurora Serverless y API Gateway para construir una solución escalable basada en una base de datos relacional administrada

## Escenario

Estás construyendo un sistema donde una aplicación consulta y actualiza datos de una base de datos relacional. Usarás **Amazon Aurora Serverless** para la base de datos, **Lambda** como capa de aplicación, y **API Gateway** para exponer los endpoints.

## Objetivo

1. Crear una base de datos **Aurora Serverless** para almacenar datos.
2. Crear una función Lambda para interactuar con la base de datos (consultar y actualizar datos).
3. Exponer dos endpoints mediante API Gateway:

   - `GET /users`: Obtiene todos los usuarios.
   - `POST /users`: Crea un nuevo usuario.

4. Probar el flujo completo desde la API.

## Pasos para Resolver

1. Configurar Aurora Serverless:

   - Ve a la consola de **RDS** y crea un clúster Aurora Serverless:

     - Motor: **MySQL/PostgreSQL-compatible Aurora**.
     - Configura una base de datos llamada `UserDB`.

   - Una vez creado el clúster, toma nota del **endpoint** de conexión.
   - En la base de datos, ejecuta este script SQL para crear una tabla:

   - `MySQL`:

   ```sql
   CREATE TABLE Users (
     id INT AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     email VARCHAR(100) NOT NULL
   );
   ```

   - `PostgreSQL`:

   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       email VARCHAR(100) NOT NULL UNIQUE,
       name VARCHAR(100) NOT NULL
   );
   ```

2. Crear la Función Lambda:

   - La función Lambda debe conectarse a la base de datos y realizar consultas. Usa este código base en Python (esta basado para mysql):

   ```python
   import pymysql
   import json
   import os

   # Configura las credenciales y el endpoint de la base de datos
   DB_HOST = os.environ['DB_HOST']
   DB_USER = os.environ['DB_USER']
   DB_PASSWORD = os.environ['DB_PASSWORD']
   DB_NAME = os.environ['DB_NAME']

   def lambda_handler(event, context):
       connection = pymysql.connect(
           host=DB_HOST,
           user=DB_USER,
           password=DB_PASSWORD,
           database=DB_NAME
       )
       try:
           if event['httpMethod'] == 'GET':
               with connection.cursor(pymysql.cursors.DictCursor) as cursor:
                   cursor.execute("SELECT * FROM Users;")
                   results = cursor.fetchall()
                   return {
                       "statusCode": 200,
                       "body": json.dumps(results)
                   }

           elif event['httpMethod'] == 'POST':
               body = json.loads(event['body'])
               name = body['name']
               email = body['email']
               with connection.cursor() as cursor:
                   cursor.execute("INSERT INTO Users (name, email) VALUES (%s, %s);", (name, email))
                   connection.commit()
                   return {
                       "statusCode": 201,
                       "body": json.dumps({"message": "User created successfully"})
                   }

           else:
               return {
                   "statusCode": 400,
                   "body": json.dumps({"message": "Unsupported HTTP method"})
               }
       finally:
           connection.close()
   ```

3. Configurar la Lambda:

   - Crea una nueva función Lambda y sube el código anterior.
   - Usa **Layer** para incluir la biblioteca `pymysql` o empaqueta el código con la dependencia.
   - Configura las siguientes variables de entorno:
     - `DB_HOST`: Endpoint de Aurora Serverless.
     - `DB_USER`: Nombre de usuario de la base de datos.
     - `DB_PASSWORD`: Contraseña de la base de datos.
     - `DB_NAME`: `UserDB`.

4. Configurar API Gateway:

   - Crea una nueva API REST en API Gateway.
   - Define los endpoints:
     - `GET /users`: Vincúlalo a la función Lambda.
     - `POST /users`: Vincúlalo también a la misma función Lambda.
   - Asegúrate de habilitar **CORS**.

5. Probar el Flujo:

   - Usa **curl** o Postman para probar los endpoints:
   - `GET /users`:

   ```bash
   curl https://<tu-api-gateway-id>.execute-api.<region>.amazonaws.com/<stage>/users
   ```

   - `POST /users`:

   ```bash
   curl -X POST https://<tu-api-gateway-id>.execute-api.<region>.amazonaws.com/<stage>/users \
    -H "Content-Type: application/json" \
    -d '{"name": "John Doe", "email": "john.doe@example.com"}'
   ```

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías credenciales sensibles (como el usuario y la contraseña de la base de datos) en un entorno real?
- ¿Qué ventajas tiene Aurora Serverless frente a RDS tradicional?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
