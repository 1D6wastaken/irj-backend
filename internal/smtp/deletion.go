package smtp

import "context"

func (s *SMTPService) SendDeletionMail(ctx context.Context, to []EmailPerson, byAdmin bool) error {
	if byAdmin {
		return s.deletionByAdmin(ctx, to)
	}

	return s.deletionByUser(ctx, to)
}

//nolint:lll
func (s *SMTPService) deletionByAdmin(ctx context.Context, to []EmailPerson) error {
	subject := "Compte supprimé - Le site du Patrimoine Saint-Jacques"

	htmlEmail := `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Compte supprimé - Le site du Patrimoine Saint-Jacques</title>
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
              <p>Nous vous informons qu’un administrateur a procédé à la <strong>suppression</strong> de votre compte sur <em>Le site du Patrimoine Saint-Jacques</em>.</p>
              <p>Toutes vos données associées ont été retirées de notre système conformément à notre politique de gestion des comptes.</p>
              <p>Si vous pensez qu’il s’agit d’une erreur ou si vous souhaitez obtenir plus d’informations, vous pouvez nous contacter en écrivant à institut.irj@gmail.com.</p>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; text-align:center;">
              <a href="mailto:institut.irj@gmail.com" 
                 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
                Contacter l'équipe
              </a>
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

Nous vous informons qu’un administrateur a procédé à la suppression de votre compte sur "Le site du Patrimoine Saint-Jacques".

Toutes vos données associées ont été retirées de notre système conformément à notre politique de gestion des comptes.

Si vous pensez qu’il s’agit d’une erreur ou si vous souhaitez obtenir plus d’informations, vous pouvez nous contacter en écrivant à institut.irj@gmail.com

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}

//nolint:lll
func (s *SMTPService) deletionByUser(ctx context.Context, to []EmailPerson) error {
	subject := "Compte supprimé - Le site du Patrimoine Saint-Jacques"

	htmlEmail := `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Suppression de compte - Confirmation</title>
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
              <p>Conformément à votre demande, nous vous confirmons que votre compte sur <em>Le site du Patrimoine Saint-Jacques</em> a été <strong>supprimé</strong>.</p>
              <p>Toutes vos données personnelles associées ont été retirées de notre système, conformément à notre politique de confidentialité.</p>
              <p>Nous vous remercions d'avoir fait partie de notre communauté et espérons vous revoir à l'avenir.</p>
            </td>
          </tr>
          <tr>
            <td style="padding-top:20px; text-align:center;">
              <a href="https://test.saintjacquesinfo.eu/#/" 
                 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
                Revenir sur le site
              </a>
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

Conformément à votre demande, nous vous confirmons que votre compte sur "Le site du Patrimoine Saint-Jacques" a été supprimé.

Toutes vos données personnelles associées ont été retirées de notre système, conformément à notre politique de confidentialité.

Nous vous remercions d'avoir fait partie de notre communauté et espérons vous revoir à l'avenir.

Revenir sur le site : https://test.saintjacquesinfo.eu/#/

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}
