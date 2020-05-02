class API {
  constructor(){
    this.get = addMethod('GET')
    this.post = addMethod('POST')
    this.put = addMethod('PUT')
    this.del = addMethod('DELETE')
  }
  get: Function;
  post: Function;
  put: Function;
  del: Function;
}

function addMethod(m: string): Function {
  const tokenName = 'token'
  const rootUrl = 'http://localhost:5001/'
  const func = async function (url: string, data: any, fields: any) {
    try {
      const token = await getToken(tokenName)
      const skip = isPublic(rootUrl + url)
      if (tokenName && !token && !skip) {
        throw new Error("no token")
      }
      const headers: {[key:string]:string} = {}
      if (tokenName && token) headers['authorization'] = 'Bearer ' + token
      const opts: {[key:string]:any} = { mode: 'cors' }
      if (m === 'POST' || m === 'PUT') {
        headers['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8'
        opts.body = new URLSearchParams(data)
      }
      if (m === 'UPLOAD') {
        const file = data
        const filename = file.name || 'name'
        const type = file.type || 'application/octet-stream'
        let formData = new FormData();
        formData.append('file', new Blob([file], { type }), filename)
        // Object.entries(fields).forEach(e => formData.append(e[0], e[1]))
        opts.body = formData
      }
      opts.headers = new Headers(headers)
      opts.method = m === 'UPLOAD' ? 'POST' : m
      if (m === 'BLOB') opts.method = 'GET'
      const r = await fetch(rootUrl + url, opts);
      if (!r.ok) {
        console.log(r)
        throw new Error('Not OK!');
      }
      let res
      if (m === 'BLOB') res = await r.blob()
      else {
        res = await r.json();
        if (res.token) {
          localStorage.setItem(tokenName, res.token)
        }
      }
      return res
    } catch (e) {
      throw e
    }
  }
  return func
}

function isPublic(url: string) {
  return url.endsWith('login') ||
    url.endsWith('errs/all')
}

export default new API()

async function getToken(name: string) {
  if (!name) return ""
  return localStorage.getItem(name)
}
