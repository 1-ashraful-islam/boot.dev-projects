function normalizeURL(url) {

  let normal_url = url.hostname;
  normal_url += url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
  
  return normal_url;
}

module.exports = {
  normalizeURL
};

