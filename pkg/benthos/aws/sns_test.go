package aws

import (
	"regexp"
	"testing"
)

func TestVerifyPayload(t *testing.T) {

	m := map[string]any{
		"Message":          "You have chosen to subscribe to the topic arn:aws:sns:eu-central-1:700242105991:foo.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
		"MessageId":        "5af996da-b9ab-4652-8375-72db5dc16700",
		"Signature":        "NYKFYRGHtyHcgBiB1DZPBkwt0RcP3bbBKYIC7bMkpIsY33o45UH1ijgrx0CeORs+arakrFL+jSKYpLtZAkiwKDVnPJ4Czywx8bTz0V1z1dwMaLfYYjUbqOrRiPxVMUeIbjGUATJwgyH8IiADj04v5OnbNpEA4nIPlELANSlf8SOkxSBx9mNmA3HThgDXpoZTDup1J9wEgdtXZQ5xrrWBuYRxeDrheYOxZzokrc7TK1RSWhy3svPYXwY2vvq0mfbDeoDe9UWLO8/0GDOZscXQ+Irds1E6yvrWoLCJ/2ktruFEtPrQvLvg/mrbDMtuNNC1tOe8Xc22gA3juvEwrSOLoQ==",
		"SignatureVersion": "1",
		"SigningCertURL":   "https://raw.githubusercontent.com/otto-de/sherlock-microservice/main/pkg/aws/testdata/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
		"SubscribeURL":     "https://sns.eu-central-1.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:eu-central-1:700242105991:foo&Token=2336412f37fb687f5d51e6e2425ba1f2583191300ec0352647350ee975394ba6e35f08620711b5121e169c171eaeb11894b7f2afaca0a976aa44978c93fcd61a73f37b2b5dec4e2b0952113dd9a09c80fd7e8a70e76ed888e2bee3a87c7f5cd8e223f1ee33ef8c8c2e4e6f0e4fa77fca",
		"Timestamp":        "2024-09-18T08:55:16.190Z",
		"Token":            "2336412f37fb687f5d51e6e2425ba1f2583191300ec0352647350ee975394ba6e35f08620711b5121e169c171eaeb11894b7f2afaca0a976aa44978c93fcd61a73f37b2b5dec4e2b0952113dd9a09c80fd7e8a70e76ed888e2bee3a87c7f5cd8e223f1ee33ef8c8c2e4e6f0e4fa77fca",
		"TopicArn":         "arn:aws:sns:eu-central-1:700242105991:foo",
		"Type":             "SubscriptionConfirmation",
	}

	vp := snsPayloadVerifier{
		hostPattern: regexp.MustCompile(".*"),
	}
	err := vp.verifyFromURL(
		m,
		"NYKFYRGHtyHcgBiB1DZPBkwt0RcP3bbBKYIC7bMkpIsY33o45UH1ijgrx0CeORs+arakrFL+jSKYpLtZAkiwKDVnPJ4Czywx8bTz0V1z1dwMaLfYYjUbqOrRiPxVMUeIbjGUATJwgyH8IiADj04v5OnbNpEA4nIPlELANSlf8SOkxSBx9mNmA3HThgDXpoZTDup1J9wEgdtXZQ5xrrWBuYRxeDrheYOxZzokrc7TK1RSWhy3svPYXwY2vvq0mfbDeoDe9UWLO8/0GDOZscXQ+Irds1E6yvrWoLCJ/2ktruFEtPrQvLvg/mrbDMtuNNC1tOe8Xc22gA3juvEwrSOLoQ==",
		"https://raw.githubusercontent.com/otto-de/sherlock-microservice/main/pkg/aws/testdata/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
		"1",
	)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}
}
