FROM python:alpine

RUN mkdir -p /opt/python-api
WORKDIR /opt/python-api
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
ENV FLASK_APP=app.py
COPY app.py app.py
ENTRYPOINT ["flask", "run", "--host=0.0.0.0"]
