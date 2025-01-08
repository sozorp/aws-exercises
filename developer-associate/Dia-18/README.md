# Día 18: Vamos a trabajar con AWS Elastic Load Balancer (ELB), Auto Scaling Groups, y EC2

## Escenario

Tu aplicación web debe manejar un tráfico variable y estar disponible en todo momento. Para lograrlo, vas a configurar un **Auto Scaling Group** que lanza instancias EC2 detrás de un **Application Load Balancer (ALB)**. La arquitectura debe ser capaz de escalar automáticamente según la carga y balancear el tráfico entre las instancias.

## Objetivo

1. Crear una **plantilla de lanzamiento** para instancias EC2.
2. Configurar un **Application Load Balancer** que distribuya el tráfico.
3. Configurar un **Auto Scaling Group** que escale dinámicamente según la carga.
4. Probar la escalabilidad lanzando una carga artificial sobre la aplicación.

## Pasos Detallados

### Paso 1: Crear una Plantilla de Lanzamiento

1. Ve a la consola de **EC2** y selecciona **Launch Templates**.
2. Crea una nueva plantilla de lanzamiento con los siguientes detalles:

   - AMI: Amazon Linux 2.
   - Tipo de instancia: `t2.micro`.
   - Configura un script de datos de usuario (user data) que instale Apache y configure una página web básica:

   ```bash
   #!/bin/bash
   sudo yum update -y
   sudo yum install -y httpd
   sudo systemctl start httpd
   sudo systemctl enable httpd
   echo "Hello from $(hostname -f)" > /var/www/html/index.html
   ```

3. Guarda la plantilla de lanzamiento.

### Paso 2: Crear un Application Load Balancer

1. Ve a la consola de **EC2** y selecciona **Load Balancers**.
2. Crea un nuevo **Application Load Balancer** con las siguientes configuraciones:

   - Tipo: **Internet-facing**.
   - Listeners: HTTP en el puerto 80.
   - Selecciona al menos dos subnets en diferentes zonas de disponibilidad.

3. Crea un **target group** para el ALB:

   - Tipo de destino: Instancias.
   - Puerto: 80.
   - Registra las instancias que lance el Auto Scaling Group en este target group.

### Paso 3: Crear un Auto Scaling Group

1. Ve a la consola de **Auto Scaling Groups** y crea un nuevo grupo.
2. Usa la plantilla de lanzamiento que creaste en el paso 1.
3. Configura el tamaño del grupo:

   - **Desired capacity**: 2.
   - **Min capacity**: 1.
   - **Max capacity**: 4.

4. Adjunta el **Application Load Balancer** y el target group creado en el paso 2.
5. Configura una política de escalado automático basada en el uso de CPU:

   - Escalar cuando el uso de CPU supere el 60%.
   - Escalar hacia abajo cuando el uso de CPU baje del 30%.

### Paso 4: Probar el Flujo Completo

1. Obtén la URL del **Application Load Balancer** y accede desde tu navegador:

    ```bash
    http://<alb-dns-name>
    ```

2. Usa **Apache Benchmark** o una herramienta similar para simular una carga en la aplicación:

    ```bash
    ab -n 10000 -c 100 http://<alb-dns-name>/
    ```

3. Observa cómo el Auto Scaling Group lanza nuevas instancias para manejar el aumento de carga.
4. Verifica en la consola de EC2 que el tráfico se balancea entre las instancias y que el Auto Scaling Group reduce el número de instancias cuando la carga disminuye.

## Preguntas de Práctica Asociadas

- ¿Qué diferencia hay entre un Application Load Balancer y un Network Load Balancer?
- ¿Cómo manejarías una implementación en múltiples regiones usando Auto Scaling y ELB?

Happy hacking! 🚀
