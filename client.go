package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

var client mqtt.Client

func connectClient(port string) {
	pwd, pubkey := signTimestamp()
	connectToLocal(port, pubkey, pwd)
}

func connectToLocal(port, pubkey, pwd string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://0.0.0.0:" + port)
	opts.SetClientID("local-client")

	opts.SetUsername(pubkey)
	opts.SetPassword(pwd)
	opts.SetCleanSession(false)
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		connectedToLocalCallback(c, pubkey)
	})

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("=> connection error:", token.Error())
	}
}

func connectedToLocalCallback(c mqtt.Client, pubkey string) {
	fmt.Println(" ===> CONNECTED!")

	client.Subscribe(pubkey+"/#", 0, func(c mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		// topic := msg.Topic()
		// if string([]rune(topic)[0]) == "/" {
		// 	topic = topic[1:]
		// }

		// var m map[string]interface{}
		// err := json.Unmarshal(msg.Payload(), &m)
		// if err != nil {
		// 	m = map[string]interface{}{"data": string(msg.Payload())}
		// }

		// out, _ := json.Marshal(m)
		// //fmt.Println(string(out))
		// token := flespiClient.Publish(topic, 0, false, string(out))
		// token.Wait()
		// if token.Error() != nil {
		// 	fmt.Println(token.Error())
		// }
	})

	pwd, pubkey := signTimestamp()
	client.Publish(pubkey+"/"+pwd, 0, false, "hi")
}

/*
CLIENT connects
*/

func signTimestamp() (string, string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("no .env file")
	}
	privKey := os.Getenv("PRIV_KEY")
	thePrivKey, err := base64.URLEncoding.DecodeString(privKey)
	if err != nil {
		thePrivKey = zekesPrivKey
	}

	priv, pub := btcec.PrivKeyFromBytes(btcec.S256(), thePrivKey)
	pubBase64 := base64.URLEncoding.EncodeToString(pub.SerializeCompressed())
	signer = newNodeSigner(priv)

	time := time.Now().Unix()
	timeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(timeBuf, uint32(time))
	sig := Sign(timeBuf, priv)

	pwdBuf := append(timeBuf, sig...)
	return base64.URLEncoding.EncodeToString(pwdBuf), pubBase64
}

var zekesPrivKey = []byte{
	0x2c, 0xd8, 0x07, 0xc9, 0x7f, 0x0e, 0x00, 0xaf,
	0x1a, 0x1f, 0xc3, 0x32, 0x8f, 0xa7, 0x63, 0xa9,
	0x26, 0x97, 0x23, 0xc8, 0xdb, 0x8f, 0xac, 0x4f,
	0x93, 0xaf, 0x52, 0xdb, 0x18, 0x6d, 0x6e, 0x90,
}

var signer *nodeSigner

// NodeSigner is an implementation of the MessageSigner interface backed by the
// identity private key of running lnd node.
type nodeSigner struct {
	privKey *btcec.PrivateKey
}

// NewNodeSigner creates a new instance of the NodeSigner backed by the target
// private key.
func newNodeSigner(key *btcec.PrivateKey) *nodeSigner {
	priv := &btcec.PrivateKey{}
	priv.Curve = btcec.S256()
	priv.PublicKey.X = key.X
	priv.PublicKey.Y = key.Y
	priv.D = key.D
	return &nodeSigner{
		privKey: priv,
	}
}

// Sign ...
func Sign(msg []byte, privKey *btcec.PrivateKey) []byte {

	msg = append(signedMsgPrefix, msg...)
	digest := chainhash.DoubleHashB(msg)
	// btcec.S256(), sig, digest

	sigBytes, err := btcec.SignCompact(btcec.S256(), privKey, digest, true)
	if err != nil {
		return nil
	}

	// sig := base64.URLEncoding.EncodeToString(sigBytes)
	return sigBytes
}

var (
	// signedMsgPrefix is a special prefix that we'll prepend to any
	// messages we sign/verify. We do this to ensure that we don't
	// accidentally sign a sighash, or other sensitive material. By
	// prepending this fragment, we mind message signing to our particular
	// context.
	signedMsgPrefix = []byte("Lightning Signed Message:")
)
