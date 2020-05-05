# sphinx-tribes

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/sphinx-tribes.png)

Decentralized message broker for public groups in Sphinx. Anyone can run a **sphinx-tribes** server, to route messages for large groups in Sphinx. 

**sphinx-tribes** clients can be anything from **sphinx-relay** nodes, to apps, websites, or IoT devices.

### Architecture

Under the hood, sphinx-tribes uses an MQTT broker to route messages to subscribing nodes. Message topics always have two parts: `{receiverPubKey}/{groupUUID}`. Only the owner of the group is allowed to publish to it: all messages from group members must be submitted to the owner as an LND keysend payment. the `groupUUID` is signed by the owner.

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/tribes.jpg)

### Authentication

Authentication is handled by the sphinx-tribes microservice. 

### run

**build**
docker build --no-cache -t sphinx-tribes .

**restart**
docker stop sphinx-tribes && docker rm sphinx-tribes && docker create -p 0.0.0.0:80:5002/tcp -p 0.0.0.0:1883:1883/tcp --name sphinx-tribes --restart on-failure sphinx-tribes:latest && docker cp config.json sphinx-tribes:/config.json && docker start sphinx-tribes

docker logs sphinx-tribes --since 10m --follow