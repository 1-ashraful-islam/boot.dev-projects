const { JSDOM } = require('jsdom')

function normalizeURL(url) {

  let normal_url = url.hostname;
  normal_url += url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
  
  return normal_url;
}

function getURLsFromHTML(htmlBody, currentURL, sameOrigin = true) {
  const urls = [];
  const dom = new JSDOM(htmlBody);
  const anchors = dom.window.document.querySelectorAll('a');

  anchors.forEach((anchor) => {
    try {
      const url= new URL(anchor.getAttribute('href'), currentURL);
      if (sameOrigin && url.origin !== currentURL.origin) {
        return;
      }
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

async function crawlPage(currentURL, baseURL, pages, crawlDelay, robots, userAgent, depth, graph) {
  //check if the baseURL has same origin as currentURL
  if (baseURL.origin !== currentURL.origin || depth < 0) {
    return;
  }
  // Check if the URL is allowed to be crawled
  if (robots && !robots.isAllowed(currentURL.href, userAgent)) {
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
    graph.addNode(currentURL.pathname);
    if (baseURL.pathname !== currentURL.pathname) {
      graph.addEdge(baseURL.pathname, currentURL.pathname);
    }
  }

  console.log(`Crawling ${currentURL}`);

  try {
    //fetch the page
    let htmlBody = await fetchPage(currentURL);

    //get all urls from the page
    let urls = getURLsFromHTML(htmlBody, currentURL, sameOrigin = true);

    //crawl each url from the page with delay
    for (let url of urls) {
      await new Promise(resolve => setTimeout(resolve, crawlDelay));
      await crawlPage(url, currentURL, pages, crawlDelay, robots, userAgent, depth - 1, graph);
    }
  } catch (error) {
    console.error(`Error crawling ${currentURL}: ${error}`);
    //delete the currentURL from pages if the error is 404
    // if (error.message.includes('404')) {
    //   pages.delete(normalizedURL);
    // }
  }

  return pages;
}

async function breadthFirstCrawl(baseURL, pages, crawlDelay, robots, userAgent, maxDepth, graph) {
  let queue = [{url: baseURL, refURL: baseURL, depth: 0}];

  let edges = new Map();


  while (queue.length > 0) {
    let {url, refURL, depth} = queue.shift(); // Get the first URL from the queue
    if (depth > maxDepth) continue;

    const normalizedURL = normalizeURL(url);
    if (pages.has(normalizedURL)) {
      pages.set(normalizedURL, pages.get(normalizedURL) + 1);
      continue;
    }
    pages.set(normalizedURL, 1);

    if (robots && !robots.isAllowed(url.href, userAgent)) {
      console.error(`Crawling disallowed for ${url}`);
      continue;
    }
    
    graph.addNode(url.pathname);
    if ((!edges.has(refURL.pathname + "->" + url.pathname) || !edges.has(url.pathname + "->" +refURL.pathname)) && refURL.pathname !== url.pathname) {
      graph.addEdge(refURL.pathname, url.pathname);
      edges.set(refURL.pathname + "->" + url.pathname, true);
    }

    console.log(`Crawling ${url}`);
    try {
      let htmlBody = await fetchPage(url);
      let urls = getURLsFromHTML(htmlBody, url, sameOrigin = true);

      urls.forEach(newUrl => {
        const normalizedNewUrl = normalizeURL(newUrl);
        if (newUrl.origin === baseURL.origin) {
          queue.push({url: newUrl, refURL: url, depth: depth + 1}); // Add new URLs to the queue
          
        }
      });
    } catch (error) {
      console.error(`Error crawling ${url}: ${error}`);
    }

    await new Promise(resolve => setTimeout(resolve, crawlDelay)); // Crawl delay
  }

  return {pages, graph};
}



module.exports = {
  normalizeURL,
  getURLsFromHTML,
  fetchPage,
  crawlPage,
  breadthFirstCrawl,
};

