require("dotenv").config();
const mandrill = require("mandrill-api/mandrill");

const mandrill_client = new mandrill.Mandrill(process.env.MANDRILL_KEY);
const inviteUrl = process.env.INVITE_URL;

exports.mail = async (req, res) => {
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
    const { request, project, email, id, from, to, role, uuid } = req.body;

    switch (request) {
      case "invite":
        try {
          await sendInvite({
            email,
            project,
            id,
            from,
            uuid
          });
          res.send("");
        } catch (error) {
          res.status(500).send(error.message);
        }
        break;
      case "accept":
        try {
          await sendReceipt({ project, id, to, from, uuid });
          res.send();
        } catch (error) {
          res.status(500).send(error.message);
        }
        break;
      default:
        res.status(403).send("Request Unauthorized");
    }
  }
};

/**
 * Send an email when an invitee joined a project
 * @param {*} project
 * @param {*} id
 * @param {*} to
 * @param {*} from
 */
function sendReceipt({ project, id, to, from }) {
  const url = `${inviteUrl}?action=accept&id=${id}&project=${project}/${uuid}&to=${to}`;
  const link = encodeURI(url);
  return new Promise((resolve, reject) => {
    mandrill_client.messages.sendTemplate(
      {
        template_name: "keystone",
        template_content: [],
        message: {
          from_email: "contact@wearedevx.com",
          from_name: "Keystone",
          to: [{ email: to }],
          subject: `${from} is ready to join ${project}`,
          merge_language: "handlebars",
          global_merge_vars: [
            {
              name: "headline",
              content: `Project ${project}`
            },
            {
              name: "tags",
              content: `#${uuid} #${id}`
            },
            {
              name: "content",
              content: `<strong>${from}</strong> has accepted your invitation. <br><br>You need to set permissions for him so he can access files from your project.
              <br><br>
              <a href="${link}">Click here</a> to start.`
            }
          ]
        }
      },
      success => {
        resolve(success);
      },
      error => {
        reject(error);
      }
    );
  });
}

async function sendInvite({ email, project, id, from, uuid }) {
  const url = `${inviteUrl}?action=join&id=${id}&project=${project}/${uuid}&from=${from}&to=${email}`;
  const link = encodeURI(url);

  return new Promise((resolve, reject) => {
    mandrill_client.messages.sendTemplate(
      {
        template_name: "keystone",
        template_content: [],
        message: {
          from_email: "contact@wearedevx.com",
          from_name: "Keystone",
          to: [{ email: email }],
          subject: `${from} invites you to project ${project}`,
          merge_language: "handlebars",
          global_merge_vars: [
            {
              name: "headline",
              content: `Project ${project}`
            },
            {
              name: "tags",
              content: `#${uuid} #${id}`
            },
            {
              name: "content",
              content: `This invite is sent by <strong>${from}</strong>. <br><br><a href="${link}">Click here to join</a> the project or ignore if you don't know the sender.`
            }
          ]
        }
      },
      success => {
        resolve(success);
      },
      error => {
        reject(error);
      }
    );
  });
}
