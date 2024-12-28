# Día 9: Vamos a combinar varios servicios para trabajar con AWS Elastic Beanstalk y automatizar el despliegue de una aplicación web sencilla

## Escenario

Tienes una aplicación web escrita en **Python Flask** que muestra datos de DynamoDB. Quieres desplegarla usando **AWS Elastic Beanstalk** para que se gestione automáticamente la infraestructura (EC2, Load Balancers, etc.).

## Objetivo

1. Desplegar una aplicación Flask en Elastic Beanstalk.
2. Configurar la aplicación para conectarse a DynamoDB y mostrar datos en un endpoint.
3. Probar la aplicación en un entorno en vivo.

## Pasos para Resolver

1. Preparar la Aplicación Flask:

   - Crea un archivo `application.py` con el siguiente contenido:

   ```python
   from flask import Flask, jsonify
   import boto3

   app = Flask(__name__)
   dynamodb = boto3.resource('dynamodb')
   table = dynamodb.Table('EventLogs')  # Asegúrate de usar la tabla DynamoDB del ejercicio anterior.

   @app.route('/')
   def home():
       return "Welcome to the Flask App on Elastic Beanstalk!"

   @app.route('/events')
   def get_events():
       # Escanear la tabla DynamoDB
       response = table.scan()
       events = response.get('Items', [])
       return jsonify(events)

   if __name__ == '__main__':
       app.run(debug=True)
   ```

   - Crea un archivo `requirements.txt` con las dependencias de la aplicación:

   ```text
   flask
   boto3
   ```

2. Configurar Elastic Beanstalk:

   - Instala la CLI de Elastic Beanstalk si aún no lo has hecho:

   ```bash
   pip install awsebcli
   ```

   - Inicializa un nuevo entorno de Elastic Beanstalk:

   ```bash
    eb init
   ```

   - Selecciona tu región.
   - Escoge el lenguaje Python.

   - Crea un entorno para la aplicación:

   ```bash
   eb create flask-app-env
   ```

3. Desplegar la Aplicación:

   - Sube el código al entorno de Elastic Beanstalk:

   ```bash
   eb deploy
   ```

4. Configurar Permisos:

   - Asegúrate de que la instancia EC2 que se crea tenga un rol con permisos para acceder a DynamoDB.

5. Probar la Aplicación:

   - Una vez que Elastic Beanstalk termine el despliegue, obtén la URL del entorno:

   ```bash
   eb open
   ```

   - Verifica:

     - Que el endpoint `/` muestra el mensaje de bienvenida.
     - Que el endpoint `/events` devuelve datos de la tabla DynamoDB.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías la escalabilidad de tu aplicación en Elastic Beanstalk?
- ¿Qué sucede si necesitas agregar variables de entorno específicas para tu aplicación?

> [!TIP]
> Puedes usar el archivo `lambda.sh` para automatizar la construcción y empaquetado de tu función Lambda.

Happy hacking! 🚀
