const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      on('before:run', () => {
        console.log('V2_BOT_URL:', process.env.V2_BOT_URL);
        console.log('V2_BOT_TOKEN:', process.env.V2_BOT_TOKEN);
      });

      on('uncaught:exception', (err, runnable) => {
        // returning false here prevents Cypress from
        // failing the test
        return false;
      });
    },
  },
});