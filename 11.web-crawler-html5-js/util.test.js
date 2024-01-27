const {printReport} = require('./util');

test('print report', () => {
  const pages = new Map();
  pages.set('blog.boot.dev', 1);
  pages.set('blog.boot.dev/courses', 2);
  pages.set('blog.boot.dev/courses/javascript', 3);
  pages.set('blog.boot.dev/courses/python', 4);
  pages.set('blog.boot.dev/courses/python/advanced', 5);

  const spy = jest.spyOn(console, 'log');
  printReport(pages);
  expect(spy).toHaveBeenCalledTimes(1);
  expect(spy).toHaveBeenCalledWith(pages);
});
