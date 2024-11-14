#!/bin/bash

# Actualizar el sistema
apt-get update
apt-get upgrade -y

# Instalar dependencias necesarias para Docker y Go
apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    software-properties-common \
    wget \
    build-essential \
    git \
    make

# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
# Agregar el usuario actual al grupo docker (opcional pero recomendado)
usermod -aG docker $USER
# (Necesitarás cerrar sesión y volver a iniciar sesión para que esto surta efecto)

# Instalar Go 1.23.2
GO_VERSION=1.23.2
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# Configurar las variables de entorno de Go
echo "export GOROOT=/usr/local/go" >> /etc/profile
echo "export GOPATH=\$HOME/go" >> /etc/profile
echo "export PATH=\$GOPATH/bin:\$GOROOT/bin:\$PATH" >> /etc/profile
source /etc/profile


# Instalar golang-migrate
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate

# Clonar el repositorio (reemplaza con la URL de tu repositorio)
git clone https://github.com/japablazatww/songs-searcher.git

# Navegar al directorio del proyecto
cd songs-searcher

# Construir y ejecutar con Make (asumiendo que tu Makefile tiene los targets necesarios)
make up_build

# Asegura tener la DB creada
make create_db

# Ejecuta las migraciones
make migrate_db_up