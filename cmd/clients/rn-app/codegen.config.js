/** @type {import('@graphql-codegen/cli').CodegenConfig} */
const config = {
  schema: process.env.GRAPHQL_ENDPOINT || 'http://localhost:8080/graphql',
  documents: ['graphql/**/*.graphql'],
  ignoreNoDocuments: true,
  generates: {
    './graphql/generated/graphql.ts': {
      plugins: [
        'typescript',
        'typescript-operations',
        'typescript-react-apollo',
      ],
      config: {
        withHooks: true,
      },
    },
  },
};

module.exports = config;
