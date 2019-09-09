FROM python:alpine

RUN mkdir -p /opt/python-api
WORKDIR /opt/python-api
COPY requirements.txt requirements.txt
RUN apk add --no-cache mariadb-connector-c-dev mariadb-connector-c ;\
    apk add --no-cache --virtual .build-deps \
        build-base \
        mariadb-dev ;\
    pip install -r requirements.txt ;\
    apk del .build-deps
RUN pip install -r requirements.txt
ENV FLASK_APP=app.py
COPY app.py app.py
ENTRYPOINT ["flask", "run", "--host=0.0.0.0"]
