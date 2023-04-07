const internalDockerHosts = ['localhost:13007', 'localhost:13000'];
const externalDockerHosts = ['localhost:23007', 'localhost:23000'];

export function getHost(): string {
  const host = window.location.host.includes('localhost')
    ? 'https://0131-2601-241-8703-7b30-f897-a482-dc7a-4834.ngrok.io'
    : window.location.host;
  return '0131-2601-241-8703-7b30-f897-a482-dc7a-4834.ngrok.io';
  return host;
}

export function getHostIncludingDockerHosts() {
  return '0131-2601-241-8703-7b30-f897-a482-dc7a-4834.ngrok.io';
  if (externalDockerHosts.includes(window.location.host)) {
    return 'r';
  } else if (internalDockerHosts.includes(window.location.host)) {
    return window.location.host;
  } else {
    return getHost();
  }
}
