// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	l4g "github.com/alecthomas/log4go"

	"github.com/cetessai2501/gasilo/model"
)

var IsLicensed bool = false
var License *model.License = &model.License{
	Features: new(model.Features),
}
var ClientLicense map[string]string = map[string]string{"IsLicensed": "false"}

var publicKey []byte = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyZmShlU8Z8HdG0IWSZ8r
tSyzyxrXkJjsFUf0Ke7bm/TLtIggRdqOcUF3XEWqQk5RGD5vuq7Rlg1zZqMEBk8N
EZeRhkxyaZW8pLjxwuBUOnXfJew31+gsTNdKZzRjrvPumKr3EtkleuoxNdoatu4E
HrKmR/4Yi71EqAvkhk7ZjQFuF0osSWJMEEGGCSUYQnTEqUzcZSh1BhVpkIkeu8Kk
1wCtptODixvEujgqVe+SrE3UlZjBmPjC/CL+3cYmufpSNgcEJm2mwsdaXp2OPpfn
a0v85XL6i9ote2P+fLZ3wX9EoioHzgdgB7arOxY50QRJO7OyCqpKFKv6lRWTXuSt
hwIDAQAB
-----END PUBLIC KEY-----`)

func LoadLicense(licenseBytes []byte) {
	if success, licenseStr := ValidateLicense(licenseBytes); success {
		license := model.LicenseFromJson(strings.NewReader(licenseStr))
		SetLicense(license)
		return
	}

	l4g.Warn(T("utils.license.load_license.invalid.warn"))
}

func SetLicense(license *model.License) bool {
	license.Features.SetDefaults()

	if !license.IsExpired() {
		License = license
		IsLicensed = true
		ClientLicense = getClientLicense(license)
		ClientCfg = getClientConfig(Cfg)
		return true
	}

	return false
}

func RemoveLicense() {
	License = &model.License{}
	IsLicensed = false
	ClientLicense = getClientLicense(License)
	ClientCfg = getClientConfig(Cfg)
}

func ValidateLicense(signed []byte) (bool, string) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(signed)))

	_, err := base64.StdEncoding.Decode(decoded, signed)
	if err != nil {
		l4g.Error(T("utils.license.validate_license.decode.error"), err.Error())
		return false, ""
	}

	if len(decoded) <= 256 {
		l4g.Error(T("utils.license.validate_license.not_long.error"))
		return false, ""
	}

	// remove null terminator
	for decoded[len(decoded)-1] == byte(0) {
		decoded = decoded[:len(decoded)-1]
	}

	plaintext := decoded[:len(decoded)-256]
	signature := decoded[len(decoded)-256:]

	block, _ := pem.Decode(publicKey)

	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		l4g.Error(T("utils.license.validate_license.signing.error"), err.Error())
		return false, ""
	}

	rsaPublic := public.(*rsa.PublicKey)

	h := sha512.New()
	h.Write(plaintext)
	d := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(rsaPublic, crypto.SHA512, d, signature)
	if err != nil {
		l4g.Error(T("utils.license.validate_license.invalid.error"), err.Error())
		return false, ""
	}

	return true, string(plaintext)
}

func GetLicenseFileFromDisk(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		l4g.Error("Failed to open license key from disk at %v err=%v", fileName, err.Error())
		return nil
	}
	defer file.Close()

	licenseBytes, err := ioutil.ReadAll(file)
	if err != nil {
		l4g.Error("Failed to read license key from disk at %v err=%v", fileName, err.Error())
		return nil
	}

	return licenseBytes
}

func GetLicenseFileLocation(fileLocation string) string {
	if fileLocation == "" {
		return FindDir("config") + "mattermost.mattermost-license"
	} else {
		return fileLocation
	}
}

func getClientLicense(l *model.License) map[string]string {
	props := make(map[string]string)

	props["IsLicensed"] = strconv.FormatBool(IsLicensed)

	if IsLicensed {
		props["Id"] = l.Id
		props["Users"] = strconv.Itoa(*l.Features.Users)
		props["LDAP"] = strconv.FormatBool(*l.Features.LDAP)
		props["MFA"] = strconv.FormatBool(*l.Features.MFA)
		props["SAML"] = strconv.FormatBool(*l.Features.SAML)
		props["Cluster"] = strconv.FormatBool(*l.Features.Cluster)
		props["Metrics"] = strconv.FormatBool(*l.Features.Metrics)
		props["GoogleOAuth"] = strconv.FormatBool(*l.Features.GoogleOAuth)
		props["Office365OAuth"] = strconv.FormatBool(*l.Features.Office365OAuth)
		props["Compliance"] = strconv.FormatBool(*l.Features.Compliance)
		props["CustomBrand"] = strconv.FormatBool(*l.Features.CustomBrand)
		props["MHPNS"] = strconv.FormatBool(*l.Features.MHPNS)
		props["PasswordRequirements"] = strconv.FormatBool(*l.Features.PasswordRequirements)
		props["IssuedAt"] = strconv.FormatInt(l.IssuedAt, 10)
		props["StartsAt"] = strconv.FormatInt(l.StartsAt, 10)
		props["ExpiresAt"] = strconv.FormatInt(l.ExpiresAt, 10)
		props["Name"] = l.Customer.Name
		props["Email"] = l.Customer.Email
		props["Company"] = l.Customer.Company
		props["PhoneNumber"] = l.Customer.PhoneNumber
	}

	return props
}

func GetClientLicenseEtag(useSanitized bool) string {
	value := ""

	lic := ClientLicense

	if useSanitized {
		lic = GetSanitizedClientLicense()
	}

	for k, v := range lic {
		value += fmt.Sprintf("%s:%s;", k, v)
	}

	return model.Etag(fmt.Sprintf("%x", md5.Sum([]byte(value))))
}

func GetSanitizedClientLicense() map[string]string {
	sanitizedLicense := make(map[string]string)

	for k, v := range ClientLicense {
		sanitizedLicense[k] = v
	}

	if IsLicensed {
		delete(sanitizedLicense, "Id")
		delete(sanitizedLicense, "Name")
		delete(sanitizedLicense, "Email")
		delete(sanitizedLicense, "PhoneNumber")
		delete(sanitizedLicense, "IssuedAt")
		delete(sanitizedLicense, "StartsAt")
		delete(sanitizedLicense, "ExpiresAt")
	}

	return sanitizedLicense
}
