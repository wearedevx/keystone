const conflictedDescriptors = {
  base: {
    content: [
      'NODE_ENV=dev',
      'SEND_MAIL=false',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_QUEUE=your-queue-name',
      'GCLOUD_LOCATION=europe-west1',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 1,
  },

  left: {
    content: [
      'NODE_ENV=dev-test',
      'SEND_MAIL=true',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_LOCATION=europe-west3',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 2,
  },
  right: {
    content: [
      'NODE_ENV=dev',
      'SEND_MAIL=false',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace-dev',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_QUEUE=your-queue-name',
      'GCLOUD_LOCATION=europe-west2',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 2,
  },
}

const conflictLessDescriptors = {
  base: {
    content: [
      'NODE_ENV=dev',
      'SEND_MAIL=false',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_QUEUE=your-queue-name',
      'GCLOUD_LOCATION=europe-west1',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 1,
  },
  left: {
    content: [
      'NODE_ENV=dev-test',
      'SEND_MAIL=true',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_QUEUE=your-queue-name',
      'GCLOUD_LOCATION=europe-west1',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 2,
  },
  right: {
    content: [
      'NODE_ENV=dev',
      'SEND_MAIL=false',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=10000',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace-2',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_QUEUE=your-queue-name',
      'GCLOUD_LOCATION=europe-west2',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 2,
  },
}

const conflictedEnvDescriptors = {
  base: {
    content: [
      'NODE_ENV=dev',
      'SEND_MAIL=false',
      '',
      'GOOGLE_APPLICATION_CREDENTIALS_FILE=config/dev/dev-datastore-credentials.json',
      '',
      'LOGIN_SALT=$2b$10$K7I4ahLkUpv1CWxeWrvPGO',
      'JWT_SECRET=jwtjsonsecretdedev',
      'JWT_TOKEN_EXPIRE_TIME=',
      '',
      'GCLOUD_DATASTORE_NAMESPACE=your-namespace',
      'GCLOUD_PROJECT=your-google-project-id',
      'GCLOUD_QUEUE=your-queue-name',
      'GCLOUD_LOCATION=europe-west1',
      'GCLOUD_APPENGINE_SERVICE=your-appengine-service',
    ].join('\n'),
    version: 1,
  },

  left: {
    content: {
      name: 'default',
      files: [
        {
          name: 'foo.txt',
          checksum: 'f7664840578baca90505f863d7cd0d88b3a41b37',
        },
        {
          checksum: 'c0c2afe4cc16b67b25497f1f46d44a9bfdf18f35',
          name: 'foo2.txt',
        },
        {
          name: 'toto.txt',
          checksum: '63ee5b8b60d7b15ffdb75cb29eb81c5c510ba2c0',
        },
        {
          checksum: '941ea76661f21eff7b3b127f2b96ae28310eb86e',
          name: 'test.txt',
        },
        {
          name: 'goo.txt',
          checksum: 'bf7e93b3b3523b1b9676e511048ab073b62458ce',
        },
      ],
    },
    version: 2,
  },
  right: {
    content: {
      name: 'default',
      files: [
        {
          name: 'foo.txt',
          checksum: 'f7664840578baca90505f863d7cd0d88b3a41b37',
        },
        {
          name: 'bar.txt',
          checksum: 'fe479fd418e19ed04c2dd850387d9a7a9a7c226c',
        },
        {
          checksum: 'c0c2afe4cc16b67b25497f1f46d44a9bfdf18f35',
          name: 'foo2.txt',
        },
        {
          name: 'toto.txt',
          checksum: '63ee5b8b60d7b15ffdb75cb29eb81c5c510ba2c0',
        },
        {
          name: 'goo.txt',
          checksum: 'bf7e93b3b3523b1b9676e511048ab073b62458ce',
        },
      ],
    },
    version: 2,
  },
}
module.exports = {
  conflictedDescriptors,
  conflictLessDescriptors,
  conflictedEnvDescriptors,
}
