# Día 17: Vamos a trabajar con AWS CloudFront y S3 para crear una distribución de contenido estático con caching y HTTPS

## Escenario

Tu empresa quiere entregar contenido estático (como imágenes, videos o archivos HTML) de forma eficiente a usuarios de todo el mundo. Usarás **Amazon CloudFront** para distribuir los archivos almacenados en un bucket S3, asegurando que el contenido se entregue mediante HTTPS.

## Objetivo

1. Configurar un bucket S3 para almacenar contenido estático.
2. Crear una distribución de **CloudFront** que apunte al bucket S3.
3. Configurar HTTPS usando el certificado gratuito de **AWS Certificate Manager (ACM)**.
4. Probar el acceso al contenido a través de CloudFront.

## Pasos Detallados

### Paso 1: Configurar el Bucket S3

1. Ve a la consola de **S3** y crea un bucket llamado `cloudfront-static-{tu-nombre}`.
2. Sube un archivo HTML o de imagen al bucket, por ejemplo, un archivo index.html`:

   ```html
   <html>
     <head>
       <title>CloudFront Demo</title>
     </head>
     <body>
       <h1>Hello from CloudFront!</h1>
     </body>
   </html>
   ```

3. Configura los permisos del bucket:
   - Asegúrate de que los objetos del bucket sean **públicos**.
   - Habilita la opción **"Static website hosting"** y establece el documento de índice como `index.html`.

### Paso 2: Crear un Certificado en ACM

1. Ve a la consola de **AWS Certificate Manager (ACM)**.
2. Solicita un nuevo certificado público.
3. Si tienes un dominio propio, agrégalo aquí (por ejemplo, `www.mydomain.com`) y verifica la propiedad del dominio.
4. Si no tienes un dominio, puedes omitir este paso y usar el dominio generado por CloudFront.

### Paso 3: Crear la Distribución de CloudFront

1. Ve a la consola de **CloudFront** y crea una nueva distribución.
2. En **Origin Domain**, selecciona el bucket S3 que creaste (`cloudfront-static-{tu-nombre}.s3.amazonaws.com`).
3. Configura el resto de las opciones:

   - **Viewer Protocol Policy**: Redirige todas las solicitudes a HTTPS.
   - **Default Root Object**: Especifica `index.html`.

4. Si configuraste un certificado en ACM, selecciónalo en la sección de **SSL Certificate**.
5. Crea la distribución y espera unos minutos a que se propague.

### Paso 4: Probar el Flujo

1. Una vez que la distribución esté activa, obtén el **Domain Name** de la distribución de CloudFront (por ejemplo, `d1a2b3c4.cloudfront.net`).
2. Accede a la URL desde tu navegador:

   ```bash
   https://d1a2b3c4.cloudfront.net
   ```

3. Verifica que se muestre el contenido del archivo `index.html`.

## Preguntas de Práctica Asociadas

- ¿Cómo podrías restringir el acceso al contenido solo a usuarios autenticados?
- ¿Qué beneficios ofrece CloudFront en comparación con servir contenido directamente desde S3?

Happy hacking! 🚀
