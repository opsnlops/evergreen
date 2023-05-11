FROM registry.access.redhat.com/ubi9/ubi:latest as archive-builder

# Grab build args from the command line
ARG GIT_HASH
ARG GOOS
ARG GOARCH
ARG GOVERSION

WORKDIR /build
COPY . .

# Make an archive of our sources
RUN tar --transform "s,^,evergreen-${GIT_HASH}/,"  -cvzf /evergreen-${GIT_HASH}.tar.gz .



FROM registry.access.redhat.com/ubi9/ubi:latest as build


# Grab build args from the command line
ARG GIT_HASH
ARG GOOS
ARG GOARCH

WORKDIR /build

COPY rpm/evergreen.spec .
COPY --from=archive-builder /evergreen-${GIT_HASH}.tar.gz SOURCES/

# Add make, since we use it to build
RUN dnf update && dnf install -y rpm-build && dnf builddep -y evergreen.spec

# Grab the latest go from Google
#ADD https://go.dev/dl/go${GOVERSION}.${GOOS}-${GOARCH}.tar.gz /
#RUN tar -C /usr/local -xzf /go${GOVERSION}.${GOOS}-${GOARCH}.tar.gz
#ENV PATH=/usr/local/go/bin:${PATH}

# Install bom so we can make an SPDX file
RUN go install sigs.k8s.io/bom/cmd/bom@latest

# Compile evergreen
RUN rpmbuild -ba \
    --define "_git_hash ${GIT_HASH}" \
    --define "_topdir /build" \
    --define "_go_arch ${GOARCH}" \
    --define "_go_os ${GOOS}" \
    evergreen.spec

# rpmbuild -ba --define "_git_hash 52189baf6" --define "_topdir /build" --define "_go_arch arm64" --define "_go_os linux" evergreen.spec


# Create a minimal UBI image
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest as runtime
WORKDIR /app

ARG GIT_HASH

# Copy evergreen into this container
COPY --from=build /build/RPMS/aarch64/evergreen-${GIT_HASH}-1.el9.aarch64.rpm .
COPY --from=build /build/SRPMS/evergreen-${GIT_HASH}-1.el9.src.rpm .

RUN rpm -i evergreen-${GIT_HASH}-1.el9.aarch64.rpm

#CMD ["/usr/bin/evergreen"]
