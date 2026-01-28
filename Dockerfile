#
# Builder image
#

FROM golang:1.24-bookworm AS build-env
ARG SOURCE=*

ADD $SOURCE /src/
WORKDIR /src/

# Unpack any tars, then try to execute a Makefile, but if the SOURCE url is
# just a tar of binaries, then there probably won't be one. Using multiple RUN
# commands to ensure any errors are caught.
RUN find . -name '*.tar.gz' -type f | xargs -rn1 tar -xzf
RUN if [ -f Makefile ]; then make static; fi
RUN cp "$(find . -name 'louketo-proxy' -type f -print -quit)" /louketo-proxy

#
# Actual image
#

FROM gcr.io/distroless/static-debian12

LABEL Name=louketo-proxy \
      Release=https://github.com/louketo/louketo-proxy \
      Url=https://github.com/louketo/louketo-proxy \
      Help=https://github.com/louketo/louketo-proxy/issues

WORKDIR /opt/louketo

COPY templates ./templates
COPY --from=build-env /louketo-proxy ./louketo-proxy

USER nonroot:nonroot
ENTRYPOINT [ "/opt/louketo/louketo-proxy" ]
