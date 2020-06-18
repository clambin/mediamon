FROM python:3.7-alpine
MAINTAINER Christophe Lambin <christophe.lambin@gmail.com>

EXPOSE 8080

RUN addgroup -S -g 1000 abc && adduser -S --uid 1000 --ingroup abc abc

WORKDIR /app
COPY Pipfile Pipfile.lock ./

RUN pip install --upgrade pip && \
    pip install pipenv && \
    pipenv install --system --ignore-pipfile

COPY *.py ./
COPY src src/

USER abc
ENTRYPOINT ["/usr/local/bin/python3", "mediamon.py"]
CMD []
