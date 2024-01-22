package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"

	"github.com/pquerna/otp"
)

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// ทำการ copy ไฟล์จาก src ไปยัง out
	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

func GenerateQRCodeURL(key *otp.Key) string {
	issuer := url.QueryEscape(key.Issuer())
	accountName := url.QueryEscape(key.AccountName())
	secret := url.QueryEscape(key.Secret())
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, accountName, secret, issuer)
}

// func DereferenceString(s *string) string { << no need anymore because we want to return null instead or ""
// 	if s != nil {
// 		return *s
// 	}
// 	return ""
// }
