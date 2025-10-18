export async function writeToClipboard(content: string) {
  if (navigator.clipboard === undefined || window.isSecureContext === false) {
    return;
  }

  await navigator.clipboard.writeText(content);
}
