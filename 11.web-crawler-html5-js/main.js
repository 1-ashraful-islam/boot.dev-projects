const {argv} = require('node:process');
const { crawlPage, printReport} = require('./crawl');

function main() {
  if (argv.length < 3 || argv.length > 3) {
    printHelp();
    return;
  }
  if (argv[2] === "help") {
    printHelp();
    return;
  }

  try {

    //TODO: check robots.txt before crawling
    const BASE_URL = new URL(argv[2]);
    console.log(`Crawling ${BASE_URL}`);
    let pages = new Map();
    crawlPage(BASE_URL, BASE_URL, pages)
      .then((pages) => {
        printReport(pages);
      });

  } catch (error) {
    console.log("Invalid URL");
    printHelp();
    return;
  }
}

function printHelp() {
  console.log("Usage: npm run start <BASE_URL>");
  console.log("Example: npm run start https://blog.boot.dev");
}

main();
