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
    try {
      const url = new URL(anchor.href, baseURL);
      urls.push(url);
    } catch (e) {
      console.error(`Invalid URL: ${anchor.href}`);
    }
  });

  return urls;
}

async function fetchPage(currentURL) {

  return fetch(currentURL)
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
      // console.log(response.text())
      return response.text();
    });

}

async function crawlPage(baseURL, currentURL, pages) {
  //check if the baseURL has same origin as currentURL
  if (baseURL.origin !== currentURL.origin) {
    return pages;
  }
  //check if the currentURL has already been crawled
  const normalizedURL = normalizeURL(currentURL);
  if (pages.has(normalizedURL)) {
    pages.set(normalizedURL, pages.get(normalizedURL) + 1);
    return pages;
  } else {
    pages.set(normalizedURL, baseURL === currentURL ? 0 :1);
  }

  try {
    //fetch the page
    let htmlBody = await fetchPage(currentURL);
    //get all urls from the page
    let urls = getURLsFromHTML(htmlBody, currentURL);
    //crawl each url from the page
    for (let url of urls) {
      await crawlPage(baseURL, url, pages);
    }
  } catch (error) {
    console.error(`Error crawling ${currentURL}: ${error}`);
  }

  return pages;
}

function printReport(pages) {
  const sortedPages = new Map([...pages.entries()].sort((a, b) => b[1] - a[1]));
  console.log(sortedPages);
}



module.exports = {
  normalizeURL,
  getURLsFromHTML,
  fetchPage,
  crawlPage,
  printReport
};

