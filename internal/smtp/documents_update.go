package smtp

import "context"

//nolint:lll
func (s *SMTPService) SendNewDocumentMail(ctx context.Context, to []EmailPerson) error {
	subject := "Un contributeur a soumis une nouvelle fiche !"

	htmlEmail := `
<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Nouvelle fiche soumise - Le site du Patrimoine Saint-Jacques</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
    <tr>
      <td align="center" style="padding:20px;">
        <a href="https://saintjacquesinfo.eu/#/" target="_blank" style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
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
              <p>Bonjour Administrateur,</p>
              <p>Un contributeur vient de soumettre une nouvelle <strong>fiche</strong> sur <em>Le site du Patrimoine Saint-Jacques</em>.</p>
              <p>Veuillez vous connecter à votre espace d’administration pour examiner et valider cette fiche avant sa publication.</p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding-top:20px;">
              <a href="https://saintjacquesinfo.eu/#/" 
                 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
                Examiner la fiche
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
</html>
`

	textContent := `Bonjour Administrateur,

Un contributeur vient de soumettre une nouvelle fiche sur "Le site du Patrimoine Saint-Jacques".

Veuillez vous connecter à votre espace d'administration pour l'examiner et la valider avant publication.
`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}

//nolint:lll
func (s *SMTPService) SendDocumentApprovedMail(ctx context.Context, to []EmailPerson, isUpdate bool) error {
	updateTxt := ""
	if isUpdate {
		updateTxt = "proposition de modification de"
	}

	subject := "Votre " + updateTxt + " fiche a été validée !"

	htmlEmail := `
<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Fiche validée - Le site du Patrimoine Saint-Jacques</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
	<tr>
	  <td align="center" style="padding:20px;">
		<a href="https://saintjacquesinfo.eu/#/" target="_blank" style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
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
			  <p>Bonne nouvelle 🎉 ! Votre` + updateTxt + ` fiche a été validée par un administrateur et est désormais visible sur <em>Le site du Patrimoine Saint-Jacques</em>.</p>
			  <p>Merci pour votre contribution à la valorisation du patrimoine Saint-Jacques !</p>
			</td>
		  </tr>
		  <tr>
			<td align="center" style="padding-top:20px;">
			  <a href="https://saintjacquesinfo.eu/#/" 
				 style="background-color:#802331; color:#ffffff; padding:12px 24px; border-radius:5px; text-decoration:none; font-weight:bold;">
				Voir la fiche
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
</html>
`

	textContent := `Bonjour,

Bonne nouvelle 🎉 ! Votre ` + updateTxt + ` fiche a été validée par un administrateur et est désormais visible sur "Le site du Patrimoine Saint-Jacques".

Merci pour votre contribution à la valorisation du patrimoine Saint-Jacques !

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}

//nolint:lll
func (s *SMTPService) SendDocumentRejectedMail(ctx context.Context, to []EmailPerson, isUpdate bool) error {
	updateTxt := ""
	if isUpdate {
		updateTxt = "proposition de modification de"
	}

	subject := "Votre " + updateTxt + " fiche n'a pas été approuvée"

	htmlEmail := `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Fiche rejetée - Le site du Patrimoine Saint-Jacques</title>
</head>
<body style="margin:0; padding:0; background-color:#fffcf8; font-family:Arial, sans-serif;">

  <!-- En-tête -->
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#802331;">
	<tr>
	  <td align="center" style="padding:20px;">
		<a href="https://saintjacquesinfo.eu/#/" target="_blank" style="color:#ffffff; font-size:20px; font-weight:bold; text-decoration:none;">
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
			  <p>Nous sommes désolés de vous informer que votre ` + updateTxt + ` fiche ne peut être approuvée dans l’état.</p>
              <p>Merci de prendre contact avec l’IRJ pour explications</p>
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

Nous sommes désolés de vous informer que votre ` + updateTxt + ` fiche ne peut être approuvée dans l’état.

Merci de prendre contact avec l’IRJ pour explications

Cet email est généré automatiquement. Merci de ne pas y répondre directement.
`

	return s.send(ctx, to, subject, htmlEmail, textContent)
}
