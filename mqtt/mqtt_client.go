package mqtt

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/big"
	"os"
	"time"

	"crypto/rand"
	"encoding/hex"
)

var (
	opts          *mqtt.ClientOptions
	thePrivateKey []byte

	// signedMsgPrefix is a special prefix that we'll prepend to any
	// messages we sign/verify. We do this to ensure that we don't
	// accidentally sign a sighash, or other sensitive material. By
	// prepending this fragment, we mind message signing to our particular
	// context.
	signedMsgPrefix = []byte("Lightning Signed Message:")
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Messages published:\n  [topic= %s\n  payload= %s\n", msg.Topic(), msg.Payload())
}

func Init() {
	//mqtt.DEBUG = log.New(os.Stdout, "[mqtt:debug]", 0)
	mqtt.ERROR = log.New(os.Stdout, "[mqtt:error]", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[mqtt:CRITICAL]", 0)
	//mqtt.WARN = log.New(os.Stdout, "[mqtt:warn]", 0)

	mqttScheme := getEnv("MQTT_SCHEME", "tcp")
	mqttPort := getEnv("MQTT_PORT", "1883")
	mqttHost := getEnv("MQTT_HOSTNAME", "host.docker.internal")
	mqttClientID, _ := randomClientName()

	password, mqttPublicKey, err := signTimestamp()
	if err != nil {
		panic(err)
	}

	brokerUri := fmt.Sprintf("%s://%s:%s", mqttScheme, mqttHost, mqttPort)

	opts = mqtt.NewClientOptions()
	opts.AddBroker(brokerUri)
	opts.SetClientID(mqttClientID)

	opts.SetUsername(mqttPublicKey)
	opts.SetPassword(password)
	opts.SetCleanSession(false)
	opts.SetOnConnectHandler(func(c mqtt.Client) { connectedToLocalCallback(c, mqttPublicKey) })
	opts.SetDefaultPublishHandler(f)
}

func Client() mqtt.Client {
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("=> connection error:", token.Error())
	}
	return client
}

func connectedToLocalCallback(c mqtt.Client, pubkey string) {
	c.Subscribe(pubkey+"/#", 0, func(c mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
	pwd, pubkey, err := signTimestamp()
	if err != nil {
		log.Println("there was an error:", err)
	}
	c.Publish(pubkey+"/"+pwd, 0, false, "hi")
}

func signTimestamp() (string, string, error) {
	if thePrivateKey == nil {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return "", "", err
		}
		thePrivateKey = b
	}

	p, pub := btcec.PrivKeyFromBytes(btcec.S256(), thePrivateKey)
	pubHex := hex.EncodeToString(pub.SerializeCompressed())

	t := time.Now().Unix()
	timeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(timeBuf, uint32(t))
	sig, err := Sign(timeBuf, p)
	if err != nil {
		panic(err)
	}

	pwdBuf := append(timeBuf, sig...)
	return base64.URLEncoding.EncodeToString(pwdBuf), pubHex, nil
}

// NodeSigner is an implementation of the MessageSigner interface backed by the
// identity private key of running lnd node.
type nodeSigner struct {
	privateKey *btcec.PrivateKey
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
		privateKey: priv,
	}
}

// Sign ...
func Sign(msg []byte, privKey *btcec.PrivateKey) ([]byte, error) {

	if privKey == nil || msg == nil {
		return nil, errors.New("bad")
	}
	msg = append(signedMsgPrefix, msg...)
	digest := chainhash.DoubleHashB(msg)
	sigBytes, err := btcec.SignCompact(btcec.S256(), privKey, digest, true)
	if err != nil {
		return nil, err
	}
	return sigBytes, nil
}

func Publish(c mqtt.Client, topic string, payload string, wait bool) {
	token := c.Publish(topic, 0, false, payload)
	if wait {
		token.Wait()
	}
}

func Disconnect(c mqtt.Client) {
	c.Disconnect(250)
}

func randomClientName() (string, error) {
	const available = "0123456789abcdefghijklmnopqrstuvwxyz-"
	r := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(available))))
		if err != nil {
			return "", err
		}
		r[i] = available[n.Int64()]
	}

	return "sphinx-tribes-" + string(r), nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
