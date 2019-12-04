package chiperutil

import (
	"fmt"
	"testing"
)

func Test_Aes(t *testing.T) {
	//str := GenerateRandSeq(32)
	str := "BpLnfgDsc2WD8F2qNfHK5a84jjJkwzDk"
	fmt.Println(str)
	encryptStr := AesEncryptToStr("3306", str)
	fmt.Println(encryptStr)
	decryptStr := AesDecryptToStr(encryptStr, str)
	fmt.Println(decryptStr)
}
