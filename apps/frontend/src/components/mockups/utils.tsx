export const formatWhatsAppText = (text: string) => {
  if (!text) return null;
  const escapeHtml = (str: string) => {
    const map: Record<string, string> = { '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' };
    return str.replace(/[&<>'"]/g, (tag) => map[tag] || tag);
  };
  let html = escapeHtml(text);
  const paymentLinkHtml = `<a href="#" class="text-[#53bdeb] underline decoration-[#53bdeb]/30 underline-offset-2 cursor-pointer break-all">https://rzp.io/i/fB9x2pL</a>`;
  html = html.replace(/\[PAYMENT_LINK\]/g, paymentLinkHtml);
  html = html.replace(/```([\s\S]*?)```/g, '<code class="font-mono text-[12.5px] bg-black/15 px-1 py-0.5 rounded">$1</code>');
  html = html.replace(/\*([^*]+)\*/g, '<strong>$1</strong>');
  html = html.replace(/_([^_]+)_/g, '<em>$1</em>');
  html = html.replace(/~([^~]+)~/g, '<del>$1</del>');
  return <span dangerouslySetInnerHTML={{ __html: html }} />;
};

export const formatSmsText = (text: string) => {
  if (!text) return null;
  const escapeHtml = (str: string) => {
    const map: Record<string, string> = { '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' };
    return str.replace(/[&<>'"]/g, (tag) => map[tag] || tag);
  };
  let html = escapeHtml(text);
  const paymentLinkHtml = `<a href="#" class="text-[#0a84ff] underline underline-offset-2 cursor-pointer break-all">https://rzp.io/i/fB9x2pL</a>`;
  html = html.replace(/\[PAYMENT_LINK\]/g, paymentLinkHtml);
  return <span dangerouslySetInnerHTML={{ __html: html }} />;
};

export const parseDraftBody = (bodyObj: any) => {
  const rawString = typeof bodyObj === 'string' ? bodyObj : (bodyObj?.String || '');
  if (!rawString) return 'Generating draft...';
  try {
    const parsed = JSON.parse(rawString);
    return parsed.body || rawString;
  } catch {
    return rawString;
  }
};
