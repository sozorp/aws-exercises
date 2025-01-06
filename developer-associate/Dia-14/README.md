# Día 14: Vamos a trabajar con AWS CodePipeline, CodeBuild, y CodeDeploy para implementar un flujo de CI/CD completo

## Escenario

Tu equipo de desarrollo necesita un flujo de integración y entrega continua (CI/CD) que automatice el proceso de construcción, pruebas y despliegue de una aplicación web basada en contenedores.

El objetivo es

1. Usar **AWS CodePipeline** para orquestar el flujo de CI/CD.
2. Configurar **CodeBuild** para compilar la imagen de contenedor y subirla a **Amazon ECR**.
3. Usar **CodeDeploy** para desplegar la aplicación en un clúster ECS con Fargate.
4. Un flujo que detecte automáticamente cambios en el repositorio y realice la implementación sin intervención manual.

## Pasos para Resolver

1. Crear un Repositorio en Github:

   - Ve a **Github** y crea un nuevo repositorio llamado `web-app-repo`.
   - Clona el repositorio localmente y agrega un archivo `Dockerfile`y un archivo`index.html` para una aplicación Nginx básica.

   **Dockerfile**:

   ```dockerfile
   FROM nginx:alpine
   COPY index.html /usr/share/nginx/html/index.html
   ```

   **index.html**:

   ```html
   <html>
     <head>
       <title>Welcome</title>
     </head>
     <body>
       <h1>Hello from CodePipeline CI/CD!</h1>
     </body>
   </html>
   ```

   - Haz un commit y sube los cambios al repositorio Github:

   ```bash
   git add .
   git commit -m "Initial commit"
   git push origin main
   ```

2. Configurar un Repositorio en ECR:

   - Crea un repositorio en **ECR** llamado `web-app`.

3. Configurar CodeBuild:

   - Crea un proyecto de **CodeBuild** con los siguientes detalles:
     - Fuente: el repositorio en Github `web-app-repo`.
     - Entorno: Usa una imagen administrada de **Ubuntu** con Docker preinstalado.
   - Agrega un archivo `buildspec.yml` al repositorio con las siguientes instrucciones de construcción:

   **buildspec.yml**:

   ```yaml
   version: 0.2

   phases:
     pre_build:
       commands:
         - echo Logging in to Amazon ECR...
         - $(aws ecr get-login --no-include-email --region $AWS_DEFAULT_REGION)
     build:
       commands:
         - echo Building the Docker image...
         - docker build -t web-app .
         - docker tag web-app:latest <account-id>.dkr.ecr.<region>.amazonaws.com/web-app:latest
     post_build:
       commands:
         - echo Pushing the Docker image...
         - docker push <account-id>.dkr.ecr.<region>.amazonaws.com/web-app:latest
         - echo Build completed on `date`
   ```

   - Haz commit y push del archivo `buildspec.yml` al repositorio.

4. Configurar ECS y CodeDeploy:

   A. **Crea un clúster ECS con Fargate**.

   1. Ve a **Amazon ECS** y crea un nuevo clúster con el tipo de lanzamiento **Fargate**.
   2. Crea una nueva tarea de definición que utilice:

      - Contenedor: `web-app`.
      - Imagen: `<account-id>.dkr.ecr.<region>.amazonaws.com/web-app:latest`.
      - Puerto mapeado: 80.

   3. Crea un servicio basado en la definición de tarea que acabas de crear.

   B. **Configurar CodeDeploy**

   1. Ve a **AWS CodeDeploy** y crea una nueva aplicación para ECS.
   2. Crea un nuevo grupo de despliegue y selecciona el clúster y servicio ECS que configuraste.
   3. Configura el archivo `appspec.yml` y el script de despliegue en el repositorio:

   **appspec.yml**:

   ```yaml
   version: 0.0
   Resources:
     - TargetService:
         Type: AWS::ECS::Service
         Properties:
           TaskDefinition: "web-app"
           LoadBalancerInfo:
             ContainerName: "web-app"
             ContainerPort: 80

       - Define un servicio ECS que use la imagen `web-app` desde ECR.
       - Configura **CodeDeploy** para ECS:
         - Crea una aplicación y un grupo de despliegue en CodeDeploy.
         - Configura el grupo de despliegue para que use el clúster y servicio ECS que creaste.
   ```

5. Configurar CodePipeline:

   1. Crea una nueva **pipeline** en **CodePipeline**:

   - **Source**: Configura el repositorio `web-app-repo` de GitHub.
   - **Build**: Usa el proyecto de **CodeBuild** que configuraste.
   - **Deploy**: Configura la acción de despliegue con **CodeDeploy** y selecciona la aplicación y grupo de despliegue creados.

## Prueba del Flujo Completo

1. Realiza un cambio en el archivo `index.html`, por ejemplo, cambia el mensaje a:

   ```html
   <h1>Hello from CI/CD Pipeline - Updated!</h1>
   ```

2. Haz un commit y push al repositorio.
3. Verifica que CodePipeline:

   - Detecta automáticamente el cambio.
   - Inicia el proceso de construcción en CodeBuild.
   - Despliega la nueva versión de la aplicación en ECS mediante CodeDeploy.

4. Accede a la aplicación desde el navegador y verifica que el mensaje actualizado se refleja correctamente.

## Consideraciones

1. **Roles y permisos**: Asegúrate de que los roles de servicio de CodePipeline, CodeBuild y CodeDeploy tengan los permisos necesarios para interactuar con ECR, ECS y S3.

2. **Logs y monitoreo**:

   - Habilita **CloudWatch Logs** para CodeBuild y ECS para depurar errores si el flujo falla.

## Preguntas de Práctica Asociadas

- ¿Cómo manejarías versiones anteriores de la aplicación si necesitas hacer un rollback?
- ¿Qué beneficios ofrece CodePipeline en comparación con herramientas externas de CI/CD?

Happy hacking! 🚀
