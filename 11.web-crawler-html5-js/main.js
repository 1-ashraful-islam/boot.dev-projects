const {argv} = require('node:process');
const robotsParser = require('robots-parser');

const {crawlPage} = require('./crawl');
const {printReport, fetchRobotsTxt} = require('./util');


async function main() {
  if (argv.length < 3 || argv.length > 3) {
    printHelp();
    return;
  }
  if (argv[2] === "help") {
    printHelp();
    return;
  }

  try {

    const BASE_URL = new URL(argv[2]);
    const robots = await fetchRobotsTxt(BASE_URL);
    const userAgent = 'Mozilla/5.0 (compatible; Crawlbot/1.0; +http://1-ashraful-islam.github.io/project/web-crawler)';

    // Crawl delay implementation
    const crawlDelay = robots?.getCrawlDelay(userAgent) * 1000 || 10;
  
    console.log(`Crawling ${BASE_URL}`);
    let pages = new Map();
    crawlPage(BASE_URL, BASE_URL, pages, crawlDelay, robots, userAgent)
      .then((pages) => {
        if (pages) {
          printReport(pages);
        }
        
      });

  } catch (error) {
    console.error(`Error crawling ${argv[2]}. Error: ${error}`);
    printHelp();
    return;
  }
}

function printHelp() {
  console.log("Usage: npm run start <BASE_URL>");
  console.log("Example: npm run start https://blog.boot.dev");
}

main();
