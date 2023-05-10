FROM registry.access.redhat.com/ubi9/ubi:latest as build

WORKDIR /build
COPY . /build

RUN dnf install -y make go

# Compile evergreen
ENV GOOS=linux
ENV GOARCH=arm64
RUN make

# Copy the binary to a predictable place so next step can find it
RUN mv clients/${GOOS}_${GOARCH}/evergreen /evergreen


# Create a minimal UBI image
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest as runtime
WORKDIR /app

# Copy evergreen into this container
COPY --from=build /evergreen .

CMD ["/app/evergreen"]
