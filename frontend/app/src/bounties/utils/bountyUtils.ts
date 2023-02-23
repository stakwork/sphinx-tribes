export function sendToRedirect(url) {
  const el = document.createElement('a');
  el.href = url;
  el.target = '_blank';
  el.click();
}
