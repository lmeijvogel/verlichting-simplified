FROM python:3

WORKDIR /src

RUN apt-get update && apt-get install tini && rm -rf /var/lib/apt/lists/*

# https://www.twilio.com/docs/usage/tutorials/how-to-set-up-your-python-and-flask-development-environment
RUN pip install --root-user-action=ignore Flask twilio requests python-dotenv gunicorn

ENTRYPOINT ["tini", "--"]
CMD ["python", "app.py"]
