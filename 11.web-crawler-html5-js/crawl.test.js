const {normalizeURL, getURLsFromHTML} = require('./crawl');

test('normalizes urls', () => {
  const url_list = [
    'https://blog.boot.dev/path/',
    'https://blog.boot.dev/path',
    'http://blog.boot.dev/path/',
    'http://blog.boot.dev/path'
  ];

  url_list.forEach((url_str) => {
    // Perform operations on each URL
    const url = new URL(url_str);
    expect(normalizeURL(url)).toBe('blog.boot.dev/path');
  });
});

test('get all urls from html', () => {
  const html = `
  <html>
  <body>
      <a href="https://blog.boot.dev"><span>Go to Boot.dev</span></a>
      <a href="/courses">Courses</a>
  </body>
  </html>
  `;
  
  const urls = getURLsFromHTML(html, "https://blog.boot.dev");
  expect(urls.length).toBe(2);
  expect(urls[0].href).toBe('https://blog.boot.dev/');
  expect(urls[1].href).toBe('https://blog.boot.dev/courses');

  //TODO: Test cases for when the href of the anchor is an invalid url string


});
