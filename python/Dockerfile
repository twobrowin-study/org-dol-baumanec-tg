FROM python:3.10-slim-buster

ENV CONNECTION_STRING ''

WORKDIR /python-docker

RUN apt update && apt install -y python3-psycopg2 libpq-dev \
    python3-waitress python3-setuptools python3-pandas \
    python3-sqlalchemy python3-flask python3-xlsxwriter

RUN pip3 install flask pandas XlsxWriter sqlalchemy setuptools

COPY . .

EXPOSE 8080
CMD [ "waitress-serve", "--port=8080" , "--call", "download:create_app"]