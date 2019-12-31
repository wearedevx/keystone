import path from 'path'

export default {
  devServer: {
    port: 8000,
    headers: {
      'Access-Control-Allow-Origin': '*',
    },
  },
  plugins: [
    [
      require.resolve('react-static-plugin-source-filesystem'),
      {
        location: path.resolve('./src/pages'),
      },
    ],
    require.resolve('react-static-plugin-reach-router'),
    require.resolve('react-static-plugin-tailwindcss'),
    require.resolve('@elbstack/react-static-plugin-dotenv'),
  ],
}
