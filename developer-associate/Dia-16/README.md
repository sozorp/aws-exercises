# Día 16: Vamos a trabajar con AWS Systems Manager (SSM) y EC2 para automatizar tareas administrativas y mejorar la gestión de servidores

## Escenario

Tu empresa tiene múltiples instancias EC2 que necesitan ser administradas de manera eficiente. Quieres automatizar la instalación de software y la ejecución de comandos en estas instancias sin tener que iniciar sesión manualmente en cada una.

Vas a usar **AWS Systems Manager** para:

1. Ejecutar comandos remotos en instancias EC2.
2. Automatizar la instalación de un servidor web Apache.
3. Configurar una tarea recurrente para monitorear el estado de las instancias.

## Objetivo

1. Configurar **AWS Systems Manager Agent (SSM Agent)** en una instancia EC2.
2. Usar **Run Command** de Systems Manager para instalar Apache.
3. Crear una **tarea automatizada** que envíe un reporte diario del estado de las instancias EC2.

## Pasos Detallados

### Paso 1: Configurar una Instancia EC2 con SSM Agent

1. Ve a la consola de **EC2** y lanza una nueva instancia.

   - Usa una AMI de Amazon Linux 2 o Ubuntu, que ya tienen el **SSM Agent** preinstalado.
   - Asegúrate de que la instancia tenga un **role de IAM** con la política `AmazonEC2RoleforSSM` adjunta.

2. Verifica que el **SSM Agent** esté ejecutándose en la instancia.

   - Si usas Amazon Linux 2:

   ```bash
   sudo systemctl status amazon-ssm-agent
   ```

### Paso 2: Usar Run Command para Instalar Apache

1. Ve a la consola de **Systems Manager** y selecciona la opción **Run Command**.
2. Haz clic en **Run a command**.
3. Elige el comando predefinido `AWS-RunShellScript`.
4. Selecciona tu instancia EC2.
5. En la sección de comandos, ingresa lo siguiente para instalar Apache y arrancarlo:

   ```bash
   sudo yum update -y
   sudo yum install -y httpd
   sudo systemctl start httpd
   sudo systemctl enable httpd
   ```

6. Ejecuta el comando y verifica que no haya errores.
7. Abre el navegador y accede a la dirección IP pública de la instancia para confirmar que el servidor Apache está en funcionamiento.

### Paso 3: Configurar una Tarea Automatizada con Systems Manager

1. Ve a la consola de **Systems Manager** y selecciona **Automation**.
2. Crea una nueva automatización.
3. Usa el documento predefinido `AWS-CreateSnapshot`.
4. Configura la automatización para que cree un snapshot diario de los volúmenes adjuntos a la instancia EC2.
5. Activa la opción de recurrencia y selecciona una frecuencia diaria.

## Prueba del Flujo Completo

1. Verifica que puedes ejecutar comandos en la instancia EC2 usando Systems Manager sin iniciar sesión manualmente.
2. Asegúrate de que el servidor Apache esté funcionando.
3. Confirma que la tarea de snapshot automatizado se ejecuta correctamente y puedes ver los snapshots en la consola de EC2.

## Preguntas de Práctica Asociadas

- ¿Qué beneficios ofrece Systems Manager frente a una conexión manual vía SSH?
- ¿Cómo podrías usar Systems Manager para ejecutar comandos en múltiples instancias al mismo tiempo?

Happy hacking! 🚀
