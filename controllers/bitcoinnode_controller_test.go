package controllers

import (
	"github.com/btcsuite/btcd/rpcclient"
	bitcoinv1alpha1 "github.com/kiln-fired/kiln-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestReconcileBitcoinNode_CreateStatefulset(t *testing.T) {

	b := &bitcoinv1alpha1.BitcoinNode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-namespace",
		},
		Spec: bitcoinv1alpha1.BitcoinNodeSpec{
			RPCServer: bitcoinv1alpha1.RPCServer{
				CertSecret: "some-secret",
				User:       "some-user",
				Password:   "some-password",
			},
		},
	}

	r := makeTestReconciler(t, b)

	assert.NotNil(t, r.statefulsetForBitcoinNode(b))
}

func Test_BlockCountCommand(t *testing.T) {
	client := makeBitcoinClient()
	count, err := client.GetBlockCount()

	if err != nil {
		t.Fatalf("failed to run RPC command: %s", err)
	}

	assert.Equal(t, int64(400), count)
}

func Test_HashesPerSec(t *testing.T) {
	client := makeBitcoinClient()
	hashCount, err := client.GetMiningInfo()

	if err != nil {
		t.Fatalf("failed to run RPC command: %s", err)
	}

	assert.Equal(t, "", hashCount)
}

func Test_Generate(t *testing.T) {
	client := makeBitcoinClient()
	err := client.SetGenerate(false, 0)

	if err != nil {
		t.Fatalf("failed to run RPC command: %s", err)
	}

	//assert.Equal(t, int64(0), peers)
}

func Test_NodeConnect(t *testing.T) {
	client := makeBitcoinClient()
	perm := "perm"
	err := client.Node("connect", "friend:18555", &perm)

	if err != nil {
		t.Fatalf("failed to run RPC command: %s", err)
	}

	//assert.Equal(t, int64(0), peers)
}

func Test_SetDebug(t *testing.T) {
	client := makeBitcoinClient()
	generate, err := client.DebugLevel("info")

	if err != nil {
		t.Fatalf("failed to run RPC command: %s", err)
	}

	assert.Equal(t, "Done.", generate)
}

func makeBitcoinClient() *rpcclient.Client {
	certs, _ := ioutil.ReadFile("/home/davgordo/Code/kiln/ca.crt")

	connCfg := &rpcclient.ConnConfig{
		Host:         "btcd-kiln.apps-crc.testing:443",
		User:         "node-user",
		Pass:         "st4cks4ts",
		Certificates: certs,
		HTTPPostMode: true,
	}

	client, _ := rpcclient.New(connCfg, nil)
	return client
}

func makeTestReconciler(t *testing.T, objs ...runtime.Object) *BitcoinNodeReconciler {
	s := scheme.Scheme
	assert.NoError(t, bitcoinv1alpha1.AddToScheme(s))

	cl := fake.NewClientBuilder().WithScheme(s).WithRuntimeObjects(objs...).Build()
	return &BitcoinNodeReconciler{
		Client: cl,
		Scheme: s,
	}
}
