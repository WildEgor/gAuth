package services

import (
	"bytes"
	"errors"
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/pquerna/otp/totp"
	"image/png"
	"time"
)

var (
	ErrGenerateOTP  = errors.New("failed to generate OTP")
	ErrSendSMS      = errors.New("failed to send SMS")
	ErrorGenerateQR = errors.New("failed to generate QR code")
)

type OTPGenerator struct {
	config *configs.OTPConfig
}

func NewOTPGenerator(
	config *configs.OTPConfig,

) *OTPGenerator {
	return &OTPGenerator{
		config: config,
	}
}

func (g *OTPGenerator) GenerateAndSMSSend(phone string) (string, error) {
	code, err := totp.GenerateCode(g.config.Secret, time.Now())
	if err != nil {
		return "", ErrGenerateOTP
	}

	// TODO: impl
	//err = g.smsSender.Send(phone, code)
	//if err != nil {
	//	return "", ErrSendSMS
	//}

	return code, nil
}

func (g *OTPGenerator) GenerateAndEmailSend(email string) (string, error) {
	code, err := totp.GenerateCode(g.config.Secret, time.Now())
	if err != nil {
		return "", ErrGenerateOTP
	}

	// TODO: impl
	//err = g.emailSender.Send(email, code)
	//if err != nil {
	//	return "", ErrSendSMS
	//}

	return code, nil
}

func (g *OTPGenerator) GenerateQR(identity string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      g.config.Issuer,
		AccountName: identity,
	})
	if err != nil {
		return buf, ErrGenerateOTP
	}

	img, err := key.Image(200, 200)
	if err != nil {
		return buf, ErrorGenerateQR
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return buf, ErrorGenerateQR
	}

	return buf, nil
}
