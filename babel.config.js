module.exports = (api) => {
  // eslint-disable-next-line no-unused-vars
  const isTest = api.env("DEV");

  return {
    presets: [
      ["@babel/preset-env", { targets: { node: "current" } }],
      "@babel/preset-typescript"
    ]
  };
};
