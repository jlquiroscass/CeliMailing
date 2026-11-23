# CeliMailing

Servicio en Go que envía los correos de CeliAPI por Amazon SES (región `eu-west-1`). Todo está en `main.go`.

- **API:** `POST /mail` en el puerto 6565 con JSON `{subject, to, body (HTML), txtBody, sender, replyTo (opcional: dirección de «Responder a»)}`; responde con el mismo JSON. CeliAPI lo llama con `celicidad.mailing.url` (`http://mailing.celicidad.net:6565/mail`, o `http://mail:6565/mail` dentro de Docker); despliegue en `CeliDeployTools/Mailing/`.
- **Rama:** `develop` es la principal.
- **Compilar:** módulo `github.com/jlquiroscass/CeliMailing` (`go.mod`, Go 1.16+); el `Dockerfile` usa `golang:1.22-alpine` fija para que una versión nueva de Go no rompa la imagen y compila siempre a `/mail` (lo que arranca `ENTRYPOINT`), se llame como se llame el módulo. En local, `go build .`.
- **Errores:** si SES rechaza el correo (remitente no verificado, mensaje rechazado), solo lo escribe en el registro y responde sin cuerpo: CeliAPI no se entera. Ante correos que no llegan, mirar primero el registro de este servicio.
- **Sin autenticación ni límite:** cualquiera que llegue al puerto 6565 puede enviar correo como Celicidad. Debe quedar solo en la red interna.
- **Credenciales:** las de AWS salen del entorno o de `cred/credentials`; no subirlas a git ni copiarlas en otros archivos.
