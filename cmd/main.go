package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	global "github.com/Subskribo-BV/dnn-fabric-api"
	"github.com/Subskribo-BV/dnn-fabric-api/api/handler"
	"github.com/Subskribo-BV/dnn-fabric-api/api/router"
	"github.com/Subskribo-BV/dnn-fabric-api/service"
	"github.com/Subskribo-BV/dnn-fabric-api/utils/auth"

	"github.com/gorilla/handlers"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/urfave/negroni"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	log.Println("initializing the service...")
	conf := global.LoadConfig()

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(conf.TlsCertPath, conf.GatewayPeer, conf.PeerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(conf.CertPath, conf.MspId)
	sign := newSign(conf.KeyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(conf.ChannelName)
	contract := network.GetContract(conf.ChaincodeName)

	s := service.New(contract)
	h := handler.New(s)
	r := router.BuildRouter(h)

	auth.Init(conf.Server.JwtKey)

	go func() {
		n := negroni.Classic()
		n.UseHandler(r)

		server := http.Server{
			Addr: fmt.Sprintf(":%d", conf.Server.Port),
			Handler: handlers.CORS(
				handlers.ExposedHeaders(conf.Server.ExposedHeaders),
				handlers.AllowedHeaders(conf.Server.AllowedHeaders),
				handlers.AllowedMethods(conf.Server.AllowedMethods),
				handlers.AllowedOrigins(conf.Server.AllowedOrigins),
			)(n),
		}

		log.Println("service listening at", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	lock := make(chan os.Signal, 1)
	signal.Notify(lock, os.Interrupt, syscall.SIGTERM)
	<-lock

	log.Printf("terminating the service in %d seconds\n", 5)
	time.Sleep(5 * time.Second)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint string) *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(certPath, mspID string) *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign(keyPath string) identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}

	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
