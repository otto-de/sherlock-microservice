package aws

import (
	"bytes"
	"context"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	_ "github.com/redpanda-data/benthos/v4/public/components/pure"

	"github.com/redpanda-data/benthos/v4/public/service"
)

var (
	_ service.Processor = (*processor)(nil)
)

func init() {
	spec := service.NewConfigSpec().
		Beta().
		Summary("Verifies a SNS message signature").
		Description(`
Verifies a message sent by AWS SNS. This processor implements https://docs.aws.amazon.com/sns/latest/dg/sns-verify-signature-of-message.html.s
`).
		Fields(
			service.NewStringField("host_pattern").
				Advanced().
				Description("field allows for overriding the host check. Intended for testing purposes only.").
				Default(awsHostPattern),
		)
	err := service.RegisterProcessor("aws_sns_message_verify", spec, ctor)
	if err != nil {
		panic(err)
	}
}

func ctor(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {

	hostPattern, err := conf.FieldString("host_pattern")
	if err != nil {
		return nil, err
	}

	if hostPattern == "" {
		hostPattern = awsHostPattern
	}

	return &processor{
		hostPattern: regexp.MustCompile(hostPattern),
	}, nil
}

type processor struct {
	hostPattern *regexp.Regexp
}

func (p *processor) Process(ctx context.Context, message *service.Message) (service.MessageBatch, error) {

	rootA, err := message.AsStructured()
	if err != nil {
		return nil, fmt.Errorf("extracting structure of message failed: %w", err)
	}

	if rootA == nil {
		return nil, errors.New("unexpected nil structure")
	}

	root, ok := rootA.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected a map, got: %T", rootA)
	}

	// root cannot be nil here

	signingCertURLA, ok := root["SigningCertURL"]
	if !ok {
		return nil, errors.New("expected SigningCertURL field")
	}

	if signingCertURLA == nil {
		return nil, errors.New("unexpected nil SigningCertURL")
	}

	signingCertURL, ok := signingCertURLA.(string)
	if !ok {
		return nil, fmt.Errorf("expected SigningCertURL as map, got: %T", signingCertURLA)
	}

	signatureVersionA, ok := root["SignatureVersion"]
	if !ok {
		return nil, errors.New("expected SignatureVersion field")
	}

	if signatureVersionA == nil {
		return nil, errors.New("unexpected nil SignatureVersion")
	}

	signatureVersion, ok := signatureVersionA.(string)
	if !ok {
		return nil, fmt.Errorf("expected SignatureVersion as string, got: %T", signatureVersionA)
	}

	signatureA, ok := root["Signature"]
	if !ok {
		return nil, errors.New("expected Signature field")
	}

	if signatureA == nil {
		return nil, errors.New("unexpected nil Signature")
	}

	signatureBase64, ok := signatureA.(string)
	if !ok {
		return nil, fmt.Errorf("expected Signature as string, got: %T", signatureA)
	}

	vp := snsPayloadVerifier{
		hostPattern: p.hostPattern,
	}
	err = vp.verifyFromURL(root, signatureBase64, signingCertURL, signatureVersion)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return service.MessageBatch{
		message,
	}, nil
}

func (p *processor) Close(ctx context.Context) error {
	return nil
}

const awsHostPattern = `^sns\.[a-zA-Z0-9\-]{3,}\.amazonaws\.com(\.cn)?$`

// buildSignature returns a byte array containing a signature usable for SNS verification
func buildSignature(root map[string]any) []byte {
	var builtSignature bytes.Buffer
	signableKeys := []string{"Message", "MessageId", "Subject", "SubscribeURL", "Timestamp", "Token", "TopicArn", "Type"}
	for _, key := range signableKeys {
		v, ok := root[key]
		if !ok {
			continue
		}
		value := v.(string)
		if value == "" {
			continue
		}
		builtSignature.WriteString(key + "\n")
		builtSignature.WriteString(value + "\n")
	}
	return builtSignature.Bytes()
}

// signatureAlgorithm returns properly Algorithm for AWS Signature Version.
func signatureAlgorithm(signatureVersion string) x509.SignatureAlgorithm {
	if signatureVersion == "2" {
		return x509.SHA256WithRSA
	}
	return x509.SHA1WithRSA
}

type snsPayloadVerifier struct {
	scheme      string
	hostPattern *regexp.Regexp
}

// verifyFromURL will verify that a payload came from SNS
func (vp *snsPayloadVerifier) verifyFromURL(root map[string]any, signatureBase64 string, signingCertURL string, signatureVersion string) error {

	scheme := "https"
	if vp.scheme != "" {
		scheme = vp.scheme
	}

	payloadSignature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return err
	}

	certURL, err := url.Parse(signingCertURL)
	if err != nil {
		return err
	}

	if certURL.Scheme != scheme {
		return fmt.Errorf("url should be using https")
	}

	if !vp.hostPattern.Match([]byte(certURL.Host)) {
		return fmt.Errorf("certificate is located on an invalid domain")
	}

	resp, err := http.Get(signingCertURL)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cert request failed with status: %s", resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return verify(root, body, signatureVersion, payloadSignature)
}

func verify(root map[string]any, pemBody []byte, signatureVersion string, payloadSignature []byte) error {
	decodedPem, _ := pem.Decode(pemBody)
	if decodedPem == nil {
		return errors.New("the decoded PEM file was empty")
	}

	parsedCertificate, err := x509.ParseCertificate(decodedPem.Bytes)
	if err != nil {
		return err
	}

	return parsedCertificate.CheckSignature(signatureAlgorithm(signatureVersion), buildSignature(root), payloadSignature)
}
