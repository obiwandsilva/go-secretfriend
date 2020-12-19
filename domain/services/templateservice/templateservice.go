package templateservice

import (
	"bytes"
	"fmt"
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
	"html/template"
)

type pageDataSecretKey struct {
	SecretKey string
	PhoneNumber string
	AlreadyExists bool
}

type pageDataSecretFriend struct {
	PhoneNumber  string
	PickedFriend entities.Friend
}

const tmplSecretFriend = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Amigo Nóia</title>
</head>
<body>
    <h1>Amigo Nóia<strong></h1>
	<h2>
  		E ae, nóia! Seu/sua nóia secretx se chama {{.PickedFriend.Name}} ({{.PickedFriend.PhoneNumber}}).
  		Guarde muito bem e com carinho o nome. Bjsss
	</h2>
</body>
</html>`

const tmplSecretKey = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Amigo Nóia</title>
</head>
<body>
    <h1>Amigo Nóia<strong></h1>
	<h2>
  		{{if .AlreadyExists}}
			Senha já gerada para o número {{.PhoneNumber}} 
		{{else}}
			E ae, nóia! Sua senha é {{.SecretKey}}
  			Guarde muito bem. Bjsss

			<button onclick="window.location.href='https://qxveiuwevb.execute-api.us-east-1.amazonaws.com/dev/{{.PhoneNumber}}/{{.SecretKey}}/amigo';">
      			Ver meu amigx
    		</button>
		{{end}}
	</h2>
</body>
</html>`

func RenderSecretKey(phoneNumber, secretKey string, alreadyExists bool) (string, error) {
	// Make and parse the HTML template
	t, err := template.New("webpage").Parse(tmplSecretKey)
	if err != nil {
		return "", fmt.Errorf("error when parsing secret key template: %w", err)
	}

	data := pageDataSecretKey{
		SecretKey: secretKey,
		PhoneNumber: phoneNumber,
		AlreadyExists: alreadyExists,
	}

	var buffer bytes.Buffer

	// Render the data and output
	err = t.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("error when rendering to buffer: %w", err)
	}

	return buffer.String(), nil
}

func RenderSecretFriend(phoneNumber string, pickedFriend entities.Friend) (string, error) {
	// Make and parse the HTML template
	t, err := template.New("webpage").Parse(tmplSecretFriend)
	if err != nil {
		return "", fmt.Errorf("error when parsing secret friend template: %w", err)
	}

	data := pageDataSecretFriend{
		PhoneNumber:  phoneNumber,
		PickedFriend: pickedFriend,
	}

	var buffer bytes.Buffer

	// Render the data and output
	err = t.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("error when rendering to buffer: %w", err)
	}

	return buffer.String(), nil
}
