const robotsParser = require('robots-parser');
const fs = require('fs');

function printReport(pages, outputFile) {
  const sortedPages = new Map([...pages.entries()].sort((a, b) => b[1] - a[1]));

  // Generate CSV string
  let csv = 'Page,Visits\n';
  sortedPages.forEach((count, page) => {
    csv += `${page},${count}\n`;
  });

  // Save CSV to file
  fs.writeFileSync(outputFile, csv);

  console.log('Report saved as report.csv');
}

async function fetchRobotsTxt(url) {
  try {
    const robotsUrl = new URL('/robots.txt', url.origin).href;
    const response = await fetch(robotsUrl);
    const robotsTxt = await response.text();
    return robotsParser(robotsUrl, robotsTxt);
  } catch (error) {
    console.error('Error fetching robots.txt:', error);
    return null;
  }
}


module.exports = {
  printReport,
  fetchRobotsTxt
};
