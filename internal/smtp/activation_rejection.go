package smtp

import (
	"context"
	"fmt"
)

//nolint:lll
func (s *SMTPService) SendActivationMail(ctx context.Context, to []EmailPerson, id string) error {
	subject := "Votre candidature a été acceptée"

	htmlEmail := fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Candidature acceptée - Confirmation d'email</title>
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
              <p>Bonne nouvelle 🎉 ! Votre candidature pour devenir <strong>contributeur</strong> sur <em>Le site du Patrimoine Saint-Jacques</em> a été acceptée.</p>
              <p>Il ne vous reste plus qu'à confirmer votre adresse email en cliquant sur le bouton ci-dessous :</p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding-top:20px;">
              <a href="https://test.saintjacquesinfo.eu/email/%s/validate"
                 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
                Confirmer mon email
              </a>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; font-size:14px; color:#555555;">
              Si le bouton ne fonctionne pas, vous pouvez copier-coller ce lien dans votre navigateur :
              <br><a href="https://test.saintjacquesinfo.eu/email/%s/validate" style="color:#802331;">https://test.saintjacquesinfo.eu/email/%s/validate</a>
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

Bonne nouvelle 🎉 ! Votre candidature pour devenir contributeur sur "Le site du Patrimoine Saint-Jacques" a été acceptée.

Il ne vous reste plus qu'à confirmer votre adresse email en cliquant sur le lien suivant :
https://test.saintjacquesinfo.eu/email/%s/validate

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`, id)

	return s.send(ctx, to, subject, htmlEmail, textContent)
}

//nolint:lll
func (s *SMTPService) SendRejectionMail(ctx context.Context, to []EmailPerson) error {
	subject := "Candidature rejetée"

	htmlEmail := `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Candidature rejetée - Le site du Patrimoine Saint-Jacques</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
    <tr>
      <td align="center" style="padding:20px;">
        <a href="https://test.saintjacquesinfo.eu/#/" target="_blank" style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
          Le site du Patrimoine Saint-Jacques
        </a>
      </td>
    </tr>
  </table>

  <!-- Contenu -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#fffcf8; padding:20px;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff; border-radius:8px; padding:20px; border:1px solid #e0e0e0;">
          <tr>
            <td style="color:#333333; font-size:16px; line-height:1.5;">
              <p>Bonjour,</p>
              <p>Après étude de votre demande, nous sommes au regret de vous informer que votre candidature pour devenir <strong>contributeur</strong> sur <em>Le site du Patrimoine Saint-Jacques</em> a été <strong>rejetée</strong>.</p>
              <p>Votre candidature a été supprimée de notre système.</p>
              <p>Nous vous remercions de l’intérêt que vous portez à notre projet et restons à votre disposition pour toute question.</p>
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
</html>`

	textContent := `Bonjour,

Après étude de votre demande, nous sommes au regret de vous informer que votre candidature pour devenir contributeur sur "Le site du Patrimoine Saint-Jacques" a été rejetée.

Votre candidature a été supprimée de notre système.

Nous vous remercions de l’intérêt que vous portez à notre projet et restons à votre disposition pour toute question.

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}
