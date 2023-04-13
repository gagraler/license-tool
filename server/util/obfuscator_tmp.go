package util

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"unicode"
)

type Obfuscator interface {
	Obfuscate(input []byte, signatureCode string) (string, error)
}

type XorObfuscator struct {
}

type RandomReplaceObfuscator struct {
	Ratio float64 // 替换比例
}

type ReverseXorObfuscator struct {
}

type Base64XorObfuscator struct {
}

func (o *XorObfuscator) Obfuscate(input []byte, signatureCode string) (string, error) {
	const targetLen = 4096

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	ciphertext := make([]byte, targetLen)
	for i := 0; i < targetLen; i++ {
		if i < len(input) {
			ciphertext[i] = input[i] ^ scBytes[i%len(scBytes)]
		} else {
			randByte := byte(rand.Intn(256))
			for randByte != '0' && randByte != '1' && !unicode.IsLetter(rune(randByte)) {
				randByte = byte(rand.Intn(256))
			}
			ciphertext[i] = randByte
		}
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (o *RandomReplaceObfuscator) Obfuscate(input []byte, signatureCode string) (string, error) {

	const targetLen = 4096

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	ciphertext := make([]byte, targetLen)

	for i := 0; i < targetLen; i++ {
		if i < len(input) {
			if rand.Float64() < o.Ratio {
				randByte := byte(rand.Intn(256))
				for randByte == input[i] || randByte == '0' || randByte == '1' || unicode.IsLetter(rune(randByte)) {
					randByte = byte(rand.Intn(256))
				}
				ciphertext[i] = randByte ^ scBytes[i%len(scBytes)]
			} else {
				ciphertext[i] = input[i] ^ scBytes[i%len(scBytes)]
			}
		} else {
			randByte := byte(rand.Intn(256))
			for randByte != '0' && randByte != '1' && !unicode.IsLetter(rune(randByte)) {
				randByte = byte(rand.Intn(256))
			}
			ciphertext[i] = randByte ^ scBytes[i%len(scBytes)]
		}
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (o *ReverseXorObfuscator) Obfuscate(input []byte, signatureCode string) (string, error) {
	const targetLen = 4096

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	ciphertext := make([]byte, targetLen)

	for i := 0; i < targetLen; i++ {
		if i < len(input) {
			ciphertext[i] = input[len(input)-1-i] ^ scBytes[i%len(scBytes)]
		} else {
			randByte := byte(rand.Intn(256))
			for randByte != '0' && randByte != '1' && !unicode.IsLetter(rune(randByte)) {
				randByte = byte(rand.Intn(256))
			}
			ciphertext[i] = randByte
		}
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (o *Base64XorObfuscator) Obfuscate(input []byte, signatureCode string) (string, error) {
	const targetLen = 4096

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	ciphertext := make([]byte, targetLen)

	for i := 0; i < targetLen; i++ {
		if i < len(input) {
			ciphertext[i] = input[i] ^ scBytes[i%len(scBytes)]
		} else {
			randByte := byte(rand.Intn(256))
			for randByte != '0' && randByte != '1' && !unicode.IsLetter(rune(randByte)) {
				randByte = byte(rand.Intn(256))
			}
			ciphertext[i] = randByte ^ scBytes[i%len(scBytes)]
		}
	}

	encoded := base64.URLEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

//
//func main() {
//	var config Config
//	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	obfuscatorMap := map[string]Obfuscator{
//		"xor":            &XorObfuscator{},
//		"random_replace": &RandomReplaceObfuscator{Ratio: config.RandomReplace.Ratio},
//		"base64_xor":     &Base64XorObfuscator{},
//		"reverse_xor":    &ReverseXorObfuscator{},
//	}
//
//	obfuscationEnabled := config.Obfuscation.Enabled
//	algorithm := strings.ToLower(config.Obfuscation.Algorithm)
//	obfuscator, ok := obfuscatorMap[algorithm]
//	if !ok {
//		fmt.Printf("unsupported obfuscation algorithm: %s\n", algorithm)
//		os.Exit(1)
//	}
//
//	for {
//		var input string
//		_, err := fmt.Scanln(&input)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		signatureCode := strconv.FormatInt(time.Now().UnixNano(), 10)
//		encoded, err := obfuscator.Obfuscate([]byte(input), signatureCode)
//
//		if err != nil {
//			fmt.Println("error while obfuscating:", err.Error())
//			continue
//		}
//
//		if obfuscationEnabled {
//			fmt.Printf("%s\n", encoded)
//		} else {
//			fmt.Printf("%s_%s\n", signatureCode, encoded)
//		}
//	}
//}
