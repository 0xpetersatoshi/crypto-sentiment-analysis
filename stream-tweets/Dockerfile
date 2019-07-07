FROM python:3.6-alpine

WORKDIR /usr/src/app
COPY ./requirements.txt .
RUN pip install -r requirements.txt
COPY ./app.py .

ENTRYPOINT [ "python", "app.py" ]
