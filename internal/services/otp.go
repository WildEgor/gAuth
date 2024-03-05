package services

import (
	"bytes"
	"errors"
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/pquerna/otp/totp"
	"image/png"
	"math/rand"
	"strconv"
)

var (
	ErrGenerateOTP  = errors.New("failed to generate OTP")
	ErrSendSMS      = errors.New("failed to send SMS")
	ErrorGenerateQR = errors.New("failed to generate QR code")
)

type CodeGenerator struct {
	Length uint8
}

func NewCodeGenerator(length uint8) *CodeGenerator {
	return &CodeGenerator{
		Length: length,
	}
}

func (g *CodeGenerator) Generate() string {
	num := rand.Intn(9000) + 1000
	return strconv.Itoa(num)
}

type OTPService struct {
	config    *configs.OTPConfig
	generator *CodeGenerator
}

func NewOTPService(
	config *configs.OTPConfig,
) *OTPService {
	return &OTPService{
		config:    config,
		generator: NewCodeGenerator(config.Length),
	}
}

func (g *OTPService) GenerateAndSMSSend(phone string) (string, error) {
	code := g.generator.Generate()

	// TODO: send SMS

	return code, nil
}

func (g *OTPService) GenerateQR(identity string) (bytes.Buffer, error) {
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
