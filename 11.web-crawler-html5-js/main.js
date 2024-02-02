const {argv} = require('node:process');
const robotsParser = require('robots-parser');
const graphviz = require('graphviz');
const yargs = require('yargs');

const {crawlPage, breadthFirstCrawl} = require('./crawl');
const {printReport, fetchRobotsTxt} = require('./util');


async function main() {
  const args = yargs
    .usage('Usage: npm run start -- --url <BASE_URL> [--depth <DEPTH>] [--bfs] [--userAgent <USER_AGENT>]')
    .option('url', {
      describe: 'The base URL to crawl',
      demandOption: true,
      type: 'string'
    })
    .option('depth', {
      describe: 'The maximum depth to crawl',
      type: 'number',
      default: 20
    })
    .option('bfs', {
      describe: 'Use breadth-first search instead of depth-first search',
      type: 'boolean',
      default: false
    })
    .option('userAgent', {
      describe: 'The user agent string to use for crawling',
      type: 'string',
      default: 'Mozilla/5.0 (compatible; Crawlbot/1.0;)'
    })
    .help()
    .argv;

  try {
    const BASE_URL = new URL(args.url);
    const depth = parseInt(args.depth);
    const depthFirst = !args.bfs;
    const userAgent = args.userAgent;
    const robots = await fetchRobotsTxt(BASE_URL);

    // Create a new digraph.
    var g = graphviz.digraph("G");

    // Crawl delay implementation
    const crawlDelay = robots?.getCrawlDelay(userAgent) * 1000 || 10;

    console.log(`Crawling ${BASE_URL}`);
    let pages = new Map();
    if (depthFirst) {
      await crawlPage(BASE_URL, BASE_URL, pages, crawlDelay, robots, userAgent, depth, g);
        
    } else {
      await breadthFirstCrawl(BASE_URL, pages, crawlDelay, robots, userAgent, depth, g)
    }
    
    if (pages) {
      printReport(pages, 'output/report.csv');
      //print the dot script
      console.log(g.to_dot());
      // Render the graph to an image
      g.render("png", "output/graph.png");
    }
  

  } catch (error) {
    console.error(`Error crawling ${args.url}. Error: ${error}`);
    yargs.showHelp();
    return;
  }
}

main();
