#!/bin/bash

# Verificar que se haya proporcionado el nombre del archivo de salida
if [ -z "$1" ]; then
  echo "No se ha proporcionado el nombre del archivo de salida"
  exit 1
fi

# Eliminar .zip si ya existe
rm -f bootstrap.zip

# Configurar GOOS y GOARCH para AWS Lambda
export GOOS=linux
export GOARCH=amd64

# Construir el binario
go build -o bootstrap $1

# Verificar que la construcción fue exitosa
if [ $? -ne 0 ]; then
  echo "Error al construir el binario"
  exit 1
fi

# Crear el archivo ZIP
zip bootstrap.zip bootstrap

# Verificar que el empaquetado fue exitoso
if [ $? -ne 0 ]; then
  echo "Error al crear el archivo ZIP"
  exit 1
fi

# Limpiar el binario
rm bootstrap

echo "Función empaquetada exitosamente: bootstrap.zip"
