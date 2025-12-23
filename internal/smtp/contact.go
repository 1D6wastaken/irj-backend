package smtp

import "context"

func (s *SMTPService) SendContactEmail(ctx context.Context, innerSubject, body, email string, to []EmailPerson) error {
	subject := "Un contributeur vous a contacté : " + innerSubject

	htmlEmail := `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Un contributeur vous a contacté</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
    <tr>
      <td align="center" style="padding:20px;">
        <a href="https://saintjacquesinfo.eu/#/" target="_blank" 
			style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
          Le site du Patrimoine Saint-Jacques
        </a>
      </td>
    </tr>
  </table>

  <!-- Contenu -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#fffcf8; padding:20px;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" 
			style="background-color:#ffffff; border-radius:8px; padding:20px; border:1px solid #e0e0e0;">
          <tr>
            <td style="color:#333333; font-size:16px; line-height:1.5;">
              <p>Bonjour,</p>
			  <p>Vous avez reçu un nouveau message d'un contributeur via le formulaire de contact :</p>
			  <p>Objet du message : ` + innerSubject + `</p>
			  <pre>` + body + `</pre>
			  <p>Vous pouvez répondre directement à cette adresse email : <a href="mailto:` + email + `">` + email + `</a></p>
			</td>
		  </tr><tr>
            <td style="padding-top:30px; font-size:12px; color:#999999; text-align:center;">
              Cet email est généré automatiquement. Merci de ne pas y répondre directement.
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>

</body>
</html>`

	textContent := `Bonjour,

Vous avez reçu un nouveau message d'un contributeur via le formulaire de contact :

Objet du message : ` + innerSubject + `

` + body + `

Vous pouvez répondre directement à cette adresse email : ` + email + `

Cet email est généré automatiquement. Merci de ne pas y répondre directement.`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}
