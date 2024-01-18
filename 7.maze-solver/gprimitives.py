from time import sleep

class Point():
  def __init__(self, x, y):
    self.x = x
    self.y = y

class Line():
  def __init__(self, start: Point, end: Point):
    self.start = start
    self.end = end
  
  def draw(self, canvas, fill_color):
    canvas.create_line(self.start.x, self.start.y, self.end.x, self.end.y, fill=fill_color, width = 2)
    canvas.pack()


class Cell():
  def __init__(self, top_left: Point, bottom_right: Point, window):
    self.__x1 = top_left.x
    self.__y1 = top_left.y
    self.__x2 = bottom_right.x
    self.__y2 = bottom_right.y
    self.has_left_wall = True
    self.has_right_wall = True
    self.has_top_wall = True
    self.has_bottom_wall = True
    self.__win = window
    self.visited = False
  
  def draw(self):

    bg_color = self.__win.background
    left_line = Line(Point(self.__x1, self.__y1), Point(self.__x1, self.__y2))
    if self.has_left_wall:
      self.__win.draw_line(left_line)
    else:
      self.__win.draw_line(left_line, bg_color)

    right_line = Line(Point(self.__x2, self.__y1), Point(self.__x2, self.__y2))
    if self.has_right_wall:
      self.__win.draw_line(right_line)
    else:
      self.__win.draw_line(right_line, bg_color)
    
    top_line = Line(Point(self.__x1, self.__y1), Point(self.__x2, self.__y1))
    if self.has_top_wall:
      self.__win.draw_line(top_line)
    else:
      self.__win.draw_line(top_line, bg_color)
    
    bottom_line = Line(Point(self.__x1, self.__y2), Point(self.__x2, self.__y2))
    if self.has_bottom_wall:
      self.__win.draw_line(bottom_line)
    else:
      self.__win.draw_line(bottom_line, bg_color)
  
  def draw_move(self, to_cell, undo = False):
    Point1 = Point((self.__x1 + self.__x2) / 2, (self.__y1 + self.__y2) / 2)
    Point2 = Point((to_cell.__x1 + to_cell.__x2) / 2, (to_cell.__y1 + to_cell.__y2) / 2)
    if undo:
      self.__win.draw_line(Line(Point1, Point2), "#aaa")
      sleep(0.0005)
    else:
      self.__win.draw_line(Line(Point1, Point2), "red")

  def _draw_emoji(self, emoji, font_size = 40):

    x = (self.__x1 + self.__x2) / 2
    y = (self.__y1 + self.__y2) / 2
    self.__win.get_canvas().create_text(x, y, text=emoji,font=('Arial', 40), anchor="center")
