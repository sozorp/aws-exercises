#!/bin/bash

# Verificar que se hayan proporcionado los dos parámetros
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Uso: $0 <archivo_go> <ruta_salida_zip>"
  exit 1
fi

# Parámetros de entrada
GO_FILE=$1
OUTPUT_DIR=$2
OUTPUT_ZIP="$OUTPUT_DIR/bootstrap.zip"

# Verificar que el archivo Go existe
if [ ! -f "$GO_FILE" ]; then
  echo "El archivo Go especificado no existe: $GO_FILE"
  exit 1
fi

# Verificar que el directorio de salida existe
if [ ! -d "$OUTPUT_DIR" ]; then
  echo "El directorio de salida especificado no existe: $OUTPUT_DIR"
  exit 1
fi

# Eliminar el archivo .zip si ya existe en la ruta de salida
rm -f "$OUTPUT_ZIP"

# Configurar GOOS y GOARCH para AWS Lambda
export GOOS=linux
export GOARCH=amd64

# Construir el binario
go build -o bootstrap "$GO_FILE"

# Verificar que la construcción fue exitosa
if [ $? -ne 0 ]; then
  echo "Error al construir el binario"
  exit 1
fi

# Crear el archivo ZIP
zip "$OUTPUT_ZIP" bootstrap

# Limpiar el binario
rm bootstrap

# Verificar que el empaquetado fue exitoso
if [ $? -ne 0 ]; then
  echo "Error al crear el archivo ZIP"
  exit 1
fi

echo "Función empaquetada exitosamente: $OUTPUT_ZIP"
