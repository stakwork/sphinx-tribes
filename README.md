# sphinx-tribes

Decentralized message broker for public groups in Sphinx. Anyone can run a sphinx-tribes server, to route messages for large groups in Sphinx. 

### Architecture

Under the hood, sphinx-tribes uses an MQTT broker to route messages to subscribing nodes. Message topics always have two parts: `{receiverPubKey}/{groupUUID}`. Only the owner of the group is allowed to publish to it: all messages from group members must be submitted to the owner as an LND keysend payment. the `groupUUID` is signed by the owner.

### Authentication

Authentication is handled by the sphinx-auth microservice. 