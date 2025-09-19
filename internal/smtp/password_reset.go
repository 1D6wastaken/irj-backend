package smtp

import (
	"context"
	"fmt"
)

//nolint:lll
func (s *SMTPService) SendPasswordResetMail(ctx context.Context, to []EmailPerson, token string) error {
	subject := "Réinitialisation de votre mot de passe"

	htmlEmail := fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Réinitialisation de votre mot de passe</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
    <tr>
      <td align="center" style="padding:20px;">
        <a href="https://test.saintjacquesinfo.eu/#/" target="_blank" style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
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
              <p>Vous avez demandé à réinitialiser votre mot de passe pour votre compte sur <em>Le site du Patrimoine Saint-Jacques</em>.</p>
              <p>Pour définir un nouveau mot de passe, veuillez cliquer sur le bouton ci-dessous&nbsp;:</p>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; text-align:center;">
              <a href="https://test.saintjacquesinfo.eu/reset/%s" 
                 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
                Réinitialiser mon mot de passe
              </a>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; font-size:14px; color:#666666;">
              Si le bouton ne fonctionne pas, copiez et collez ce lien dans votre navigateur&nbsp;:<br>
              <a href="https://test.saintjacquesinfo.eu/reset/%s" style="color:#802331;">https://test.saintjacquesinfo.eu/reset/%s</a>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; font-size:14px; color:#666666;">
              Pour des raisons de sécurité, ce lien expirera dans 30 minutes.
            </td>
          </tr>
          <tr>
            <td style="padding-top:30px; font-size:12px; color:#999999; text-align:center;">
              Si vous n'êtes pas à l'origine de cette demande, vous pouvez ignorer cet email.
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>

</body>
</html>`, token, token, token)

	textContent := fmt.Sprintf(`Bonjour,

Vous avez demandé à réinitialiser votre mot de passe pour votre compte sur "Le site du Patrimoine Saint-Jacques".

Pour définir un nouveau mot de passe, cliquez sur ce lien :
https://test.saintjacquesinfo.eu/reset/%s

Si le lien ne fonctionne pas, copiez-le et collez-le dans votre navigateur.

Ce lien expirera dans 30 minutes.

Si vous n'êtes pas à l'origine de cette demande, ignorez cet email.

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`, token)

	return s.send(ctx, to, subject, htmlEmail, textContent)
}
