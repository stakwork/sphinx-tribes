module.exports = (on, config) => {

    on('before:browser:launch', (browser = {}, args) => {
  
      if (browser.name === 'chrome') {
        args.push('--remote-debugging-port=9222')
  
        // whatever you return here becomes the new args
        return args
      }
  
    })
  }