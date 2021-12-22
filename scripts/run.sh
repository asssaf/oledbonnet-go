: ${IMAGE_NAME:=asssaf/oledbonnet:latest}
docker run --rm -it --privileged --device /dev/gpiomem --device /dev/i2c-1 "$IMAGE_NAME" $*
