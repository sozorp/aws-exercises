# Día 6: Trabajaremos con API Gateway, Lambda, y Cognito para implementar una solución de autenticación

## Escenario

Estás construyendo una API REST segura que permite a los usuarios autenticados acceder a sus datos personales. Utilizarás **Amazon Cognito** para gestionar la autenticación y **API Gateway** para exponer tu API con un endpoint protegido.

## Objetivo

1. Configurar un **User Pool** en Amazon Cognito para autenticar usuarios.
2. Crear un endpoint REST con **API Gateway** que esté protegido por el User Pool.
3. Crear una función Lambda que se active desde el endpoint protegido y devuelva información personalizada del usuario autenticado.
4. Probar el flujo completo de autenticación y autorización.

## Pasos para Resolver

1. Configurar el User Pool en Cognito:

   - En la consola de Cognito, crea un nuevo **User Pool** llamado `MyUserPool`.
   - Habilita el método de autenticación por correo electrónico.
   - Configura un cliente de la aplicación sin un secreto de cliente.
   - Guarda el **ID del User Pool** y el **ID del Cliente** para usarlos más adelante.

2. Crear un Usuario de Prueba:

   - En el User Pool, crea un usuario de prueba con un correo electrónico válido.
   - Envía el correo de confirmación y verifica al usuario.

3. Crear un Endpoint REST con API Gateway:

   - En la consola de **API Gateway**, crea una nueva API REST.
   - Agrega un recurso con un método `GET`.
   - Configura este método para que esté protegido por el User Pool:
     - En **Método de Solicitud**, selecciona **Autorización Cognito**.
     - Vincula el User Pool que creaste anteriormente.

4. Crear la Función Lambda:
   - Esta función recuperará los datos del usuario autenticado y los devolverá.
   - Usa el siguiente código base en Python:

```python
def lambda_handler(event, context):
    # Obtener información del usuario autenticado
    user = event['requestContext']['authorizer']['claims']
    return {
        "statusCode": 200,
        "body": f"Hello, {user['email']}! Your user ID is {user['sub']}."
    }
```

5. Configurar la Integración de Lambda en API Gateway:

   - En el método `GET` de tu API Gateway, configura la integración para llamar a la función Lambda.
   - Asegúrate de que la función Lambda tenga permisos para ser invocada por API Gateway.

6. Probar el Flujo Completo:
   - Usa el cliente de la aplicación del User Pool para autenticarte y obtener un token de ID.
     - Esto se puede hacer con el AWS CLI o el SDK:

```bash
aws cognito-idp initiate-auth \
    --auth-flow USER_PASSWORD_AUTH \
    --client-id <ID_DEL_CLIENTE> \
    --auth-parameters USERNAME=<usuario>,PASSWORD=<contraseña>
```

Envía una solicitud `GET` al endpoint de tu API Gateway con el token como encabezado `Authorization`:

```bash
curl -H "Authorization: <TOKEN>" https://<tu-api-gateway-id>.execute-api.<region>.amazonaws.com/<stage>/<resource>
```

## Preguntas de Práctica Asociadas

- ¿Qué diferencia hay entre un token de ID y un token de acceso en Cognito?
- ¿Cómo puedes manejar roles y permisos adicionales en Cognito para controlar el acceso a recursos específicos?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
