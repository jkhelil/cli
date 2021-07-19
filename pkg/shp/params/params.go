package params

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	buildclientset "github.com/shipwright-io/build/pkg/client/clientset/versioned"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// Params is a place for Shipwright CLI to store its runtime parameters including configured dynamic
// client and global flags.
type Params struct {
	client       dynamic.Interface
	clientset    kubernetes.Interface
	shpClientset buildclientset.Interface

	configFlags *genericclioptions.ConfigFlags
	namespace   string
}

// AddFlags accepts flags and adds program global flags to it
func (p *Params) AddFlags(flags *pflag.FlagSet) {
	p.configFlags.AddFlags(flags)
}

// Client returns preconfigured dynamic client with overrides
// from global flags and kubernetes configuration set by user
func (p *Params) Client() (dynamic.Interface, error) {
	if p.client != nil {
		return p.client, nil
	}

	clientConfig := p.configFlags.ToRawKubeConfigLoader()

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	p.namespace, _, err = clientConfig.Namespace()
	if err != nil {
		return nil, err
	}

	p.client, err = dynamic.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "could not create Dynamic client")
	}

	return p.client, nil
}

// ClientSet returns a kubernetes clientset.
func (p *Params) ClientSet() (kubernetes.Interface, error) {
	if p.clientset != nil {
		return p.clientset, nil
	}

	clientConfig := p.configFlags.ToRawKubeConfigLoader()
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	p.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return p.clientset, nil
}

// ShipwrightClientSet returns a Shipwright Clientset
func (p *Params) ShipwrightClientSet() (buildclientset.Interface, error) {
	if p.shpClientset != nil {
		return p.shpClientset, nil
	}
	clientConfig := p.configFlags.ToRawKubeConfigLoader()
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	p.shpClientset, err = buildclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return p.shpClientset, nil
}

// Namespace returns kubernetes namespace with alle the overrides
// from command line and kubernetes config
func (p *Params) Namespace() string {
	return p.namespace
}

// NewParams creates a new instance of ShipwrightParams and returns it as
// an interface value
func NewParams() *Params {
	p := &Params{}
	p.configFlags = genericclioptions.NewConfigFlags(true)

	return p
}
