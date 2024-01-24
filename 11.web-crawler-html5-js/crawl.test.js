const {normalizeURL} = require('./crawl');

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
