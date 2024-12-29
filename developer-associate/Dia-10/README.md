# Día 10: Vamos a crear un flujo de trabajo más complejo que incluya AWS CloudFormation para definir y desplegar infraestructura como código

## Escenario

Necesitas desplegar una infraestructura básica que incluya:

1. Un bucket S3 para almacenar datos.
2. Una función Lambda que procese eventos desde ese bucket.
3. Una cola SQS para manejar errores en caso de que la Lambda falle.

Usarás **CloudFormation** para crear y gestionar todos estos recursos.

## Objetivo

1. Escribir una plantilla de CloudFormation que:

   - Cree un bucket S3 llamado `cf-example-bucket-{tu-nombre}`.
   - Cree una función Lambda que procese eventos desde el bucket.
   - Configure un **Dead-Letter Queue (DLQ)** en SQS para los fallos de Lambda.

2. Desplegar la plantilla.
3. Probar el flujo enviando archivos al bucket y verificando los resultados.

## Pasos para Resolver

1. Escribir la Plantilla de CloudFormation:

   - Crea un archivo llamado `template.yaml` con el siguiente contenido:

   ```yaml
   AWSTemplateFormatVersion: "2010-09-09"
   Resources:
     S3Bucket:
       Type: AWS::S3::Bucket
       Properties:
         BucketName: cf-example-bucket-yourname
         NotificationConfiguration:
           LambdaConfigurations:
             - Event: s3:ObjectCreated:*
               Function: !GetAtt LambdaFunction.Arn
       DependsOn: LambdaPermission

     SQSQueue:
       Type: AWS::SQS::Queue
       Properties:
         QueueName: cf-example-dlq-yourname

     LambdaExecutionRole:
       Type: AWS::IAM::Role
       Properties:
         AssumeRolePolicyDocument:
           Version: "2012-10-17"
           Statement:
             - Effect: Allow
               Principal:
                 Service: lambda.amazonaws.com
               Action: sts:AssumeRole
         Policies:
           - PolicyName: LambdaS3Policy
             PolicyDocument:
               Version: "2012-10-17"
               Statement:
                 - Effect: Allow
                   Action:
                     - s3:PutObject
                     - s3:GetObject
                   Resource: "arn:aws:s3:::cf-example-bucket-yourname/*"
                 - Effect: Allow
                   Action:
                     - sqs:SendMessage
                   Resource: !GetAtt SQSQueue.Arn

     LambdaFunction:
       Type: AWS::Lambda::Function
       Properties:
         Handler: index.handler
         Runtime: python3.9
         Role: !GetAtt LambdaExecutionRole.Arn
         Code:
           ZipFile: |
             import boto3
             import json
             import os

             def handler(event, context):
                 try:
                     print("Event: ", event)
                     # Simulación de procesamiento
                     print("Processing S3 event...")
                 except Exception as e:
                     # Simula un error para que el mensaje vaya al DLQ
                     print(f"Error: {str(e)}")
                     raise e
         DeadLetterConfig:
           TargetArn: !GetAtt SQSQueue.Arn
         Environment:
           Variables:
             LOG_LEVEL: ERROR

     LambdaPermission:
       Type: AWS::Lambda::Permission
       Properties:
         FunctionName: !Ref LambdaFunction
         Action: lambda:InvokeFunction
         Principal: s3.amazonaws.com
         SourceArn: "arn:aws:s3:::cf-example-bucket-yourname"
   ```

2. Desplegar la Plantilla:

   - Sube la plantilla a CloudFormation usando AWS CLI:

   ```bash
   aws cloudformation create-stack --stack-name ExampleStack --template-body file://template.yaml --capabilities CAPABILITY_NAMED_IAM
   ```

3. Probar el Flujo:

   - Sube un archivo al bucket S3:

   ```bash
   aws s3 cp test-file.txt s3://cf-example-bucket-yourname/
   ```

   - Verifica:
     - Los logs de la función Lambda en CloudWatch.
     - Que los errores de la Lambda (si los hay) se envían a la cola SQS (`cf-example-dlq-yourname`).

## Preguntas de Práctica Asociadas

- ¿Cómo usarías parámetros en CloudFormation para hacer la plantilla más reutilizable?
- ¿Cómo manejarías actualizaciones en una pila de CloudFormation sin interrumpir los servicios?

Happy hacking! 🚀
