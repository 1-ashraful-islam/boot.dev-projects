const { JSDOM } = require('jsdom')

function normalizeURL(url) {

  let normal_url = url.hostname;
  normal_url += url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
  
  return normal_url;
}

function getURLsFromHTML(htmlBody, baseURL) {
  const urls = [];
  const dom = new JSDOM(htmlBody);
  const anchors = dom.window.document.querySelectorAll('a');

  anchors.forEach((anchor) => {
    const url = new URL(anchor.href, baseURL);
    urls.push(url);
  });

  return urls;
}

module.exports = {
  normalizeURL,
  getURLsFromHTML
};

