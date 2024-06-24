#!/usr/bin/env bash

set -e
set -u
set -o pipefail

function remove_images() {
  images=(
    #old images on Docker HUB
    vulpemventures/bitcoin:latest
    vulpemventures/electrs:latest
    vulpemventures/esplora:latest
    vulpemventures/nigiri-chopsticks:latest
   
    # new images on GH container registry
    ghcr.io/vulpemventures/bitcoin:latest
    ghcr.io/vulpemventures/electrs:latest
    ghcr.io/vulpemventures/esplora:latest
    ghcr.io/vulpemventures/nigiri-chopsticks:latest
    ghcr.io/catalogfi/garden_sol:latest
    ghcr.io/catalogfi/orderbook:latest
    ghcr.io/catalogfi/cobi:latest
  )
  for image in ${images[*]}; do
    if [ "$(docker images -q $image)" != "" ]; then
      docker rmi $image 1>/dev/null
      echo "successfully deleted $image"
    fi
  done
}

##/=====================================\
##|      DETECT PLATFORM      |
##\=====================================/
case $OSTYPE in
darwin*) OS="darwin" ;;
linux-gnu*) OS="linux" ;;
*)
  echo "OS $OSTYPE not supported by the installation script"
  exit 1
  ;;
esac

case $(uname -m) in
amd64) ARCH="amd64" ;;
arm64) ARCH="arm64" ;;
x86_64) ARCH="amd64" ;;
*)
  echo "Architecture $ARCH not supported by the installation script"
  exit 1
  ;;
esac

BIN="/usr/local/bin"

##/=====================================\
##|     CLEAN OLD INSTALLATION |
##\=====================================/

if [ "$(command -v merry)" != "" ]; then
  echo "Merry is already installed and will be deleted."
  # check if Docker is running
  if [ -z "$(docker info 2>&1 >/dev/null)" ]; then
    :
  else
    echo
    echo "Info: when uninstalling an old Merry version Docker must be running."
    echo
    echo "Be sure to start the Docker daemon before launching this installation script."
    echo
    #exit 1
  fi

  echo "Stopping Merry..."
  if [ -z "$(merry stop --delete &>/dev/null)" ]; then
    :
  fi

  echo "Removing Merry..."
  sudo rm -f $BIN/merry
  sudo rm -rf ~/.merry

  echo "Removing local images..."
  remove_images
fi

##/=====================================\
##|     FETCH LATEST RELEASE      |
##\=====================================/
MERRY_URL="https://github.com/catalogfi/merry/releases"
LATEST_RELEASE_URL="$MERRY_URL/latest"

echo "Fetching $LATEST_RELEASE_URL..."

TAG=$(curl -sL -H 'Accept: application/json' $LATEST_RELEASE_URL | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')

echo "Latest release tag = $TAG"

RELEASE_URL="$MERRY_URL/download/$TAG/merry-$OS-$ARCH"

echo "Fetching $RELEASE_URL..."

curl -sL $RELEASE_URL >merry

echo "Moving binary to $BIN..."
sudo mv merry $BIN

echo "Setting binary permissions..."
sudo chmod +x $BIN/merry

echo "Checking for Docker and Docker compose..."
if [ "$(command -v docker)" == "" ]; then
  echo "Warning: Merry uses Docker and it seems not to be installed, check the official documentation for downloading it."
  if [Â "$OS" = "darwin" ]; then
    echo "https://docs.docker.com/v17.12/docker-for-mac/install/#download-docker-for-mac"
  else
    echo "https://docs.docker.com/v17.12/install/linux/docker-ce/ubuntu/"
  fi
fi

echo ""
echo "Merry Catalog installed!"