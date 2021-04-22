# Container image that runs your code
FROM cossacklabs/ci-py-go-themis:latest

# RUN apt update
# RUN apt install wget gnupg2 build-essential openssl-devl -y
# RUN wget -qO - https://pkgs-ce.cossacklabs.com/gpg | apt-key add -
# RUN echo "deb https://pkgs-ce.cossacklabs.com/stable/ubuntu focal main" | tee /etc/apt/sources.list.d/cossacklabs.list
# RUN echo /etc/apt/sources.list.d/cossacklabs.list
# RUN apt update
# RUN apt-cache search libthemis
# RUN apt install libthemis

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]
