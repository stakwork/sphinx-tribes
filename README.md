# sphinx-tribes

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/sphinx-tribes.png)

Decentralized message broker for public groups in Sphinx. Anyone can run a **sphinx-tribes** server, to route group messages.

**sphinx-tribes** clients can be **sphinx-relay** nodes, apps, websites, or IoT devices.

### How

**sphinx-tribes** is an MQTT broker that any node can subscribe to. Message topics always have two parts: `{receiverPubKey}/{groupUUID}`. Only the owner of the group is allowed to publish to it: all messages from group members must be submitted to the owner as an Lightning keysend payment. The group `uuid` is a timestamp signed by the owner.

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/tribes.jpg)

### Authentication

Authentication is handled by [sphinx-auth](https://github.com/stakwork/sphinx-auth)

### build

docker build --no-cache -t sphinx-tribes .
