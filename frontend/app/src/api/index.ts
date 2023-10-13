import { getHost } from '../config/host';

function addMethod(m: string): (url: string, data?: any, incomingHeaders?: any) => void {
  const host = getHost();
  const rootUrl =
    host.includes('localhost') || host.includes('internal')
      ? `http://${host}/`
      : `https://${host}/`;

  const func = async function (url: string, data: any, incomingHeaders: any) {
    try {
      const headers: { [key: string]: string } = {};
      const opts: { [key: string]: any } = { mode: 'cors' };
      if (m === 'POST' || m === 'PUT') {
        if (!incomingHeaders) {
          headers['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
          opts.body = new URLSearchParams(data);
        } else {
          if (
            incomingHeaders &&
            incomingHeaders['Content-Type'] &&
            incomingHeaders['Content-Type'] === 'application/json'
          ) {
            opts.body = JSON.stringify(data);
          }
        }
      }
      if (m === 'UPLOAD') {
        const file = data;
        const filename = file.name || 'name';
        const type = file.type || 'application/octet-stream';
        const formData = new FormData();
        formData.append('file', new Blob([file], { type }), filename);
        opts.body = formData;
      }
      opts.headers = new Headers(headers);
      opts.method = m === 'UPLOAD' ? 'POST' : m;
      if (m === 'BLOB') opts.method = 'GET';
      const r = await fetch(rootUrl + url, opts);
      if (!r.ok) {
        throw new Error('Not OK!');
      }
      let res;
      if (m === 'BLOB') res = await r.blob();
      else {
        res = await r.json();
      }
      return res;
    } catch (e) {
      throw e;
    }
  };
  return func;
}

class API {
  constructor() {
    this.get = addMethod('GET');
    this.post = addMethod('POST');
    this.put = addMethod('PUT');
    this.del = addMethod('DELETE');
  }
  get: (url: string, data?: any, incomingHeaders?: any) => any;
  post: (url: string, data?: any, incomingHeaders?: any) => any;
  put: (url: string, data?: any, incomingHeaders?: any) => any;
  del: (url: string, data?: any, incomingHeaders?: any) => any;
}

export default new API();
