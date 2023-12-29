FROM alpine:3.19 AS deploy-image

ENV USER=pg-health-check
ENV UID=1000
ENV GID=1000

WORKDIR $HOME

RUN apk update                                                                 \
    && apk add --no-cache                                                      \
         bash                                                                  \
         ca-certificates                                                       \
         curl                                                                  \
         git                                                                   \
         sudo                                                                  \
         tzdata                                                                \
    && addgroup -g "$GID" -S "$USER"                                           \
    && adduser                                                                 \
       --disabled-password                                                     \
       --gecos ""                                                              \
       --home "$HOME"                                                          \
       --ingroup "$USER"                                                       \
       --no-create-home                                                        \
       --uid "$UID" "$USER"                                                    \
    && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER             \
    && chmod 0440 /etc/sudoers.d/$USER

USER ${UID}:${GID}

COPY --chown=${UID}:${GID} ./build-file/default-config.yaml /etc/pg-health-check/config.yaml
COPY pg-health-check /usr/local/bin/pg-health-check

EXPOSE 8080
CMD [ "/usr/local/bin/pg-health-check", "-c", "/etc/pg-health-check/config.yaml" ]