package oauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	OAUTH_VERSION              = "1.0"
	OAUTH_SIGNATURE_METHOD     = "HMAC-SHA1"
	OAUTH_CONSUMER_KEY_KEY     = "oauth_consumer_key"
	OAUTH_TOKEN_KEY            = "oauth_token"
	OAUTH_SIGNATURE_KEY        = "oauth_signature"
	OAUTH_SIGNATURE_METHOD_KEY = "oauth_signature_method"
	OAUTH_TIMESTAMP_KEY        = "oauth_timestamp"
	OAUTH_NONCE_KEY            = "oauth_nonce"
	OAUTH_VERSION_KEY          = "oauth_version"
	OAUTH_CALLBACK_KEY         = "oauth_callback"
)

type OAuth struct {
	OPS               OAuthParameters
	ParameterString   string
	OauthHeaderValues map[string]string
	RequestApiUrl     string
	RequestMethod     string
	ConsumerSecretKey string
	OauthToken        string
	OauthSecretKey    string
}

type OAuthParameters struct {
	KeyValues        map[string]string
	EncodedKeyValues map[string]string
}

func (oap *OAuthParameters) init() {
	oap.KeyValues = make(map[string]string, 20)
	oap.KeyValues[OAUTH_CONSUMER_KEY_KEY] = ""
	oap.KeyValues[OAUTH_TOKEN_KEY] = ""
	oap.KeyValues[OAUTH_CALLBACK_KEY] = ""
	oap.KeyValues[OAUTH_NONCE_KEY] = ""
	oap.KeyValues[OAUTH_SIGNATURE_KEY] = ""
	oap.KeyValues[OAUTH_SIGNATURE_METHOD_KEY] = ""
	oap.KeyValues[OAUTH_TIMESTAMP_KEY] = ""
	oap.KeyValues[OAUTH_VERSION_KEY] = ""

	oap.EncodedKeyValues = make(map[string]string, 20)
}

func (oap *OAuthParameters) clearKeyValues() {
	for k, _ := range oap.KeyValues {
		if k == OAUTH_CALLBACK_KEY ||
			k == OAUTH_CONSUMER_KEY_KEY ||
			k == OAUTH_NONCE_KEY ||
			k == OAUTH_SIGNATURE_METHOD_KEY ||
			k == OAUTH_TIMESTAMP_KEY ||
			k == OAUTH_TOKEN_KEY ||
			k == OAUTH_VERSION_KEY {
			continue
		}

		oap.KeyValues[k] = ""
	}
}

func (oap *OAuthParameters) clearEncodedKeyValues() {
	oap.EncodedKeyValues = make(map[string]string, 20)
}

func (oa *OAuth) Init() {
	oa.OPS.init()
}

func (oa *OAuth) Clear() {
	oa.OPS.clearEncodedKeyValues()
	oa.OPS.clearKeyValues()
}

func (oa *OAuth) clearOauthHeaderValues() {
	oa.OauthHeaderValues = make(map[string]string, 20)
}

func (oa *OAuth) SetParameter(key string, value string) {
	oa.OPS.KeyValues[key] = value
}

func (oa *OAuth) SetApiBaseUrl(baseUrl string) {
	oa.RequestApiUrl = baseUrl
}

func (oa *OAuth) GetApiBaseUrl() string {
	return oa.RequestApiUrl
}

func (oa *OAuth) SetRequestMethod(requestMethod string) {
	oa.RequestMethod = strings.ToUpper(requestMethod)
}

func (oa *OAuth) SetConsumerSecretKey(consumerSecretKey string) {
	oa.ConsumerSecretKey = consumerSecretKey
}

func (oa *OAuth) SetConsumerKey(consumerKey string) {
	oa.OPS.KeyValues[OAUTH_CONSUMER_KEY_KEY] = consumerKey
}

func (oa *OAuth) SetOauthToken(oauthToken string) {
	oa.OauthToken = oauthToken
}

func (oa *OAuth) SetOauthTokenSecret(oauthTokenSecret string) {
	oa.OauthSecretKey = oauthTokenSecret
}

func (oa *OAuth) SetCallback(callback string) {
	oa.OPS.KeyValues[OAUTH_CALLBACK_KEY] = callback
}

func (oa *OAuth) GetCallback() string {
	return oa.OPS.KeyValues[OAUTH_CALLBACK_KEY]
}

func (oa *OAuth) GenerateNonce() string {
	nonce := sha256.Sum256([]byte(oa.GenerateTimestamp()))
	replacer := strings.NewReplacer("_", "", "+", "", "%", "", "=", "", "/", "")
	return replacer.Replace(base64.StdEncoding.EncodeToString(nonce[:32]))
	//return "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg"
}

func (oa *OAuth) GenerateTimestamp() string {
	return fmt.Sprintf("%10d", time.Now().Unix())
	//return "1318622958"
}

func (oa *OAuth) isOauthHeaderKey(key string) bool {
	switch key {
	case OAUTH_CALLBACK_KEY:
		return true
	case OAUTH_CONSUMER_KEY_KEY:
		return true
	case OAUTH_NONCE_KEY:
		return true
	case OAUTH_SIGNATURE_KEY:
		return true
	case OAUTH_SIGNATURE_METHOD_KEY:
		return true
	case OAUTH_TIMESTAMP_KEY:
		return true
	case OAUTH_TOKEN_KEY:
		return true
	case OAUTH_VERSION_KEY:
		return true
	default:
		return false
	}
}

func (oa *OAuth) encode(value string) string {
	return strings.ReplaceAll(url.QueryEscape(value), "+", "%20")
}

func (oa *OAuth) prepareOauthParameters(requiredOauthToken bool) (err error) {
	if oa.OPS.KeyValues[OAUTH_CONSUMER_KEY_KEY] == "" {
		oa.OPS.KeyValues[OAUTH_CONSUMER_KEY_KEY] = os.Getenv("OAUTH_CONSUMER_KEY")
	}

	if oa.OPS.KeyValues[OAUTH_TOKEN_KEY] == "" {
		oa.OPS.KeyValues[OAUTH_TOKEN_KEY] = oa.OauthToken
	}

	oa.OPS.KeyValues[OAUTH_NONCE_KEY] = oa.GenerateNonce()
	oa.OPS.KeyValues[OAUTH_TIMESTAMP_KEY] = oa.GenerateTimestamp()
	oa.OPS.KeyValues[OAUTH_SIGNATURE_METHOD_KEY] = OAUTH_SIGNATURE_METHOD
	oa.OPS.KeyValues[OAUTH_VERSION_KEY] = OAUTH_VERSION

	if oa.OPS.KeyValues[OAUTH_CONSUMER_KEY_KEY] == "" {
		return errors.New(OAUTH_CONSUMER_KEY_KEY + " is empty.")
	}

	if requiredOauthToken == true {
		if oa.OPS.KeyValues[OAUTH_TOKEN_KEY] == "" {
			return errors.New(OAUTH_TOKEN_KEY + " is empty.")
		}
	}

	//encode key - value
	oa.OPS.clearEncodedKeyValues()
	oa.clearOauthHeaderValues()
	var keys []string = make([]string, len(oa.OPS.KeyValues))
	for k, v := range oa.OPS.KeyValues {
		if k != "" && v != "" {
			ek := oa.encode(k)
			ev := oa.encode(v)

			oa.OPS.EncodedKeyValues[ek] = ev
			keys = append(keys, ek)

			if oa.isOauthHeaderKey(k) == true {
				oa.OauthHeaderValues[ek] = ev
			}
		}
	}

	sort.Strings(keys)

	oa.ParameterString = ""
	var ParameterStringSlice []string
	for _, k := range keys {
		if k != "" && oa.OPS.EncodedKeyValues[k] != "" {
			keyValue := fmt.Sprintf("%s=%s", k, oa.OPS.EncodedKeyValues[k])
			ParameterStringSlice = append(ParameterStringSlice, keyValue)
		}
	}

	oa.ParameterString = strings.Join(ParameterStringSlice, "&")

	encodedRequestMethod := oa.encode(oa.RequestMethod)
	encodedRequestUrl := oa.encode(oa.RequestApiUrl)

	oa.ParameterString = fmt.Sprintf("%s&%s&%s", encodedRequestMethod, encodedRequestUrl, oa.encode(oa.ParameterString))

	//os.Create("./signaturebasestr.txt")
	//ioutil.WriteFile("./signaturebasestr.txt", []byte(oa.ParameterString), 0777)

	if oa.ConsumerSecretKey == "" {
		return errors.New("CONSUMER SECRET KEY is empty.")
	}

	signKey := fmt.Sprintf("%s&%s", oa.ConsumerSecretKey, oa.OauthSecretKey)
	signature := oa.calculateSignature(oa.ParameterString, signKey)

	//os.Create("./signkey.txt")
	//ioutil.WriteFile("./signkey.txt", []byte(signKey+"\n"+signature+"\n"+url.QueryEscape(signature)), 0777)

	oa.OPS.KeyValues[OAUTH_SIGNATURE_KEY] = signature
	oa.OauthHeaderValues[oa.encode(OAUTH_SIGNATURE_KEY)] = oa.encode(signature)

	return nil
}

func (oa *OAuth) calculateSignature(parameterString string, signKey string) string {
	signKeyByte := []byte(signKey)
	h := hmac.New(sha1.New, signKeyByte)
	h.Write([]byte(parameterString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (oa *OAuth) buildOauthHeaderValue(keyValues *map[string]string) string {
	var parameterStringSlice []string
	for k, v := range *keyValues {
		if k != "" && v != "" {
			parameterStringSlice = append(parameterStringSlice, fmt.Sprintf("%s=\"%s\"", k, v))
		}
	}
	return strings.Join(parameterStringSlice, ", ")
}

func (oa *OAuth) SetAuthorizationHeader(req *http.Request) (err error) {
	err = oa.prepareOauthParameters(false)
	if err != nil {
		return err
	}
	headerValue := oa.buildOauthHeaderValue(&oa.OauthHeaderValues)
	//headerValue = `oauth_nonce="K7ny27JTpKVsTgdyLdDfmQQWVLERj2zAK5BslRsqyw", oauth_callback="http%3A%2F%2Fmyapp.com%3A3005%2Ftwitter%2Fprocess_callback", oauth_signature_method="HMAC-SHA1", oauth_timestamp="1300228849", oauth_consumer_key="OqEqJeafRSF11jBMStrZz", oauth_signature="Pc%2BMLdv028fxCErFyi8KXFM%2BddU%3D", oauth_version="1.0"`
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "OAuth", headerValue))
	//os.Create("./header.txt")
	//ioutil.WriteFile("./header.txt", []byte(headerValue), 0777)

	return nil
}
