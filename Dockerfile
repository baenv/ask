FROM jetpackio/devbox:latest

WORKDIR /code
USER root:root
RUN mkdir -p /code && chown ${DEVBOX_USER}:${DEVBOX_USER} /code
USER ${DEVBOX_USER}:${DEVBOX_USER}

COPY --chown=${DEVBOX_USER}:${DEVBOX_USER} . .

RUN devbox install

RUN devbox run install

RUN devbox run build

CMD ["devbox", "run", "start"]
