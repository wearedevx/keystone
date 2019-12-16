require("dotenv").config();
const mandrill = require("mandrill-api/mandrill");

const mandrill_client = new mandrill.Mandrill(process.env.MANDRILL_KEY);
const inviteUrl = process.env.INVITE_URL;

exports.mail = (req, res) => {
  // Set CORS headers for preflight requests
  // Allows GETs from any origin with the Content-Type header
  // and caches preflight response for 3600s

  res.set("Access-Control-Allow-Origin", "*");

  if (req.method === "OPTIONS") {
    // Send response to OPTIONS requests
    res.set("Access-Control-Allow-Methods", "GET");
    res.set("Access-Control-Allow-Headers", "Content-Type");
    res.set("Access-Control-Max-Age", "3600");
    res.status(204).send("");
  } else {
    console.log("body", req.body);
    const { request, project, email, id, from, to, role } = req.body;

    switch (request) {
      case "invite":
        // mails.map(mail => {
        sendMailRequestInvite(email, project, id, from, role);
        res.send("");
        // })

        break;
      case "accept":
        // const to = req.body.to

        // console.log(to,'|',from)

        sendMailAcceptInvite(project, id, to, from);
        res.send("");
        break;
      default:
        console.log("request", request);
        throw new Error("Unknown mail request");
    }
  }
};

function sendMailAcceptInvite(project, id, to, from) {
  // const url = `${inviteUrl}?action=accept&id=${id}&project=${project}&to=${to}`
  // const link = encodeURI(url)
  try {
    mandrill_client.messages.send(
      {
        message: {
          from_email: "contact@wearedevx.com",
          from_name: "Keystone",
          to: [{ email: to }],
          subject: `${from} is ready to join ${project}`,
          text: `${from} (${id}) has accepted your invitation. Type the following commands in your terminal to encrypt the files for him.`
        }
      },
      success => {
        console.log("success", success);
      },
      error => {
        console.log("error", error);
      }
    );
  } catch (error) {}
}

function sendMailRequestInvite(email, project, id, from, role) {
  const url = `${inviteUrl}?action=join&id=${id}&project=${project}&from=${from}&to=${email}`;
  const link = encodeURI(url);

  try {
    mandrill_client.messages.send(
      {
        message: {
          from_email: "contact@wearedevx.com",
          from_name: "Keystone",
          to: [{ email: email }],
          subject: `${from} invites you to join ${project}`,
          text: `You received an invite to the project ${project}. Click the link to join: ${link}`
        }
      },
      success => {
        console.log("success", success);
      },
      error => {
        console.log("error", error);
      }
    );
  } catch (error) {
    throw error;
  }
}
