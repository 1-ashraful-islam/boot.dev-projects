const {normalizeURL, getURLsFromHTML, fetchPage} = require('./crawl');

afterEach(() => {
  jest.restoreAllMocks();
});

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

test('fetchPage fetch successful html',async () => {
  
  const responseBody = '<html><body><h1>Test</h1></body></html>';
  const responseOptions = {
    status: 200,
    ok: true,
    headers: {
      'Content-Type': 'text/html'
    },
  };

  global.fetch = jest.fn();

  fetch.mockReturnValue(
    Promise.resolve(
      new Response(responseBody, responseOptions)
    )
  );

  const output = await fetchPage('https://blog.boot.dev');
  expect(output).toBe('<html><body><h1>Test</h1></body></html>');
  expect(fetch).toHaveBeenCalledWith('https://blog.boot.dev');

});


test('fetchPage throws 400+ error', async () => {
  const responseBody = 'Bad Request';
  const responseOptions = {
    status: 400,
    ok: false,
    headers: {
      'Content-Type': 'text/html'
    },
  };

  global.fetch = jest.fn();

  fetch.mockReturnValue(
    Promise.resolve(
      new Response(responseBody, responseOptions)
    )
  );

  await expect(() => fetchPage('https://blog.boot.dev')).rejects.toThrow('Failed to fetch https://blog.boot.dev. Status: 400');
});

test('fetchPage throws content not text/html', async () => {
  const responseBody = 'Bad Request';
  const responseOptions = {
    status: 200,
    ok: true,
    headers: {
      'Content-Type': 'text/plain'
    },
  };

  global.fetch = jest.fn();

  fetch.mockReturnValue(
    Promise.resolve(
      new Response(responseBody, responseOptions)
    )
  );

  await expect(() => fetchPage('https://blog.boot.dev')).rejects.toThrow('Failed to fetch https://blog.boot.dev. Content-Type is not text/html');
});


//TODO: test fetchPage with invalid URL string
