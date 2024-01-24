const {argv} = require('node:process')

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
    const BASE_URL = new URL(argv[2]);
    console.log(`Crawling ${BASE_URL}`);
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
