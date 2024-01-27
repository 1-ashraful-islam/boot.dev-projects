const { JSDOM } = require('jsdom')

function normalizeURL(url) {

  let normal_url = url.hostname;
  normal_url += url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
  
  return normal_url;
}

function getURLsFromHTML(htmlBody, currentURL) {
  const urls = [];
  const dom = new JSDOM(htmlBody);
  const anchors = dom.window.document.querySelectorAll('a');

  anchors.forEach((anchor) => {
    try {
      const url= new URL(anchor.getAttribute('href'), currentURL);
      urls.push(url);
    } catch (e) {
      console.error(`Invalid URL: ${anchor.href}`);
    }
  });

  return urls;
}


async function fetchPage(currentURL) {

  return fetch(currentURL.href, {redirect: 'follow'})
    .then((response) => {
      if (!response.ok || response.status !== 200) {
        throw new Error(`Failed to fetch ${currentURL}. Status: ${response.status}`);
      }

      const contentType = response.headers.get('Content-Type');
      
      if (!contentType || !contentType.includes('text/html')) {
        throw new Error(`Failed to fetch ${currentURL}. Content-Type is not text/html`);
      }

      return response;
    })
    .then((response) => {
      // if (response.redirected) {
      //   console.log(`Redirected to ${response.url} from ${currentURL}`);
      // }
      // console.log(response.text())
      return response.text();
    });

}

async function crawlPage(currentURL, baseURL, pages, crawlDelay, robots, userAgent) {
  //check if the baseURL has same origin as currentURL
  if (baseURL.origin !== currentURL.origin) {
    return;
  }
  // Check if the URL is allowed to be crawled
  if (!robots.isAllowed(currentURL.href, userAgent)) {
    console.error(`Crawling disallowed for ${currentURL}`);
    return;
  }
  //check if the currentURL has already been crawled
  const normalizedURL = normalizeURL(currentURL);
  if (pages.has(normalizedURL)) {
    pages.set(normalizedURL, pages.get(normalizedURL) + 1);
    return pages;
  } else {
    pages.set(normalizedURL, baseURL === currentURL ? 0 :1);
  }

  console.log(`Crawling ${currentURL}`);

  try {
    //fetch the page
    let htmlBody = await fetchPage(currentURL);


    //get all urls from the page
    let urls = getURLsFromHTML(htmlBody, currentURL);
    //crawl each url from the page with delay
    for (let url of urls) {
      await new Promise(resolve => setTimeout(resolve, crawlDelay));
      await crawlPage(url, currentURL, pages, crawlDelay, robots, userAgent);
    }
  } catch (error) {
    console.error(`Error crawling ${currentURL}: ${error}`);
    //delete the currentURL from pages
    // pages.delete(normalizedURL);
  }

  return pages;
}



module.exports = {
  normalizeURL,
  getURLsFromHTML,
  fetchPage,
  crawlPage,
};

