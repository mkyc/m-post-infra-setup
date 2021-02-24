FROM quay.io/ansible/ansible-runner:stable-2.10-devel

COPY ./resources/requirements.yml .

RUN ansible-galaxy install -r requirements.yml
