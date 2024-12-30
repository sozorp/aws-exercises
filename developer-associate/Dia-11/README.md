# Día 11: Vamos a explorar Amazon ECS (Elastic Container Service) con Fargate

## Escenario

Quieres desplegar una aplicación basada en contenedores utilizando **AWS ECS con Fargate**. La aplicación será un servidor web Nginx básico. Configurarás ECS para que administre el contenedor, y probarás el acceso desde un navegador.

## Objetivo

1. Crear un **Dockerfile** para una aplicación web básica (usando Nginx).
2. Subir la imagen a **Amazon Elastic Container Registry (ECR)**.
3. Configurar un clúster ECS con Fargate para ejecutar la aplicación.
4. Probar el acceso a la aplicación desde un navegador.

## Pasos para Resolver

1. Crear un Dockerfile:

   - En tu directorio de trabajo, crea un archivo llamado Dockerfile con el siguiente contenido:

   ```dockerfile
   FROM nginx:alpine
   COPY index.html /usr/share/nginx/html/index.html
   ```

   - Crea un archivo `index.html` con un mensaje básico:

   ```html
   <html>
     <head>
       <title>Welcome</title>
     </head>
     <body>
       <h1>Hello from ECS Fargate!</h1>
     </body>
   </html>
   ```

2. Construir y Subir la Imagen a ECR:

   - Inicia sesión en ECR:

   ```bash
   aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account-id>.dkr.ecr.<region>.amazonaws.com
   ```

   - Crea un repositorio en ECR:

   ```bash
   aws ecr create-repository --repository-name nginx-web-app
   ```

   - Construye y etiqueta la imagen:

   ```bash
   docker build -t nginx-web-app .
   docker tag nginx-web-app:latest <account-id>.dkr.ecr.<region>.amazonaws.com/nginx-web-app:latest
   ```

   - Sube la imagen al repositorio:

   ```bash
   docker push <account-id>.dkr.ecr.<region>.amazonaws.com/nginx-web-app:latest
   ```

3. Configurar un Clúster ECS con Fargate:

   - Crea un clúster ECS:

   ```bash
   aws ecs create-cluster --cluster-name WebAppCluster
   ```

   - Define una tarea en ECS usando el siguiente archivo JSON (`task-definition.json`):

   ```json
   {
     "family": "nginx-web-app",
     "networkMode": "awsvpc",
     "containerDefinitions": [
       {
         "name": "nginx",
         "image": "<account-id>.dkr.ecr.<region>.amazonaws.com/nginx-web-app:latest",
         "memory": 512,
         "cpu": 256,
         "essential": true,
         "portMappings": [
           {
             "containerPort": 80,
             "hostPort": 80,
             "protocol": "tcp"
           }
         ]
       }
     ],
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "256",
     "memory": "512",
     "executionRoleArn": "arn:aws:iam::<account-id>:role/ecsTaskExecutionRole"
   }
   ```

   - Si no tienes el rol `ecsTaskExecutionRole`, [revisa este link](https://docs.aws.amazon.com/es_es/AmazonECS/latest/developerguide/task_execution_IAM_role.html)

   - Registra la tarea:

   ```bash
   aws ecs register-task-definition --cli-input-json file://task-definition.json
   ```

4. Crear un Servicio ECS:

   - Lanza el servicio ECS basado en la definición de tarea:

   ```bash
   aws ecs create-service \
    --cluster WebAppCluster \
    --service-name NginxService \
    --task-definition nginx-web-app \
    --desired-count 1 \
    --launch-type FARGATE \
    --network-configuration "awsvpcConfiguration={subnets=[<subnet-id>],securityGroups=[<sg-id>],assignPublicIp=ENABLED}"
   ```

   - Asegúrate de que

     - **Grupo de Seguridad**: Permita tráfico entrante en el puerto 80 desde tu IP de origen o cualquier IP (`0.0.0.0/0` para acceso público).
     - **Subnets**: La tarea está en una subnet con acceso a Internet (ya sea una subnet pública o mediante una NAT Gateway).

5. Probar la Aplicación:

   - Obtén la dirección IP pública de la tarea:

   ```bash
   aws ecs list-tasks --cluster WebAppCluster
   aws ecs describe-tasks --cluster WebAppCluster --tasks <task-id>
   ```

   - Accede a la dirección IP pública en tu navegador y verifica que ves el mensaje `Hello from ECS Fargate!`.

## Preguntas de Práctica Asociadas

- ¿Cuáles son las diferencias clave entre ECS y EKS (Elastic Kubernetes Service)?
- ¿Cómo podrías escalar este servicio automáticamente para manejar un aumento de tráfico?

Happy hacking! 🚀
