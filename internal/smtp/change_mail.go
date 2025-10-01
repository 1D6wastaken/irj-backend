package smtp

import (
	"context"
	"fmt"
)

//nolint:lll
func (s *SMTPService) SendEmailConfirmationMail(ctx context.Context, to []EmailPerson, id string) error {
	subject := "Confirmation de votre nouvelle adresse email"

	htmlEmail := fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Confirmation de votre nouvelle adresse email</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
    <tr>
      <td align="center" style="padding:20px;">
        <a href="https://saintjacquesinfo.eu/#/" target="_blank" style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
          Le site du Patrimoine Saint-Jacques
        </a>
      </td>
    </tr>
  </table>

  <!-- Contenu -->
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#fffcf8; padding:20px;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff; border-radius:8px; padding:20px; border:1px solid #e0e0e0;">
          <tr>
            <td style="color:#333333; font-size:16px; line-height:1.5;">
              <p>Bonjour,</p>
              <p>Vous avez récemment modifié l'adresse email associée à votre compte sur <em>Le site du Patrimoine Saint-Jacques</em>.</p>
              <p>Pour valider cette modification, merci de confirmer votre nouvelle adresse en cliquant sur le bouton ci-dessous :</p>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; text-align:center;">
              <a href="https://saintjacquesinfo.eu/email/%s/validate" 
                 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
                Confirmer mon adresse email
              </a>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; font-size:14px; color:#666666;">
              Si le bouton ne fonctionne pas, copiez et collez ce lien dans votre navigateur :
              <br>
              <a href="https://saintjacquesinfo.eu/email/%s/validate" style="color:#802331;">https://saintjacquesinfo.eu/email/%s/validate</a>
            </td>
          </tr>
          <tr>
            <td style="padding-top:30px; font-size:12px; color:#999999; text-align:center;">
              Cet email est généré automatiquement. Merci de ne pas y répondre directement.
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>

</body>
</html>`, id, id, id)

	textContent := fmt.Sprintf(`Bonjour,

Vous avez récemment modifié l'adresse email associée à votre compte sur "Le site du Patrimoine Saint-Jacques".

Pour valider cette modification, merci de confirmer votre nouvelle adresse en cliquant sur ce lien :
https://saintjacquesinfo.eu/email/%s/validate

Si vous n'êtes pas à l'origine de cette modification, veuillez nous contacter immédiatement.

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`, id)

	return s.send(ctx, to, subject, htmlEmail, textContent)
}
