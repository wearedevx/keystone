module.exports = {
  ERROR_CODES: {
    PullWhileFilesModified: 'PullWhileFilesModified',
    InvitationFailed: 'InvitationFailed',
    ProjectNameExists: 'ProjectNameExists',
    ConfigFileExists: 'ConfigFileExists',
    FailedToFetch: 'FailedToFetch',
    PullBeforeYouPush: 'PullBeforeYouPush',
    NeedToBeAdminOrContributor: 'NeedToBeAdminOrContributor',
    NeedToBeAdmin: 'NeedToBeAdmin',
    MissingParams: 'MissingParams',
    AccountMismatch: 'AccountMismatch',
    NoUsername: 'NoUsername',
    InvalidProjectName: 'InvalidProjectName',
    Conflict: 'Conflict',
    MissingEnv: 'MissingEnv',
    UserNotInProject: 'UserNotInProject'
  },
  KEYSTONE_MAIL:
    process.env.KEYSTONE_MAIL ||
    'https://us-central1-keystone-245200.cloudfunctions.net/keystone-mail',
  KEYSTONE_WEB: process.env.KEYSTONE_WEB || 'https://keystone.sh',
  INVITATIONS_STORE: 'invitations.json',
  ROLES: {
    ADMINS: 'admins',
    CONTRIBUTORS: 'contributors',
    READERS: 'readers',
  },
  PUBKEY: 'public.key',
  KEYSTONE_CONFIG_PATH: '.ksconfig',
  KEYSTONE_ENV_CONFIG_PATH: 'envconfig',
  PROJECTS_STORE: process.env.PROJECTS_STORE || 'projects.json',
  KEYSTONE_HIDDEN_FOLDER: process.env.KEYSTONE_HIDDEN_FOLDER || '.keystone',
  SHARED_MEMBER: '{{shared}}',
  LOGIN_KEY_PREFIX: '{{login}}',
  SHARE_FILENAME: 'keystone-link.json',
  SESSION_FILENAME: process.env.SESSION_FILENAME || 'session.json',
  EVENTS: { CONFLICT: 'CONFLICT' },
}
