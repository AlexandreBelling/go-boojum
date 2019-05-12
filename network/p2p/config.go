package p2p

import (
	"fmt"
	"io/ioutil"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	crypto		"github.com/libp2p/go-libp2p-crypto"
	libp2p		"github.com/libp2p/go-libp2p"
)

func init() {

	viper.SetDefault(libp2pClientHostViperKey, libp2pClientHostDefault)
	_ = viper.BindEnv(libp2pClientHostViperKey, libp2pClientHostEnv)

	viper.SetDefault(libp2pClientPortViperKey, libp2pClientPortDefault)
	_ = viper.BindEnv(libp2pClientPortViperKey, libp2pClientPortEnv)

	viper.SetDefault(libp2pClientPrivateKeyFileViperKey, libp2pClientPrivateKeyFileDefault)
	_ = viper.BindEnv(libp2pClientPrivateKeyFileViperKey, libp2pClientPrivateKeyFileEnv)

}

var (

	libp2pClientHostFlag				= "libp2p-client-host" 
	libp2pClientHostViperKey			= "libp2p.client.host"
	libp2pClientHostDefault				= "127.0.0.1"
	libp2pClientHostEnv					= "LIBP2P_CLIENT_HOST"

	libp2pClientPortFlag				= "libp2P-client-port"
	libp2pClientPortViperKey			= "libp2P.client.port"
	libp2pClientPortDefault				= "9000"
	libp2pClientPortEnv					= "LIBP2P_CLIENT_PORT"

	libp2pClientPrivateKeyFileFlag		= "libp2p-client-private-key-file"
	libp2pClientPrivateKeyFileViperKey	= "libp2p.client.private.key.file"
	libp2pClientPrivateKeyFileDefault	= ""
	libp2pClientPrivateKeyFileEnv		= "LIBP2P_CLIENT_PRIVATE_KEY_FILE"

)

// HostFlags initialize the flags for the libp2p host listenning address.
func HostFlags(f *pflag.FlagSet) {

	desc := fmt.Sprintf(`Host the p2p is listening from. 
Environment variable %q. 
Defaults to %q.`, libp2pClientHostEnv, libp2pClientHostDefault)
	f.String(libp2pClientHostFlag, libp2pClientHostDefault, desc)
	_ = viper.BindPFlag(libp2pClientHostViperKey, f.Lookup(libp2pClientHostFlag))
}

// PortFlags initialize the flags for the libp2p port listenning address.
func PortFlags(f *pflag.FlagSet) {

	desc := fmt.Sprintf(`Port the p2p is listening from. 
Environment variable %q. 
Defaults to %q.`, libp2pClientPortEnv, libp2pClientPortDefault)
	f.String(libp2pClientPortFlag, libp2pClientPortDefault, desc)
	_ = viper.BindPFlag(libp2pClientPortViperKey, f.Lookup(libp2pClientPortFlag))
}

// PrivateKeyFileFlags initialize the flags for private key used by the client to prove his identity.
func PrivateKeyFileFlags(f *pflag.FlagSet) {

	desc := fmt.Sprintf(`PrivateKeyFile the p2p is listening from. 
Environment variable %q. 
Defaults to %q.`, libp2pClientPrivateKeyFileEnv, libp2pClientPrivateKeyFileDefault)
	f.String(libp2pClientPrivateKeyFileFlag, libp2pClientPrivateKeyFileDefault, desc)
	_ = viper.BindPFlag(libp2pClientPrivateKeyFileViperKey, f.Lookup(libp2pClientPrivateKeyFileFlag))
}

func initFlags(f *pflag.FlagSet) {

	HostFlags(f)
	PortFlags(f)
	PrivateKeyFileFlags(f)
}

func readPrivateKey() (crypto.PrivKey, error) {

	privFile := viper.GetString(libp2pClientPrivateKeyFileViperKey)
	if privFile == "" {
		return nil, fmt.Errorf("Unable to set identity to p2p node. Did not specify a private key file")
	}

	privSlice, err := ioutil.ReadFile(privFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read private key file. %q", err.Error())
	}

	priv, err := crypto.UnmarshalSecp256k1PrivateKey(privSlice)
	if err != nil {
		return nil, fmt.Errorf(`Could not unmarshal private key %q`, err.Error())
	}

	return priv, nil
}

// GetListenAddress ...
func GetListenAddress() libp2p.Option {
	host := viper.GetString(libp2pClientHostViperKey)
	port := viper.GetString(libp2pClientPortViperKey)

	return libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/%v/tcp/%v", 
			host, 
			port,
		),
	)
}

// GetIdentity ...
func GetIdentity() (libp2p.Option, error) {
	priv, err := readPrivateKey()
	if err != nil {
		return nil, err
	}
	return libp2p.Identity(priv), nil
}

